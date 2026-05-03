package main

import (
	"github.com/malfoit/SimpleProject/internal/config"
	userHandler "github.com/malfoit/SimpleProject/internal/handler/user"
	userRepo "github.com/malfoit/SimpleProject/internal/repository/user"
	userService "github.com/malfoit/SimpleProject/internal/service/user"
	desc "github.com/malfoit/SimpleProject/pkg/user/v1"
)

type container struct {
	config      *config.Config
	userHandler desc.UserV1Server
}

func newContainer() *container {
	cfg := config.NewConfig()

	repo := userRepo.NewRepository()
	svc := userService.NewService(repo)
	handler := userHandler.NewHandler(svc)

	return &container{
		config:      cfg,
		userHandler: handler,
	}
}
