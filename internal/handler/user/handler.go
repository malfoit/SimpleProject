package user

import (
	"context"

	"github.com/malfoit/SimpleProject/internal/model"
	"github.com/malfoit/SimpleProject/internal/service"
	desc "github.com/malfoit/SimpleProject/pkg/user/v1"
)

type handler struct {
	desc.UnimplementedUserV1Server
	userService service.UserService
}

func NewHandler(userService service.UserService) desc.UserV1Server {
	return &handler{
		userService: userService,
	}
}

func (h *handler) Create(ctx context.Context, req *desc.CreateUserRequest) (*desc.CreateUserResponse, error) {
	user := model.User{
		Name:            req.GetName(),
		Email:           req.GetEmail(),
		Password:        req.GetPassword(),
		PasswordConfirm: req.GetPasswordConfirm(),
	}

	email, err := h.userService.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return &desc.CreateUserResponse{
		Email: email,
	}, nil
}
