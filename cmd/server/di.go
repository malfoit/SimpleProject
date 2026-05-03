package main

import (
	"github.com/malfoit/SimpleProject/internal/config"
	"github.com/malfoit/SimpleProject/internal/handler/user"
	githubRepo "github.com/malfoit/SimpleProject/internal/repository/github"
	userRepo "github.com/malfoit/SimpleProject/internal/repository/user"
	githubService "github.com/malfoit/SimpleProject/internal/service/github"
	syncService "github.com/malfoit/SimpleProject/internal/service/sync"
	userService "github.com/malfoit/SimpleProject/internal/service/user"
	desc "github.com/malfoit/SimpleProject/pkg/user/v1"
)

type Container struct {
	Config        *config.Config
	UserHandler   desc.UserV1Server
	GitHubService *githubService.Service
	SyncService   *syncService.Service
}

func NewContainer() *Container {
	cfg := config.NewConfig()

	uRepo := userRepo.NewRepository()
	uSvc := userService.NewService(uRepo)
	uHnd := user.NewHandler(uSvc)

	gRepo := githubRepo.NewRepo(cfg.GitHub.Token, cfg.GitHub.Owner, cfg.GitHub.Repo)
	gSvc := githubService.NewService(gRepo)

	sSvc := syncService.NewService(uRepo, uRepo) // пока оба локальные

	return &Container{
		Config:        cfg,
		UserHandler:   uHnd,
		GitHubService: gSvc,
		SyncService:   sSvc,
	}
}
