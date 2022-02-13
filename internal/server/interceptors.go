package server

import (
	"context"
	"crypto/subtle"
	"github.com/gorymoon/adoptahero-ebs/internal/jwt"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strings"
	"time"
)

const (
	authHeaderPrefix string = "Bearer "
)

type ContextKeyGameClaim struct{}
type ContextKeyFrontendClaim struct{}

type StreamContextWrapper interface {
	grpc.ServerStream
	SetContext(context.Context)
}

type wrappedContextStream struct {
	grpc.ServerStream
	ctx context.Context
}

func newWrappedContextStream(s grpc.ServerStream) StreamContextWrapper {
	return &wrappedContextStream{s, s.Context()}
}

func (w *wrappedContextStream) Context() context.Context {
	return w.ctx
}

func (w *wrappedContextStream) SetContext(ctx context.Context) {
	w.ctx = ctx
}

func (s *Server) unaryAuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	start := time.Now()

	if info.FullMethod != "/blt.adoptahero.Frontend/RequestServiceJWT" {
		if ctx, err = s.authorizeFrontend(ctx); err != nil {
			return nil, err
		}
	}

	h, err := handler(ctx, req)

	log.Info().
		Str("ctx", "interceptor").
		Str("method", info.FullMethod).
		TimeDiff("duration", time.Now(), start).
		Err(err).
		Msg("GRPC Request")

	return h, err
}

// Check for jwt, add a metadata flag for mod or twitch
func (s *Server) streamAuthInterceptor(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	start := time.Now()

	w := newWrappedContextStream(stream)

	if err := s.authorizeGame(w); err != nil {
		return err
	}

	err := handler(srv, w)

	log.Info().
		Str("ctx", "interceptor").
		Str("method", info.FullMethod).
		TimeDiff("duration", time.Now(), start).
		Err(err).
		Msg("GRPC Stream Request")

	return err
}

func (s *Server) authorizeFrontend(ctx context.Context) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx, status.Errorf(codes.InvalidArgument, "Retrieving metadata is failed")
	}

	authHeader, ok := md["authorization"]
	if !ok {
		return ctx, status.Errorf(codes.Unauthenticated, "Authorization token is not supplied")
	}

	token := strings.TrimPrefix(authHeader[0], authHeaderPrefix)
	claims, err := jwt.VerifyFrontendJWT(token, s.secret)

	if err != nil {
		return ctx, status.Errorf(codes.Unauthenticated, err.Error())
	}

	if !claims.VerifyIssuer(s.issuer, true) {
		return ctx, status.Error(codes.Unauthenticated, "token is invalid")
	}

	return context.WithValue(ctx, ContextKeyFrontendClaim{}, claims), nil
}

func (s *Server) authorizeGame(w StreamContextWrapper) error {
	md, ok := metadata.FromIncomingContext(w.Context())
	if !ok {
		return status.Errorf(codes.InvalidArgument, "Retrieving metadata is failed")
	}

	authHeader, ok := md["authorization"]
	if !ok {
		return status.Errorf(codes.Unauthenticated, "Authorization token is not supplied")
	}

	token := strings.TrimPrefix(authHeader[0], authHeaderPrefix)
	claims, err := jwt.VerifyGameJWT(token, s.secret)

	if err != nil {
		return status.Errorf(codes.Unauthenticated, err.Error())
	}

	channel, err := s.channelStore.GetChannel(claims.Subject)
	if err != nil {
		return status.Error(codes.NotFound, err.Error())
	}

	if !claims.VerifyIssuer(s.issuer, true) {
		return status.Error(codes.Unauthenticated, "token is invalid")
	}

	if subtle.ConstantTimeCompare([]byte(claims.ID), []byte(channel.Uuid)) == 0 {
		return status.Error(codes.Unauthenticated, "token doesn't match")
	}

	w.SetContext(context.WithValue(w.Context(), ContextKeyGameClaim{}, claims))
	return nil
}
