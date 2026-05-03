package user

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/malfoit/SimpleProject/internal/service"
	desc "github.com/malfoit/SimpleProject/pkg/user/v1"
)

type handler struct {
	desc.UnimplementedUserV1Server
	userService service.UserService
}

func NewHandler(userService service.UserService) desc.UserV1Server {
	return &handler{userService: userService}
}

func (h *handler) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	info := req.GetUserInfo()
	id, err := h.userService.Create(ctx,
		info.GetName(),
		info.GetEmail(),
		req.GetPassword(),
		req.GetPasswordConfirm(),
	)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &desc.CreateResponse{Id: id}, nil
}

func (h *handler) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	u, err := h.userService.Get(ctx, req.GetId())
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return &desc.GetResponse{
		User: &desc.User{
			Id: u.ID,
			UserInfo: &desc.UserInfo{
				Name:  u.UserInfo.Name,
				Email: u.UserInfo.Email,
			},
			CreatedAt: timestamppb.New(u.CreatedAt),
			UpdatedAt: timestamppb.New(u.UpdatedAt),
		},
	}, nil
}

func (h *handler) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	var name, email *string
	if v := req.GetName(); v != nil {
		s := v.GetValue()
		name = &s
	}
	if v := req.GetEmail(); v != nil {
		s := v.GetValue()
		email = &s
	}
	if err := h.userService.Update(ctx, req.GetId(), name, email); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (h *handler) UpdatePassword(ctx context.Context, req *desc.UpdatePasswordRequest) (*emptypb.Empty, error) {
	err := h.userService.UpdatePassword(ctx, req.GetId(), req.GetPassword(), req.GetPasswordConfirm())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (h *handler) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	if err := h.userService.Delete(ctx, req.GetId()); err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (h *handler) ValidateCredentials(ctx context.Context, req *desc.ValidateCredentialsRequest) (*desc.ValidateCredentialsResponse, error) {
	userID, valid, err := h.userService.ValidateCredentials(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, status.Error(codes.Internal, "credential validation failed")
	}
	return &desc.ValidateCredentialsResponse{
		Valid:  valid,
		UserId: userID,
	}, nil
}
