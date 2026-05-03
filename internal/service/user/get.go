package user

import (
	"context"
	"errors"

	"github.com/malfoit/SimpleProject/internal/model"
	userRepo "github.com/malfoit/SimpleProject/internal/repository/user"
)

// Get возвращает пользователя по ID.
//
// Шаги:
//  1. Вызови s.repo.GetByID(ctx, id)
//  2. Если ошибка — проверь через errors.Is на userRepo.ErrNotFound
//     и верни понятную ошибку (не пропускай sentinel наружу)
//  3. Верни пользователя
func (s *userService) Get(ctx context.Context, id string) (*model.User, error) {
	// TODO: реализуй метод

	_ = userRepo.ErrNotFound
	return nil, errors.New("not implemented")
}
