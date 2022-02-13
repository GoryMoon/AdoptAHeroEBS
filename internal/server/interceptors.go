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
	"time"
)

func (s *Server) unaryAuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	start := time.Now()

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

	if err := s.authorize(stream.Context()); err != nil {
		return err
	}

	err := handler(srv, stream)

	log.Info().
		Str("ctx", "interceptor").
		Str("method", info.FullMethod).
		TimeDiff("duration", time.Now(), start).
		Err(err).
		Msg("GRPC Stream Request")

	return err
}

func (s *Server) authorize(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Errorf(codes.InvalidArgument, "Retrieving metadata is failed")
	}

	authHeader, ok := md["authorization"]
	if !ok {
		return status.Errorf(codes.Unauthenticated, "Authorization token is not supplied")
	}

	token := authHeader[0]
	claims, err := jwt.VerifyGameJWT(token, s.secret)

	if err != nil {
		return status.Errorf(codes.Unauthenticated, err.Error())
	}

	channel, err := s.channelStore.GetChannel(claims.Subject)
	if err != nil {
		return status.Error(codes.NotFound, err.Error())
	}

	if subtle.ConstantTimeCompare([]byte(claims.ID), []byte(channel.Uuid)) == 0 {
		return status.Error(codes.Unauthenticated, "")
	}

	md.Set("channelID", channel.Id)
	return nil
}
