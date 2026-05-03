package sync

import (
	"context"

	"github.com/malfoit/SimpleProject/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, name, email, password string) error
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	List(ctx context.Context) ([]*model.User, error)
}

type Service struct {
	localRepo  UserRepository // наша локальная мапа
	remoteRepo UserRepository // удалённый источник (БД, другой сервис)
}

func NewService(local, remote UserRepository) *Service {
	return &Service{
		localRepo:  local,
		remoteRepo: remote,
	}
}

// SyncUsers синхронизирует пользователей из remote в local
func (s *Service) SyncUsers(ctx context.Context) (int, error) {
	// Получаем пользователей из удалённого источника
	remoteUsers, err := s.remoteRepo.List(ctx)
	if err != nil {
		return 0, err
	}

	synced := 0
	for _, user := range remoteUsers {
		// Проверяем, есть ли уже такой пользователь
		_, err := s.localRepo.GetByEmail(ctx, user.Email)
		if err != nil {
			// Если нет - создаём
			err = s.localRepo.Create(ctx, user.Name, user.Email, user.Password)
			if err != nil {
				continue
			}
			synced++
		}
	}

	return synced, nil
}
