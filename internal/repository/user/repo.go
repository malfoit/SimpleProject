package user

import (
	"context"
	"errors"
	"sync"

	"github.com/malfoit/SimpleProject/internal/model"
	"github.com/malfoit/SimpleProject/internal/repository"
)

type repo struct {
	users map[string]*model.User
	mu    sync.RWMutex
}

func NewRepository() repository.UserRepo {
	return &repo{
		users: make(map[string]*model.User),
	}
}

func (r *repo) Create(ctx context.Context, name, email, password string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[email]; exists {
		return errors.New("user already exists")
	}

	r.users[email] = &model.User{
		Name:     name,
		Email:    email,
		Password: password,
	}
	return nil
}

// Новый метод
func (r *repo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[email]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

// Новый метод
func (r *repo) List(ctx context.Context) ([]*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]*model.User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}
	return users, nil
}
