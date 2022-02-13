package stores

import (
	"github.com/gorymoon/adoptahero-ebs/internal/protos"
	"github.com/gorymoon/adoptahero-ebs/pkg/db"
)

type ChannelStore struct {
	base db.Store
}

func NewChannelStore(kvDB *db.KvDB) *ChannelStore {
	return &ChannelStore{
		db.Store{
			Name: "channel",
			Db:   kvDB,
		},
	}
}

func (s *ChannelStore) GetChannel(id string) (*protos.Channel, error) {
	channel := &protos.Channel{}
	err := s.base.GetValue(id, channel)
	return channel, err
}

func (s *ChannelStore) SetChannel(id string, channel *protos.Channel) error {
	return s.base.SetValue(id, channel, 0)
}

func (s *ChannelStore) NewChannel(id string, name string, uuid string, connected bool) (*protos.Channel, error) {
	channel := &protos.Channel{
		Id:        id,
		Name:      name,
		Uuid:      uuid,
		Connected: connected,
	}
	err := s.SetChannel(id, channel)
	return channel, err
}

func (s *ChannelStore) SetKeyOnChannel(id string, uuid string) error {
	channel, err := s.GetChannel(id)
	if err != nil {
		return err
	}
	channel.Uuid = uuid
	return s.SetChannel(id, channel)
}

func (s *ChannelStore) SetConnectionOnChannel(id string, connected bool) error {
	channel, err := s.GetChannel(id)
	if err != nil {
		return err
	}
	channel.Connected = connected
	return s.SetChannel(id, channel)
}

func (s *ChannelStore) DeleteChannel(id string) error {
	err := s.base.Delete(id)
	return err
}
