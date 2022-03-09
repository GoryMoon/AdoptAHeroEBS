package server

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	jwt2 "github.com/gorymoon/adoptahero-ebs/internal/jwt"
	pb "github.com/gorymoon/adoptahero-ebs/internal/protos"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"time"
)

func (s *Server) UpdateData(stream pb.GameConnection_UpdateDataServer) error {
	startTime := time.Now()
	heroes := make(map[string]*pb.HeroData)
	ctx := stream.Context()
	channel := GetGameClaimFromContext(ctx).Subject
	err := s.channelStore.SetConnectionOnChannel(channel, true)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}
	log.Info().Str("ctx", "game").Str("channel", channel).Msg("Connection opened")

	running := true
	for running {
		select {
		case <-ctx.Done():
			running = false
			continue
		default:
		}

		msg, err := stream.Recv()
		if err == io.EOF {
			running = false
			continue
		}
		if err != nil {
			log.Error().Str("ctx", "game").Err(err).Msg("Error receiving hero data")
			continue
		}

		// If the current batch is done store it
		if msg.GetBatchDone() {
			for k := range heroes {
				log.Info().Str("hero", k).Msgf("Updating hero")
			}
			err := s.heroStore.SetHeroes(heroes)
			if err != nil {
				return status.Error(codes.InvalidArgument, err.Error())
			}
			heroes = make(map[string]*pb.HeroData)
			continue
		}

		key := fmt.Sprintf("%s_%s", channel, msg.GetData().GetName())
		heroes[key] = msg.GetData()
	}

	endTime := time.Now()
	log.Info().Str("ctx", "game").Str("channel", channel).TimeDiff("duration", endTime, startTime).Msg("Connection closed")
	err = stream.Send(&pb.CountResponse{
		Count:       int32(len(heroes)),
		ElapsedTime: int32(endTime.Sub(startTime).Seconds()),
	})

	channelErr := s.channelStore.SetConnectionOnChannel(channel, false)
	if channelErr != nil {
		return status.Error(codes.Internal, channelErr.Error())
	}
	return err
}

func (s *Server) RemoveHeroes(stream pb.GameConnection_RemoveHeroesServer) error {
	startTime := time.Now()
	var names []string
	channel := GetGameClaimFromContext(stream.Context())

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			err := s.heroStore.DeleteHeroes(names)
			if err != nil {
				return status.Error(codes.InvalidArgument, err.Error())
			}

			endTime := time.Now()
			return stream.SendAndClose(&pb.CountResponse{
				Count:       int32(len(names)),
				ElapsedTime: int32(endTime.Sub(startTime).Seconds()),
			})
		}
		if err != nil {
			return err
		}
		names = append(names, fmt.Sprintf("%s_%s", channel, msg.GetName()))
	}
}

func GetGameClaimFromContext(ctx context.Context) *jwt.RegisteredClaims {
	value := ctx.Value(ContextKeyGameClaim{})
	if value == nil {
		log.Error().Str("ctx", "game").Msg("Channel id metadata error")
		return nil
	}
	return value.(*jwt.RegisteredClaims)
}

func GetFrontendClaimFromContext(ctx context.Context) *jwt2.FrontendJWT {
	name := ctx.Value(ContextKeyFrontendClaim{})
	if name == nil {
		log.Error().Str("ctx", "game").Msg("Channel name metadata error")
		return nil
	}
	return name.(*jwt2.FrontendJWT)
}
