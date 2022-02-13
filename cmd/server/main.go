package main

import (
	"encoding/base64"
	"github.com/gorymoon/adoptahero-ebs/internal/server"
	"github.com/gorymoon/adoptahero-ebs/pkg/db"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func main() {
	level, err := zerolog.ParseLevel(GetEnv("LOG_LEVEL", "info"))
	if err != nil {
		log.Fatal().Str("ctx", "main").Err(err).Msg("Invalid log level")
	}
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(level)

	port, err := strconv.Atoi(GetEnv("PORT", "50051"))
	if err != nil {
		log.Fatal().Str("ctx", "main").Err(err).Msg("PORT must be a number")
	}

	secretString := GetEnv("SECRET", "")
	if len(secretString) <= 0 {
		log.Fatal().Str("ctx", "main").Msg("SECRET must be set")
	}

	secret, err := base64.StdEncoding.DecodeString(secretString)
	if err != nil {
		log.Fatal().Str("ctx", "main").Err(err).Msg("SECRET must be a base64 encoded value")
	}

	kvDB := &db.KvDB{
		DBPath: GetEnv("DB_LOC", "run/db"),
	}
	srv := server.CreateNewServer(kvDB, GetEnv("HOST", ""), port, secret)
	defer srv.Shutdown()

	// Setups a signal to listen for termination
	go func() {
		signChan := make(chan os.Signal, 1)

		signal.Notify(signChan, os.Interrupt, syscall.SIGTERM)
		sig := <-signChan
		log.Info().Str("ctx", "main").Str("sig", sig.String()).Msg("Shutting down")

		srv.Shutdown()
	}()

	if err := srv.Run(); err != nil {
		log.Fatal().Str("ctx", "main").Err(err).Send()
	}
	log.Info().Str("ctx", "main").Msg("Goodbye!")
}

func GetEnv(key string, fallback string) string {
	if env, ok := os.LookupEnv(key); ok {
		return env
	}
	return fallback
}
