package user

import (
	"context"
	"errors"

	"google.golang.org/protobuf/types/known/emptypb"

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

// Create обрабатывает gRPC-запрос на создание пользователя.
//
// Шаги:
//  1. Достань UserInfo из req.GetUserInfo(), затем Name и Email через GetName/GetEmail
//  2. Вызови h.userService.Create(ctx, name, email, password, passwordConfirm)
//  3. Если ошибка — верни gRPC-статус с кодом codes.InvalidArgument
//     Используй: status.Error(codes.InvalidArgument, err.Error())
//     Пакеты: "google.golang.org/grpc/codes", "google.golang.org/grpc/status"
//  4. Верни &desc.CreateResponse{Id: id}
func (h *handler) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	// TODO: реализуй метод
	return nil, errors.New("not implemented")
}

// Get обрабатывает gRPC-запрос на получение пользователя по ID.
//
// Шаги:
//  1. Вызови h.userService.Get(ctx, req.GetId())
//  2. Если ошибка — верни статус codes.NotFound
//  3. Собери ответ &desc.GetResponse{User: &desc.User{...}}:
//     - Id, UserInfo (Name + Email)
//     - CreatedAt и UpdatedAt через timestamppb.New(u.CreatedAt)
//     Пакет: "google.golang.org/protobuf/types/known/timestamppb"
func (h *handler) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	// TODO: реализуй метод
	return nil, errors.New("not implemented")
}

// Update обрабатывает gRPC-запрос на обновление имени и/или email.
//
// Шаги:
//  1. Name и Email в proto — это *wrapperspb.StringValue (опциональные поля).
//     Если req.GetName() != nil — извлеки строку через .GetValue() и возьми указатель на неё.
//     Если nil — передай nil в сервис.
//     Аналогично для Email.
//  2. Вызови h.userService.Update(ctx, req.GetId(), name, email)
//  3. Если ошибка — верни статус codes.InvalidArgument
//  4. Верни &emptypb.Empty{}, nil
func (h *handler) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	// TODO: реализуй метод
	return nil, errors.New("not implemented")
}

// UpdatePassword обрабатывает gRPC-запрос на смену пароля.
//
// Шаги:
//  1. Вызови h.userService.UpdatePassword(ctx, id, password, passwordConfirm)
//  2. Если ошибка — верни статус codes.InvalidArgument
//  3. Верни &emptypb.Empty{}, nil
func (h *handler) UpdatePassword(ctx context.Context, req *desc.UpdatePasswordRequest) (*emptypb.Empty, error) {
	// TODO: реализуй метод
	return nil, errors.New("not implemented")
}

// Delete обрабатывает gRPC-запрос на удаление пользователя.
//
// Шаги:
//  1. Вызови h.userService.Delete(ctx, req.GetId())
//  2. Если ошибка — верни статус codes.NotFound
//  3. Верни &emptypb.Empty{}, nil
func (h *handler) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	// TODO: реализуй метод
	return nil, errors.New("not implemented")
}

// ValidateCredentials обрабатывает gRPC-запрос на проверку учётных данных.
//
// Шаги:
//  1. Вызови h.userService.ValidateCredentials(ctx, email, password)
//  2. Если ошибка — верни статус codes.Internal с ОБЩИМ сообщением
//     (не раскрывай детали внутренней ошибки клиенту!)
//  3. Верни &desc.ValidateCredentialsResponse{Valid: valid, UserId: userID}
func (h *handler) ValidateCredentials(ctx context.Context, req *desc.ValidateCredentialsRequest) (*desc.ValidateCredentialsResponse, error) {
	// TODO: реализуй метод
	return nil, errors.New("not implemented")
}
