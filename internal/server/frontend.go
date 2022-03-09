package server

import (
	"context"
	"fmt"
	"github.com/dgraph-io/badger/v3"
	"github.com/google/uuid"
	"github.com/gorymoon/adoptahero-ebs/internal/jwt"
	pb "github.com/gorymoon/adoptahero-ebs/internal/protos"
	"github.com/nicklaw5/helix"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) RequestServiceJWT(ctx context.Context, msg *pb.RequestJWTMessage) (*pb.JWTResponse, error) {
	client := s.CreateTwitchClient()
	client.SetUserAccessToken(msg.Token)
	users, err := client.GetUsers(&helix.UsersParams{})
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	user := users.Data.Users[0]

	token, err := jwt.NewFrontendJWT(user.DisplayName, user.ID, msg.Token, s.issuer, s.secret)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &pb.JWTResponse{
		Token: token,
	}, nil
}

func (s *Server) GetHeroData(_ context.Context, msg *pb.RequestHeroMessage) (*pb.HeroData, error) {
	key := fmt.Sprintf("%s_%s", msg.GetChannel(), msg.GetName())

	hero, err := s.heroStore.GetHero(key)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return hero, nil
}

func (s *Server) GetConnectionStatus(_ context.Context, msg *pb.ConnectionStatusMessage) (*pb.ConnectionStatusResponse, error) {
	channel, err := s.channelStore.GetChannel(msg.Channel)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &pb.ConnectionStatusResponse{
		Channel:   channel.Id,
		Connected: channel.Connected,
	}, nil
}

func (s *Server) GetGameJWT(ctx context.Context, _ *pb.RequestGameJWTMessage) (*pb.JWTResponse, error) {
	channelID := GetFrontendClaimFromContext(ctx).Subject
	channel, err := s.channelStore.GetChannel(channelID)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	token, err := jwt.NewGameJWT(channel.Id, channel.Uuid, s.issuer, s.secret)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &pb.JWTResponse{
		Token: token,
	}, nil
}

func (s *Server) NewGameJWT(ctx context.Context, _ *pb.RequestGameJWTMessage) (*pb.JWTResponse, error) {
	uuidToken := uuid.NewString()
	claim := GetFrontendClaimFromContext(ctx)

	channelID := claim.Subject
	channel, err := s.channelStore.GetChannel(channelID)
	if err != nil {
		if err == badger.ErrKeyNotFound {
			_, err := s.channelStore.NewChannel(channelID, claim.Name, uuidToken, false)
			if err != nil {
				return nil, status.Error(codes.InvalidArgument, err.Error())
			}
		} else {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
	} else {
		channel.Uuid = uuidToken
		err = s.channelStore.SetChannel(channelID, channel)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
	}

	token, err := jwt.NewGameJWT(channelID, uuidToken, s.issuer, s.secret)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &pb.JWTResponse{
		Token: token,
	}, nil
}
