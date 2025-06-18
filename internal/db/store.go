package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq" // Драйвер Postgres
)

type Store struct {
	db *sql.DB
}

func NewPostgresStore(ctx context.Context, connString string) (*Store, error) {
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	// Проверка подключения
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("db ping failed: %w", err)
	}

	return &Store{db: db}, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

// Пример метода репозитория
func (s *Store) CreateUser(ctx context.Context, name, email string) error {
	_, err := s.db.ExecContext(ctx,
		"INSERT INTO users(name, email) VALUES($1, $2)",
		name, email,
	)
	return err
}