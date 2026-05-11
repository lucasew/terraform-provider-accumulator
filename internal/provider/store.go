package provider

import "github.com/google/uuid"

type AccumulatorStore interface {
	CreateGroup(name string, typeName string) (uuid.UUID, error)
	DeleteGroup(id uuid.UUID) error
	PutItem(groupID uuid.UUID, key string, value any) error
	DeleteItem(groupID uuid.UUID, key string) error
	GroupData(groupID uuid.UUID) (map[string]any, error)
}

type MemoryStore struct{}

func NewStore() AccumulatorStore {
	return &MemoryStore{}
}

func (s *MemoryStore) CreateGroup(name string, typeName string) (uuid.UUID, error) {
	panic("unimplemented")
}

func (s *MemoryStore) DeleteGroup(id uuid.UUID) error {
	panic("unimplemented")
}

func (s *MemoryStore) PutItem(groupID uuid.UUID, key string, value any) error {
	panic("unimplemented")
}

func (s *MemoryStore) DeleteItem(groupID uuid.UUID, key string) error {
	panic("unimplemented")
}

func (s *MemoryStore) GroupData(groupID uuid.UUID) (map[string]any, error) {
	panic("unimplemented")
}
