package user_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"

	handleruser "github.com/malfoit/SimpleProject/internal/handler/user"
	"github.com/malfoit/SimpleProject/internal/model"
	desc "github.com/malfoit/SimpleProject/pkg/user/v1"
)

// ---------------------------------------------------------------------------
// Mock service
// ---------------------------------------------------------------------------

type mockService struct {
	createFn             func(ctx context.Context, name, email, password, passwordConfirm string) (string, error)
	getFn                func(ctx context.Context, id string) (*model.User, error)
	updateFn             func(ctx context.Context, id string, name, email *string) error
	updatePasswordFn     func(ctx context.Context, id, password, passwordConfirm string) error
	deleteFn             func(ctx context.Context, id string) error
	validateCredentialsFn func(ctx context.Context, email, password string) (string, bool, error)
}

func (m *mockService) Create(ctx context.Context, name, email, password, passwordConfirm string) (string, error) {
	if m.createFn != nil {
		return m.createFn(ctx, name, email, password, passwordConfirm)
	}
	return "", nil
}
func (m *mockService) Get(ctx context.Context, id string) (*model.User, error) {
	if m.getFn != nil {
		return m.getFn(ctx, id)
	}
	return nil, nil
}
func (m *mockService) Update(ctx context.Context, id string, name, email *string) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, id, name, email)
	}
	return nil
}
func (m *mockService) UpdatePassword(ctx context.Context, id, password, passwordConfirm string) error {
	if m.updatePasswordFn != nil {
		return m.updatePasswordFn(ctx, id, password, passwordConfirm)
	}
	return nil
}
func (m *mockService) Delete(ctx context.Context, id string) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, id)
	}
	return nil
}
func (m *mockService) ValidateCredentials(ctx context.Context, email, password string) (string, bool, error) {
	if m.validateCredentialsFn != nil {
		return m.validateCredentialsFn(ctx, email, password)
	}
	return "", false, nil
}

// grpcCode извлекает gRPC код из ошибки.
func grpcCode(err error) codes.Code {
	if s, ok := status.FromError(err); ok {
		return s.Code()
	}
	return codes.Unknown
}

// ---------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------

func TestCreate_PassesAllFieldsToService(t *testing.T) {
	var gotName, gotEmail, gotPass, gotConfirm string

	svc := &mockService{
		createFn: func(_ context.Context, name, email, password, passwordConfirm string) (string, error) {
			gotName, gotEmail, gotPass, gotConfirm = name, email, password, passwordConfirm
			return "uid-1", nil
		},
	}
	h := handleruser.NewHandler(svc)

	resp, err := h.Create(context.Background(), &desc.CreateRequest{
		UserInfo:        &desc.UserInfo{Name: "Alice", Email: "alice@example.com"},
		Password:        "password123",
		PasswordConfirm: "password123",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotName != "Alice" {
		t.Errorf("name: want 'Alice', got %q", gotName)
	}
	if gotEmail != "alice@example.com" {
		t.Errorf("email: want 'alice@example.com', got %q", gotEmail)
	}
	if gotPass != "password123" {
		t.Errorf("password not passed correctly")
	}
	if gotConfirm != "password123" {
		t.Errorf("passwordConfirm not passed correctly")
	}
	if resp.GetId() != "uid-1" {
		t.Errorf("id: want 'uid-1', got %q", resp.GetId())
	}
}

func TestCreate_ServiceError_ReturnsInvalidArgument(t *testing.T) {
	svc := &mockService{
		createFn: func(_ context.Context, _, _, _, _ string) (string, error) {
			return "", errors.New("name too short")
		},
	}
	h := handleruser.NewHandler(svc)

	_, err := h.Create(context.Background(), &desc.CreateRequest{
		UserInfo: &desc.UserInfo{Name: "Al", Email: "al@example.com"},
		Password: "pass", PasswordConfirm: "pass",
	})

	if err == nil {
		t.Fatal("expected error")
	}
	if grpcCode(err) != codes.InvalidArgument {
		t.Errorf("want codes.InvalidArgument, got %v", grpcCode(err))
	}
}

func TestCreate_NilUserInfo_DoesNotPanic(t *testing.T) {
	svc := &mockService{
		createFn: func(_ context.Context, name, email, _, _ string) (string, error) {
			if name != "" || email != "" {
				t.Errorf("expected empty name/email for nil UserInfo, got name=%q email=%q", name, email)
			}
			return "id", nil
		},
	}
	h := handleruser.NewHandler(svc)

	// nil UserInfo — proto GetUserInfo() returns nil, GetName/GetEmail return ""
	_, err := h.Create(context.Background(), &desc.CreateRequest{
		UserInfo: nil, Password: "password123", PasswordConfirm: "password123",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// ---------------------------------------------------------------------------
// Get
// ---------------------------------------------------------------------------

func TestGet_PassesIDToService(t *testing.T) {
	var gotID string
	now := time.Now().Truncate(time.Second)

	svc := &mockService{
		getFn: func(_ context.Context, id string) (*model.User, error) {
			gotID = id
			return &model.User{
				ID:        "uid-1",
				UserInfo:  model.UserInfo{Name: "Alice", Email: "alice@example.com"},
				CreatedAt: now,
				UpdatedAt: now,
			}, nil
		},
	}
	h := handleruser.NewHandler(svc)

	resp, err := h.Get(context.Background(), &desc.GetRequest{Id: "uid-1"})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotID != "uid-1" {
		t.Errorf("id passed to service: want 'uid-1', got %q", gotID)
	}
	if resp.GetUser().GetId() != "uid-1" {
		t.Errorf("response id: want 'uid-1', got %q", resp.GetUser().GetId())
	}
}

func TestGet_MapsAllUserFieldsToResponse(t *testing.T) {
	now := time.Date(2025, 6, 1, 12, 0, 0, 0, time.UTC)

	svc := &mockService{
		getFn: func(_ context.Context, _ string) (*model.User, error) {
			return &model.User{
				ID:        "uid-42",
				UserInfo:  model.UserInfo{Name: "Bob", Email: "bob@example.com"},
				CreatedAt: now,
				UpdatedAt: now.Add(time.Hour),
			}, nil
		},
	}
	h := handleruser.NewHandler(svc)

	resp, _ := h.Get(context.Background(), &desc.GetRequest{Id: "uid-42"})
	u := resp.GetUser()

	if u.GetId() != "uid-42" {
		t.Errorf("id: want 'uid-42', got %q", u.GetId())
	}
	if u.GetUserInfo().GetName() != "Bob" {
		t.Errorf("name: want 'Bob', got %q", u.GetUserInfo().GetName())
	}
	if u.GetUserInfo().GetEmail() != "bob@example.com" {
		t.Errorf("email: want 'bob@example.com', got %q", u.GetUserInfo().GetEmail())
	}
	if !u.GetCreatedAt().AsTime().Equal(now) {
		t.Errorf("created_at mismatch")
	}
	if !u.GetUpdatedAt().AsTime().Equal(now.Add(time.Hour)) {
		t.Errorf("updated_at mismatch")
	}
}

func TestGet_ServiceError_ReturnsNotFound(t *testing.T) {
	svc := &mockService{
		getFn: func(_ context.Context, _ string) (*model.User, error) {
			return nil, errors.New("user not found")
		},
	}
	h := handleruser.NewHandler(svc)

	_, err := h.Get(context.Background(), &desc.GetRequest{Id: "ghost"})

	if err == nil {
		t.Fatal("expected error")
	}
	if grpcCode(err) != codes.NotFound {
		t.Errorf("want codes.NotFound, got %v", grpcCode(err))
	}
}

// ---------------------------------------------------------------------------
// Update
// ---------------------------------------------------------------------------

func TestUpdate_BothFieldsSet_PassedAsPointers(t *testing.T) {
	var gotID string
	var gotName, gotEmail *string

	svc := &mockService{
		updateFn: func(_ context.Context, id string, name, email *string) error {
			gotID = id
			gotName = name
			gotEmail = email
			return nil
		},
	}
	h := handleruser.NewHandler(svc)

	_, err := h.Update(context.Background(), &desc.UpdateRequest{
		Id:    "uid-1",
		Name:  wrapperspb.String("NewName"),
		Email: wrapperspb.String("new@example.com"),
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotID != "uid-1" {
		t.Errorf("id: want 'uid-1', got %q", gotID)
	}
	if gotName == nil || *gotName != "NewName" {
		t.Errorf("name: want 'NewName', got %v", gotName)
	}
	if gotEmail == nil || *gotEmail != "new@example.com" {
		t.Errorf("email: want 'new@example.com', got %v", gotEmail)
	}
}

func TestUpdate_OnlyNameSet_EmailIsNil(t *testing.T) {
	var gotName, gotEmail *string

	svc := &mockService{
		updateFn: func(_ context.Context, _ string, name, email *string) error {
			gotName = name
			gotEmail = email
			return nil
		},
	}
	h := handleruser.NewHandler(svc)

	_, _ = h.Update(context.Background(), &desc.UpdateRequest{
		Id:   "uid-1",
		Name: wrapperspb.String("OnlyName"),
	})

	if gotName == nil || *gotName != "OnlyName" {
		t.Errorf("expected name='OnlyName', got %v", gotName)
	}
	if gotEmail != nil {
		t.Errorf("expected email=nil, got %v", gotEmail)
	}
}

func TestUpdate_OnlyEmailSet_NameIsNil(t *testing.T) {
	var gotName, gotEmail *string

	svc := &mockService{
		updateFn: func(_ context.Context, _ string, name, email *string) error {
			gotName = name
			gotEmail = email
			return nil
		},
	}
	h := handleruser.NewHandler(svc)

	_, _ = h.Update(context.Background(), &desc.UpdateRequest{
		Id:    "uid-1",
		Email: wrapperspb.String("only@example.com"),
	})

	if gotName != nil {
		t.Errorf("expected name=nil, got %v", gotName)
	}
	if gotEmail == nil || *gotEmail != "only@example.com" {
		t.Errorf("expected email='only@example.com', got %v", gotEmail)
	}
}

func TestUpdate_BothFieldsAbsent_BothNilPassedToService(t *testing.T) {
	var gotName, gotEmail *string

	svc := &mockService{
		updateFn: func(_ context.Context, _ string, name, email *string) error {
			gotName = name
			gotEmail = email
			return nil
		},
	}
	h := handleruser.NewHandler(svc)

	_, _ = h.Update(context.Background(), &desc.UpdateRequest{Id: "uid-1"})

	if gotName != nil || gotEmail != nil {
		t.Errorf("expected both nil, got name=%v email=%v", gotName, gotEmail)
	}
}

func TestUpdate_ServiceError_ReturnsInvalidArgument(t *testing.T) {
	svc := &mockService{
		updateFn: func(_ context.Context, _ string, _, _ *string) error {
			return errors.New("email already taken")
		},
	}
	h := handleruser.NewHandler(svc)

	_, err := h.Update(context.Background(), &desc.UpdateRequest{
		Id:    "uid-1",
		Email: wrapperspb.String("taken@example.com"),
	})

	if err == nil {
		t.Fatal("expected error")
	}
	if grpcCode(err) != codes.InvalidArgument {
		t.Errorf("want codes.InvalidArgument, got %v", grpcCode(err))
	}
}

// ---------------------------------------------------------------------------
// UpdatePassword
// ---------------------------------------------------------------------------

func TestUpdatePassword_PassesAllFieldsToService(t *testing.T) {
	var gotID, gotPass, gotConfirm string

	svc := &mockService{
		updatePasswordFn: func(_ context.Context, id, password, passwordConfirm string) error {
			gotID, gotPass, gotConfirm = id, password, passwordConfirm
			return nil
		},
	}
	h := handleruser.NewHandler(svc)

	_, err := h.UpdatePassword(context.Background(), &desc.UpdatePasswordRequest{
		Id:              "uid-1",
		Password:        "newpassword1",
		PasswordConfirm: "newpassword1",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotID != "uid-1" {
		t.Errorf("id: want 'uid-1', got %q", gotID)
	}
	if gotPass != "newpassword1" {
		t.Errorf("password not passed correctly")
	}
	if gotConfirm != "newpassword1" {
		t.Errorf("passwordConfirm not passed correctly")
	}
}

func TestUpdatePassword_ServiceError_ReturnsInvalidArgument(t *testing.T) {
	svc := &mockService{
		updatePasswordFn: func(_ context.Context, _, _, _ string) error {
			return errors.New("passwords do not match")
		},
	}
	h := handleruser.NewHandler(svc)

	_, err := h.UpdatePassword(context.Background(), &desc.UpdatePasswordRequest{
		Id: "uid-1", Password: "aaa", PasswordConfirm: "bbb",
	})

	if err == nil {
		t.Fatal("expected error")
	}
	if grpcCode(err) != codes.InvalidArgument {
		t.Errorf("want codes.InvalidArgument, got %v", grpcCode(err))
	}
}

// ---------------------------------------------------------------------------
// Delete
// ---------------------------------------------------------------------------

func TestDelete_PassesIDToService(t *testing.T) {
	var gotID string

	svc := &mockService{
		deleteFn: func(_ context.Context, id string) error {
			gotID = id
			return nil
		},
	}
	h := handleruser.NewHandler(svc)

	_, err := h.Delete(context.Background(), &desc.DeleteRequest{Id: "uid-1"})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotID != "uid-1" {
		t.Errorf("id: want 'uid-1', got %q", gotID)
	}
}

func TestDelete_ServiceError_ReturnsNotFound(t *testing.T) {
	svc := &mockService{
		deleteFn: func(_ context.Context, _ string) error {
			return errors.New("user not found")
		},
	}
	h := handleruser.NewHandler(svc)

	_, err := h.Delete(context.Background(), &desc.DeleteRequest{Id: "ghost"})

	if err == nil {
		t.Fatal("expected error")
	}
	if grpcCode(err) != codes.NotFound {
		t.Errorf("want codes.NotFound, got %v", grpcCode(err))
	}
}

// ---------------------------------------------------------------------------
// ValidateCredentials
// ---------------------------------------------------------------------------

func TestValidateCredentials_PassesEmailAndPasswordToService(t *testing.T) {
	var gotEmail, gotPass string

	svc := &mockService{
		validateCredentialsFn: func(_ context.Context, email, password string) (string, bool, error) {
			gotEmail, gotPass = email, password
			return "uid-1", true, nil
		},
	}
	h := handleruser.NewHandler(svc)

	resp, err := h.ValidateCredentials(context.Background(), &desc.ValidateCredentialsRequest{
		Email:    "alice@example.com",
		Password: "password123",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotEmail != "alice@example.com" {
		t.Errorf("email: want 'alice@example.com', got %q", gotEmail)
	}
	if gotPass != "password123" {
		t.Errorf("password not passed correctly")
	}
	if !resp.GetValid() {
		t.Error("want valid=true")
	}
	if resp.GetUserId() != "uid-1" {
		t.Errorf("user_id: want 'uid-1', got %q", resp.GetUserId())
	}
}

func TestValidateCredentials_InvalidCredentials_ValidFalseNoError(t *testing.T) {
	svc := &mockService{
		validateCredentialsFn: func(_ context.Context, _, _ string) (string, bool, error) {
			return "", false, nil
		},
	}
	h := handleruser.NewHandler(svc)

	resp, err := h.ValidateCredentials(context.Background(), &desc.ValidateCredentialsRequest{
		Email: "alice@example.com", Password: "wrongpass",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.GetValid() {
		t.Error("want valid=false")
	}
	if resp.GetUserId() != "" {
		t.Errorf("want empty user_id, got %q", resp.GetUserId())
	}
}

func TestValidateCredentials_ServiceError_ReturnsInternal_MessageHidden(t *testing.T) {
	svc := &mockService{
		validateCredentialsFn: func(_ context.Context, _, _ string) (string, bool, error) {
			return "", false, errors.New("db connection lost — sensitive detail")
		},
	}
	h := handleruser.NewHandler(svc)

	_, err := h.ValidateCredentials(context.Background(), &desc.ValidateCredentialsRequest{
		Email: "alice@example.com", Password: "pass",
	})

	if err == nil {
		t.Fatal("expected error")
	}
	if grpcCode(err) != codes.Internal {
		t.Errorf("want codes.Internal, got %v", grpcCode(err))
	}
	// Внутренняя ошибка не должна протекать клиенту
	if s, _ := status.FromError(err); s.Message() == "db connection lost — sensitive detail" {
		t.Error("internal error detail must not be exposed to the client")
	}
}
