package db

import (
	"context"
	"test_backend/internal/models"
)

type UserRepository interface {
	CreateUser(ctx context.Context, name, email string) (*models.User, error)
	GetUserByID(ctx context.Context, id int) (*models.User, error)
}