package provider

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
)

type AccumulatorStore interface {
	CreateGroup() (uuid.UUID, error)
	DeleteGroup(id uuid.UUID) error
	PutItem(groupID uuid.UUID, key string, value any) error
	DeleteItem(groupID uuid.UUID, key string) error
	GroupData(groupID uuid.UUID) (map[string]any, error)
}

type MemoryStore struct {
	mu     sync.RWMutex
	groups map[uuid.UUID]map[string]any
}

func NewStore() AccumulatorStore {
	return &MemoryStore{
		groups: map[uuid.UUID]map[string]any{},
	}
}

func (s *MemoryStore) CreateGroup() (uuid.UUID, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := uuid.New()
	s.groups[id] = map[string]any{}

	return id, nil
}

func (s *MemoryStore) DeleteGroup(id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.groups[id]; !ok {
		return fmt.Errorf("group %s not found", id)
	}

	delete(s.groups, id)
	return nil
}

func (s *MemoryStore) PutItem(groupID uuid.UUID, key string, value any) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	group, ok := s.groups[groupID]
	if !ok {
		return fmt.Errorf("group %s not found", groupID)
	}

	if _, exists := group[key]; exists {
		return fmt.Errorf("duplicate key %q in group %s", key, groupID)
	}

	group[key] = value
	return nil
}

func (s *MemoryStore) DeleteItem(groupID uuid.UUID, key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	group, ok := s.groups[groupID]
	if !ok {
		return fmt.Errorf("group %s not found", groupID)
	}

	delete(group, key)
	return nil
}

func (s *MemoryStore) GroupData(groupID uuid.UUID) (map[string]any, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	group, ok := s.groups[groupID]
	if !ok {
		return nil, fmt.Errorf("group %s not found", groupID)
	}

	out := make(map[string]any, len(group))
	for key, value := range group {
		out[key] = value
	}

	return out, nil
}
