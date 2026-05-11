package provider

import (
	"log/slog"

	"github.com/google/uuid"
)

type LoggingStore struct {
	name   string
	logger *slog.Logger
	next   AccumulatorStore
}

type storeSnapshotter interface {
	Snapshot() map[string]map[string]any
}

func NewLoggingStore(name string, next AccumulatorStore) AccumulatorStore {
	return &LoggingStore{
		name:   name,
		logger: slog.Default().With("component", "store", "store_name", name),
		next:   next,
	}
}

func (s *LoggingStore) CreateGroup() (uuid.UUID, error) {
	s.logger.Debug("CreateGroup.begin")
	id, err := s.next.CreateGroup()
	s.logger.Debug("CreateGroup.end", "group_id", id.String(), "err", err, "snapshot", s.snapshot())
	return id, err
}

func (s *LoggingStore) EnsureGroup() (uuid.UUID, map[string]any, error) {
	s.logger.Debug("EnsureGroup.begin")
	id, data, err := s.next.EnsureGroup()
	s.logger.Debug("EnsureGroup.end", "group_id", id.String(), "data", data, "err", err, "snapshot", s.snapshot())
	return id, data, err
}

func (s *LoggingStore) DeleteGroup(id uuid.UUID) error {
	s.logger.Debug("DeleteGroup.begin", "group_id", id.String())
	err := s.next.DeleteGroup(id)
	s.logger.Debug("DeleteGroup.end", "group_id", id.String(), "err", err, "snapshot", s.snapshot())
	return err
}

func (s *LoggingStore) PutItem(groupID uuid.UUID, key string, value any) error {
	s.logger.Debug("PutItem.begin", "group_id", groupID.String(), "key", key, "value", value)
	err := s.next.PutItem(groupID, key, value)
	s.logger.Debug("PutItem.end", "group_id", groupID.String(), "key", key, "err", err, "snapshot", s.snapshot())
	return err
}

func (s *LoggingStore) DeleteItem(groupID uuid.UUID, key string) error {
	s.logger.Debug("DeleteItem.begin", "group_id", groupID.String(), "key", key)
	err := s.next.DeleteItem(groupID, key)
	s.logger.Debug("DeleteItem.end", "group_id", groupID.String(), "key", key, "err", err, "snapshot", s.snapshot())
	return err
}

func (s *LoggingStore) GroupData(groupID uuid.UUID) (map[string]any, error) {
	s.logger.Debug("GroupData.begin", "group_id", groupID.String())
	data, err := s.next.GroupData(groupID)
	s.logger.Debug("GroupData.end", "group_id", groupID.String(), "data", data, "err", err, "snapshot", s.snapshot())
	return data, err
}

func (s *LoggingStore) snapshot() map[string]map[string]any {
	snapshotter, ok := s.next.(storeSnapshotter)
	if !ok {
		return nil
	}

	return snapshotter.Snapshot()
}
