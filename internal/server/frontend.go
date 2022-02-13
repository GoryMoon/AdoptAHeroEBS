package server

import (
	"context"
	"fmt"
	"github.com/dgraph-io/badger/v3"
	"github.com/google/uuid"
	"github.com/gorymoon/adoptahero-ebs/internal/jwt"
	pb "github.com/gorymoon/adoptahero-ebs/internal/protos"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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

func (s *Server) RequestGameJWT(_ context.Context, msg *pb.RequestJWTMessage) (*pb.JWTResponse, error) {
	uuidToken := uuid.NewString()

	channel, err := s.channelStore.GetChannel(msg.Channel)
	if err != nil {
		if err == badger.ErrKeyNotFound {
			_, err := s.channelStore.NewChannel(msg.GetChannel(), msg.GetName(), uuidToken, false)
			if err != nil {
				return nil, status.Error(codes.InvalidArgument, err.Error())
			}
		} else {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
	} else {
		channel.Uuid = uuidToken
		err = s.channelStore.SetChannel(msg.GetChannel(), channel)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
	}

	token, err := jwt.NewGameJWT(msg.GetChannel(), uuidToken, s.secret)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &pb.JWTResponse{
		Token: token,
	}, nil
}
