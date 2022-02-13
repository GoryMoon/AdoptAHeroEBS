package db

import (
	"errors"
	"github.com/dgraph-io/badger/v3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"strings"
	"time"
)

type KvDB struct {
	db     *badger.DB
	opened bool
	DBPath string
}

func (b *KvDB) Open() error {
	dbOpts := badger.DefaultOptions(b.DBPath)
	dbOpts.ValueLogFileSize = 128 << 20 // 128MB
	dbOpts.IndexCacheSize = 128 << 20   // 128MB
	dbOpts.BaseTableSize = 8 << 20      // 8MB
	dbOpts.CompactL0OnClose = true
	dbOpts.Logger = &ZeroLogger{Logger: log.Logger}

	db, err := badger.Open(dbOpts)
	if err != nil {
		return err
	}
	b.db = db
	b.opened = true
	return nil
}

func (b KvDB) GetDB() (*badger.DB, error) {
	if b.opened {
		return b.db, nil
	}
	return nil, errors.New("database isn't open, you need to open it before calling this")
}

func (b *KvDB) Close() {
	b.db.Close()
	b.opened = false
}

func (b *KvDB) RunGC() {
	for {
		time.Sleep(10 * time.Minute)
		if err := b.db.RunValueLogGC(0.7); err != nil {
			if err != badger.ErrNoRewrite {
				log.Fatal().Err(err).Send()
			}
		}
	}
}

type ZeroLogger struct {
	zerolog.Logger
}

func (z *ZeroLogger) Errorf(s string, i ...interface{}) {
	z.Error().Str("ctx", "badger").Msgf(strings.TrimSuffix(s, "\n"), i...)
}

func (z *ZeroLogger) Warningf(s string, i ...interface{}) {
	z.Warn().Str("ctx", "badger").Msgf(strings.TrimSuffix(s, "\n"), i...)
}

func (z *ZeroLogger) Infof(s string, i ...interface{}) {
	z.Info().Str("ctx", "badger").Msgf(strings.TrimSuffix(s, "\n"), i...)
}

func (z *ZeroLogger) Debugf(s string, i ...interface{}) {
	z.Debug().Str("ctx", "badger").Msgf(strings.TrimSuffix(s, "\n"), i...)
}
