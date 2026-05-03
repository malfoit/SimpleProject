package user

import (
	"github.com/malfoit/SimpleProject/internal/repository"
	"github.com/malfoit/SimpleProject/internal/service"
)

type userService struct {
	repo repository.UserRepo
}

func NewService(repo repository.UserRepo) service.UserService {
	return &userService{repo: repo}
}
