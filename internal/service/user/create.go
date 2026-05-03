package user

import (
	"context"
	"errors"
	"strings"

	"github.com/malfoit/SimpleProject/internal/model"
)

func (svc *service) Create(ctx context.Context, user model.User) (email string, err error) {

	if len(user.Name) < 3 {
		return "", errors.New("имя должно быть длинее 3 символов")
	}
	if len(user.Name) > 20 {
		return "", errors.New("имя не должно превышать 20 символов")
	}

	if strings.TrimSpace(user.Email) == "" {
		return "", errors.New("почта не может быть пустой")
	}
	if !strings.Contains(user.Email, "@") {
		return "", errors.New("неверный формат")
	}
	emailParts := strings.Split(user.Email, "@")
	if len(emailParts) != 2 || emailParts[0] == "" || emailParts[1] == "" {
		return "", errors.New("неверный формат")
	}
	if !strings.Contains(emailParts[1], ".") {
		return "", errors.New("неверный формат")
	}

	if len(user.Password) < 6 {
		return "", errors.New("пароль должен содержать не менее 6 символов")
	}
	if len(user.Password) > 50 {
		return "", errors.New("пароль не должен содержать более 50 символов")
	}

	if user.Password != user.PasswordConfirm {
		return "", errors.New("пароли не совпадают")
	}

	err = svc.repo.Create(ctx, user.Name, user.Email, user.Password)
	if err != nil {
		return "", err
	}

	return user.Email, nil
}
