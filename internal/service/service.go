package service

import (
	"context"

	"github.com/malfoit/SimpleProject/internal/model"
)

type UserService interface {
	Create(ctx context.Context, name, email, password, passwordConfirm string) (id string, err error)
	Get(ctx context.Context, id string) (*model.User, error)
	Update(ctx context.Context, id string, name, email *string) error
	UpdatePassword(ctx context.Context, id, password, passwordConfirm string) error
	Delete(ctx context.Context, id string) error
	ValidateCredentials(ctx context.Context, email, password string) (userID string, valid bool, err error)
}
