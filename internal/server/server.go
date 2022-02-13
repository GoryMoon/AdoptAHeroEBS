package server

import (
	"fmt"
	pb "github.com/gorymoon/adoptahero-ebs/internal/protos"
	"github.com/gorymoon/adoptahero-ebs/internal/stores"
	"github.com/gorymoon/adoptahero-ebs/pkg/db"
	"github.com/nicklaw5/helix"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"net"
)

type Server struct {
	port         int
	host         string
	secret       []byte
	issuer       string
	twitchClient string
	twitchSecret string
	kvDV         *db.KvDB
	channelStore *stores.ChannelStore
	heroStore    *stores.HeroStore
	grpcServer   *grpc.Server
	twitch       *helix.Client
	pb.UnimplementedGameConnectionServer
	pb.UnimplementedFrontendServer
}

func (s *Server) Run() error {
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.host, s.port))
	if err != nil {
		return errors.Wrap(err, "Failed to listen")
	}

	s.grpcServer = grpc.NewServer(
		grpc.UnaryInterceptor(s.unaryAuthInterceptor),
		grpc.StreamInterceptor(s.streamAuthInterceptor),
	)
	pb.RegisterGameConnectionServer(s.grpcServer, s)
	pb.RegisterFrontendServer(s.grpcServer, s)

	log.Info().Str("ctx", "server").Str("host", listen.Addr().String()).Msg("Grpc listening")

	if err := s.grpcServer.Serve(listen); err != nil {
		log.Error().Err(err).Msg("Failed to serve: ")
		return err
	}
	return nil
}

func (s *Server) Shutdown() {
	s.grpcServer.GracefulStop()
	s.kvDV.Close()
}

func (s *Server) CreateTwitchClient() *helix.Client {
	twitch, err := helix.NewClient(&helix.Options{
		ClientID:     s.twitchClient,
		ClientSecret: s.twitchSecret,
		UserAgent:    "BLT AdoptAHero",
	})
	if err != nil {
		log.Fatal().Str("ctx", "server").Err(err).Send()
	}
	return twitch
}

func CreateNewServer(kvDB *db.KvDB, host string, port int, issuer string, secret []byte, twitchClient string, twitchSecret string) *Server {
	err := kvDB.Open()
	if err != nil {
		log.Fatal().Str("ctx", "server").Err(err).Send()
	}
	go kvDB.RunGC()

	return &Server{
		host:         host,
		port:         port,
		secret:       secret,
		issuer:       issuer,
		twitchClient: twitchClient,
		twitchSecret: twitchSecret,
		kvDV:         kvDB,
		channelStore: stores.NewChannelStore(kvDB),
		heroStore:    stores.NewHeroStore(kvDB),
	}
}
