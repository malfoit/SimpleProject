package user

import (
	"context"
	"errors"

	userRepo "github.com/malfoit/SimpleProject/internal/repository/user"
)

// UpdatePassword меняет пароль пользователя.
//
// Шаги:
//  1. Проверь длину password: от 8 до 72 символов
//     (72 — максимум, который обрабатывает bcrypt)
//  2. Проверь, что password == passwordConfirm
//  3. Захэшируй пароль через bcrypt.GenerateFromPassword
//     (пакет "golang.org/x/crypto/bcrypt", bcrypt.DefaultCost)
//  4. Вызови s.repo.UpdatePasswordHash(ctx, id, string(hash))
//  5. Если userRepo.ErrNotFound — верни читаемую ошибку
func (s *userService) UpdatePassword(ctx context.Context, id, password, passwordConfirm string) error {
	// TODO: реализуй валидацию и смену пароля

	_ = userRepo.ErrNotFound
	return errors.New("not implemented")
}
