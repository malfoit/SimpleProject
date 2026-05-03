package repository

import (
	"context"

	"github.com/malfoit/SimpleProject/internal/model"
)

type UserRepo interface {
	Create(ctx context.Context, name, email string, password string) error
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	List(ctx context.Context) ([]*model.User, error)
}
