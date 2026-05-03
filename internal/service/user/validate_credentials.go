package user

import (
	"context"

	userRepo "github.com/malfoit/SimpleProject/internal/repository/user"
)

// ValidateCredentials проверяет пару email+password.
//
// Шаги:
//  1. Найди пользователя через s.repo.GetByEmail(ctx, email)
//  2. Если userRepo.ErrNotFound — верни ("", false, nil)
//     Неверный email не должен раскрывать детали клиенту
//  3. Если другая ошибка — пробрось её
//  4. Сравни password с сохранённым хэшем через
//     bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
//     (пакет "golang.org/x/crypto/bcrypt")
//  5. Если хэши не совпадают — верни ("", false, nil)
//  6. Если совпадают — верни (u.ID, true, nil)
func (s *userService) ValidateCredentials(ctx context.Context, email, password string) (string, bool, error) {
	// TODO: реализуй метод

	_ = userRepo.ErrNotFound
	return "", false, nil
}
