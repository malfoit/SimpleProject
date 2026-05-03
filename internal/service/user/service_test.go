package user_test

import (
	"context"
	"testing"

	userRepo "github.com/malfoit/SimpleProject/internal/repository/user"
	"github.com/malfoit/SimpleProject/internal/service/user"
)

func TestCreate_Success(t *testing.T) {
	repo := userRepo.NewRepository()
	svc := user.NewService(repo)
	ctx := context.Background()

	id, err := svc.Create(ctx, "Alice", "alice@example.com", "password123", "password123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id == "" {
		t.Fatal("expected non-empty id")
	}
}

func TestCreate_DuplicateEmail(t *testing.T) {
	repo := userRepo.NewRepository()
	svc := user.NewService(repo)
	ctx := context.Background()

	_, err := svc.Create(ctx, "Alice", "alice@example.com", "password123", "password123")
	if err != nil {
		t.Fatalf("first create failed: %v", err)
	}

	_, err = svc.Create(ctx, "Bob", "alice@example.com", "password123", "password123")
	if err == nil {
		t.Fatal("expected error on duplicate email")
	}
}

func TestCreate_PasswordMismatch(t *testing.T) {
	repo := userRepo.NewRepository()
	svc := user.NewService(repo)
	ctx := context.Background()

	_, err := svc.Create(ctx, "Alice", "alice@example.com", "password123", "different")
	if err == nil {
		t.Fatal("expected error on password mismatch")
	}
}

func TestCreate_ShortPassword(t *testing.T) {
	repo := userRepo.NewRepository()
	svc := user.NewService(repo)
	ctx := context.Background()

	_, err := svc.Create(ctx, "Alice", "alice@example.com", "pass", "pass")
	if err == nil {
		t.Fatal("expected error on short password")
	}
}

func TestCreate_InvalidEmail(t *testing.T) {
	repo := userRepo.NewRepository()
	svc := user.NewService(repo)
	ctx := context.Background()

	_, err := svc.Create(ctx, "Alice", "not-an-email", "password123", "password123")
	if err == nil {
		t.Fatal("expected error on invalid email")
	}
}

func TestCreate_ShortName(t *testing.T) {
	repo := userRepo.NewRepository()
	svc := user.NewService(repo)
	ctx := context.Background()

	_, err := svc.Create(ctx, "Al", "alice@example.com", "password123", "password123")
	if err == nil {
		t.Fatal("expected error on short name")
	}
}

func TestGet_NotFound(t *testing.T) {
	repo := userRepo.NewRepository()
	svc := user.NewService(repo)
	ctx := context.Background()

	_, err := svc.Get(ctx, "nonexistent-id")
	if err == nil {
		t.Fatal("expected error on missing user")
	}
}

func TestGet_Success(t *testing.T) {
	repo := userRepo.NewRepository()
	svc := user.NewService(repo)
	ctx := context.Background()

	id, _ := svc.Create(ctx, "Alice", "alice@example.com", "password123", "password123")

	u, err := svc.Get(ctx, id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if u.UserInfo.Email != "alice@example.com" {
		t.Errorf("unexpected email: %s", u.UserInfo.Email)
	}
}

func TestUpdate_Success(t *testing.T) {
	repo := userRepo.NewRepository()
	svc := user.NewService(repo)
	ctx := context.Background()

	id, _ := svc.Create(ctx, "Alice", "alice@example.com", "password123", "password123")

	newName := "Alice Updated"
	if err := svc.Update(ctx, id, &newName, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	u, _ := svc.Get(ctx, id)
	if u.UserInfo.Name != newName {
		t.Errorf("name not updated: got %s", u.UserInfo.Name)
	}
}

func TestDelete_Success(t *testing.T) {
	repo := userRepo.NewRepository()
	svc := user.NewService(repo)
	ctx := context.Background()

	id, _ := svc.Create(ctx, "Alice", "alice@example.com", "password123", "password123")

	if err := svc.Delete(ctx, id); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err := svc.Get(ctx, id)
	if err == nil {
		t.Fatal("expected error after deletion")
	}
}

func TestValidateCredentials_Success(t *testing.T) {
	repo := userRepo.NewRepository()
	svc := user.NewService(repo)
	ctx := context.Background()

	id, _ := svc.Create(ctx, "Alice", "alice@example.com", "password123", "password123")

	gotID, valid, err := svc.ValidateCredentials(ctx, "alice@example.com", "password123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !valid {
		t.Fatal("expected valid=true")
	}
	if gotID != id {
		t.Errorf("expected id %s, got %s", id, gotID)
	}
}

func TestValidateCredentials_WrongPassword(t *testing.T) {
	repo := userRepo.NewRepository()
	svc := user.NewService(repo)
	ctx := context.Background()

	_, _ = svc.Create(ctx, "Alice", "alice@example.com", "password123", "password123")

	_, valid, err := svc.ValidateCredentials(ctx, "alice@example.com", "wrongpassword")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if valid {
		t.Fatal("expected valid=false for wrong password")
	}
}

func TestValidateCredentials_UnknownEmail(t *testing.T) {
	repo := userRepo.NewRepository()
	svc := user.NewService(repo)
	ctx := context.Background()

	_, valid, err := svc.ValidateCredentials(ctx, "nobody@example.com", "password123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if valid {
		t.Fatal("expected valid=false for unknown email")
	}
}

func TestUpdatePassword_Success(t *testing.T) {
	repo := userRepo.NewRepository()
	svc := user.NewService(repo)
	ctx := context.Background()

	id, _ := svc.Create(ctx, "Alice", "alice@example.com", "password123", "password123")

	if err := svc.UpdatePassword(ctx, id, "newpassword1", "newpassword1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, valid, _ := svc.ValidateCredentials(ctx, "alice@example.com", "newpassword1")
	if !valid {
		t.Fatal("expected valid=true after password update")
	}

	_, valid, _ = svc.ValidateCredentials(ctx, "alice@example.com", "password123")
	if valid {
		t.Fatal("expected valid=false for old password")
	}
}
