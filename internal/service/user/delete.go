package user

import (
	"context"
	"errors"

	userRepo "github.com/malfoit/SimpleProject/internal/repository/user"
)

// Delete удаляет пользователя по ID.
//
// Шаги:
//  1. Вызови s.repo.Delete(ctx, id)
//  2. Если userRepo.ErrNotFound — верни читаемую ошибку
//  3. Остальные ошибки пробрось как есть
func (s *userService) Delete(ctx context.Context, id string) error {
	// TODO: реализуй метод

	_ = userRepo.ErrNotFound
	return errors.New("not implemented")
}
