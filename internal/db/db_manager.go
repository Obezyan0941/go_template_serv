package db_manager

import (
	"context"
	"fmt"

	"github.com/go-pg/pg/v10"
)

func NewPostgresConnection(db_data DBConfig) (*pg.DB, error) {
	db := pg.Connect(&pg.Options{
		Addr:     db_data.Addr,
		User:     db_data.User,
		Password: db_data.Password,
		Database: db_data.Database,
	})

	ctx := context.Background()
	if err := db.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}