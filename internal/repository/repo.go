package repository

import (
	"context"

	"github.com/malfoit/SimpleProject/internal/model"
)

type UserRepo interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	Update(ctx context.Context, id string, name, email *string) error
	UpdatePasswordHash(ctx context.Context, id, passwordHash string) error
	Delete(ctx context.Context, id string) error
}
