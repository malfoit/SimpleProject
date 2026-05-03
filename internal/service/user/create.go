package user

import (
	"context"
	"errors"

	"github.com/malfoit/SimpleProject/internal/model"
	userRepo "github.com/malfoit/SimpleProject/internal/repository/user"
)

// Create создаёт нового пользователя после валидации и хэширования пароля.
//
// Шаги:
//  1. Обрежь пробелы у name и email (strings.TrimSpace)
//  2. Валидируй name: длина от 3 до 50 символов
//  3. Валидируй email через mail.ParseAddress из пакета "net/mail"
//  4. Валидируй password: длина от 8 до 72 символов
//  5. Проверь, что password == passwordConfirm
//  6. Захэшируй пароль через bcrypt.GenerateFromPassword (пакет "golang.org/x/crypto/bcrypt")
//     Используй bcrypt.DefaultCost
//  7. Создай *model.User с UserInfo и PasswordHash
//  8. Вызови s.repo.Create — если вернул userRepo.ErrAlreadyExists,
//     оберни в читаемую ошибку (не пропускай sentinel наружу)
//  9. Верни user.ID (репозиторий заполняет его при сохранении)
func (s *userService) Create(ctx context.Context, name, email, password, passwordConfirm string) (string, error) {
	// TODO: реализуй валидацию и создание пользователя

	_ = model.User{} // убери эту строку когда начнёшь реализацию
	_ = userRepo.ErrAlreadyExists
	return "", errors.New("not implemented")
}
