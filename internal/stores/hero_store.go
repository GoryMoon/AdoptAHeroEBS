package stores

import (
	"github.com/gorymoon/adoptahero-ebs/internal/protos"
	"github.com/gorymoon/adoptahero-ebs/pkg/db"
	"google.golang.org/protobuf/proto"
	"time"
)

type HeroStore struct {
	base db.Store
}

func NewHeroStore(kvDB *db.KvDB) *HeroStore {
	return &HeroStore{
		base: db.Store{
			Name: "hero",
			Db:   kvDB,
		},
	}
}

func (s *HeroStore) GetHero(id string) (*protos.HeroData, error) {
	channel := &protos.HeroData{}
	err := s.base.GetValue(id, channel)
	return channel, err
}

func (s *HeroStore) SetHero(id string, hero *protos.HeroData) error {
	return s.base.SetValue(id, hero, time.Hour)
}

func (s *HeroStore) SetHeroes(heroes map[string]*protos.HeroData) error {
	keys := make([]string, 0, len(heroes))
	for k := range heroes {
		keys = append(keys, k)
	}
	return s.base.BulkSetValue(keys, func(key string) proto.Message {
		return heroes[key]
	}, time.Hour)
}

func (s *HeroStore) DeleteHero(id string) error {
	err := s.base.Delete(id)
	return err
}

func (s *HeroStore) DeleteHeroes(ids []string) error {
	err := s.base.BulkDelete(ids)
	return err
}
