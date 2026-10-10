package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"logline/internal/domain"
)

type Store struct {
	db *sql.DB
}

func New(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) Ping(ctx context.Context) error {
	return s.db.PingContext(ctx)
}

func (s *Store) InsertLog(ctx context.Context, e domain.LogEntry) error {
	data, err := json.Marshal(e.Data)
	if err != nil {
		return fmt.Errorf("marshaling entry data: %w", err)
	}

	const query = `
		INSERT INTO logs (level, message, service, timestamp, data)
		VALUES ($1, $2, $3, $4, $5)`

	if _, err := s.db.ExecContext(ctx, query,
		e.Level, e.Message, e.Service, e.Timestamp, data,
	); err != nil {
		return fmt.Errorf("inserting log: %w", err)
	}
	return nil
}
