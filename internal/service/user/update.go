package user

import (
	"context"
	"errors"

	userRepo "github.com/malfoit/SimpleProject/internal/repository/user"
)

// Update обновляет имя и/или email пользователя.
//
// Шаги:
//  1. Если name != nil:
//     - обрежь пробелы (strings.TrimSpace)
//     - проверь длину: от 3 до 50 символов
//     - запиши обрезанное значение обратно в *name
//  2. Если email != nil:
//     - обрежь пробелы
//     - проверь формат через mail.ParseAddress (пакет "net/mail")
//     - запиши обрезанное значение обратно в *email
//  3. Вызови s.repo.Update(ctx, id, name, email)
//  4. Обработай ошибки репозитория:
//     - userRepo.ErrNotFound    → "user not found"
//     - userRepo.ErrAlreadyExists → "email already taken"
//     - остальные → пробрось как есть
func (s *userService) Update(ctx context.Context, id string, name, email *string) error {
	// TODO: реализуй валидацию и обновление

	_ = userRepo.ErrNotFound
	_ = userRepo.ErrAlreadyExists
	return errors.New("not implemented")
}
