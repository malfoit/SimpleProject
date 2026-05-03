package user_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/malfoit/SimpleProject/internal/model"
	"github.com/malfoit/SimpleProject/internal/repository/user"
)

func ptr(s string) *string { return &s }

func newUser(name, email string) *model.User {
	return &model.User{
		UserInfo:     model.UserInfo{Name: name, Email: email},
		PasswordHash: "hash",
	}
}

func seedUser(t *testing.T, r interface {
	Create(context.Context, *model.User) error
}, name, email string) *model.User {
	t.Helper()
	u := newUser(name, email)
	if err := r.Create(context.Background(), u); err != nil {
		t.Fatalf("seed Create: %v", err)
	}
	return u
}

// ---------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------

func TestCreate_AssignsID(t *testing.T) {
	r := user.NewRepository()
	u := newUser("Alice", "alice@example.com")

	if err := r.Create(context.Background(), u); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if u.ID == "" {
		t.Error("ID must be set after Create")
	}
}

func TestCreate_AssignsTimestamps(t *testing.T) {
	r := user.NewRepository()
	before := time.Now()
	u := newUser("Alice", "alice@example.com")

	_ = r.Create(context.Background(), u)

	if u.CreatedAt.Before(before) {
		t.Error("CreatedAt must be >= time before Create")
	}
	if !u.UpdatedAt.Equal(u.CreatedAt) {
		t.Error("UpdatedAt must equal CreatedAt on creation")
	}
}

func TestCreate_DuplicateEmail_ReturnsErrAlreadyExists(t *testing.T) {
	r := user.NewRepository()
	seedUser(t, r, "Alice", "alice@example.com")

	err := r.Create(context.Background(), newUser("Bob", "alice@example.com"))
	if !errors.Is(err, user.ErrAlreadyExists) {
		t.Errorf("expected ErrAlreadyExists, got %v", err)
	}
}

func TestCreate_TwoDistinctUsers_BothStored(t *testing.T) {
	r := user.NewRepository()
	u1 := seedUser(t, r, "Alice", "alice@example.com")
	u2 := seedUser(t, r, "Bob", "bob@example.com")

	if u1.ID == u2.ID {
		t.Error("IDs must be unique")
	}
}

// ---------------------------------------------------------------------------
// GetByID
// ---------------------------------------------------------------------------

func TestGetByID_ReturnsStoredUser(t *testing.T) {
	r := user.NewRepository()
	u := seedUser(t, r, "Alice", "alice@example.com")

	got, err := r.GetByID(context.Background(), u.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.UserInfo.Email != "alice@example.com" {
		t.Errorf("unexpected email: %s", got.UserInfo.Email)
	}
}

func TestGetByID_NotFound_ReturnsErrNotFound(t *testing.T) {
	r := user.NewRepository()

	_, err := r.GetByID(context.Background(), "nonexistent")
	if !errors.Is(err, user.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestGetByID_ReturnsCopy(t *testing.T) {
	r := user.NewRepository()
	u := seedUser(t, r, "Alice", "alice@example.com")

	got, _ := r.GetByID(context.Background(), u.ID)
	got.UserInfo.Name = "mutated"

	got2, _ := r.GetByID(context.Background(), u.ID)
	if got2.UserInfo.Name == "mutated" {
		t.Error("GetByID must return a copy, not a pointer to internal state")
	}
}

// ---------------------------------------------------------------------------
// GetByEmail
// ---------------------------------------------------------------------------

func TestGetByEmail_ReturnsStoredUser(t *testing.T) {
	r := user.NewRepository()
	seedUser(t, r, "Alice", "alice@example.com")

	got, err := r.GetByEmail(context.Background(), "alice@example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.UserInfo.Name != "Alice" {
		t.Errorf("unexpected name: %s", got.UserInfo.Name)
	}
}

func TestGetByEmail_NotFound_ReturnsErrNotFound(t *testing.T) {
	r := user.NewRepository()

	_, err := r.GetByEmail(context.Background(), "nobody@example.com")
	if !errors.Is(err, user.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestGetByEmail_ReturnsCopy(t *testing.T) {
	r := user.NewRepository()
	seedUser(t, r, "Alice", "alice@example.com")

	got, _ := r.GetByEmail(context.Background(), "alice@example.com")
	got.UserInfo.Name = "mutated"

	got2, _ := r.GetByEmail(context.Background(), "alice@example.com")
	if got2.UserInfo.Name == "mutated" {
		t.Error("GetByEmail must return a copy, not a pointer to internal state")
	}
}

// ---------------------------------------------------------------------------
// Update
// ---------------------------------------------------------------------------

func TestUpdate_Name_ChangesName(t *testing.T) {
	r := user.NewRepository()
	u := seedUser(t, r, "Alice", "alice@example.com")

	if err := r.Update(context.Background(), u.ID, ptr("Alice V2"), nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, _ := r.GetByID(context.Background(), u.ID)
	if got.UserInfo.Name != "Alice V2" {
		t.Errorf("expected name 'Alice V2', got '%s'", got.UserInfo.Name)
	}
}

func TestUpdate_Email_ChangesEmailAndIndex(t *testing.T) {
	r := user.NewRepository()
	u := seedUser(t, r, "Alice", "alice@example.com")

	if err := r.Update(context.Background(), u.ID, nil, ptr("new@example.com")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// accessible by new email
	got, err := r.GetByEmail(context.Background(), "new@example.com")
	if err != nil {
		t.Fatalf("expected user at new email: %v", err)
	}
	if got.ID != u.ID {
		t.Error("ID mismatch after email update")
	}

	// old email must be gone
	_, err = r.GetByEmail(context.Background(), "alice@example.com")
	if !errors.Is(err, user.ErrNotFound) {
		t.Error("old email index must be removed after email update")
	}
}

func TestUpdate_Email_Taken_ReturnsErrAlreadyExists(t *testing.T) {
	r := user.NewRepository()
	u1 := seedUser(t, r, "Alice", "alice@example.com")
	seedUser(t, r, "Bob", "bob@example.com")

	err := r.Update(context.Background(), u1.ID, nil, ptr("bob@example.com"))
	if !errors.Is(err, user.ErrAlreadyExists) {
		t.Errorf("expected ErrAlreadyExists, got %v", err)
	}
}

func TestUpdate_SameEmail_NoError(t *testing.T) {
	r := user.NewRepository()
	u := seedUser(t, r, "Alice", "alice@example.com")

	if err := r.Update(context.Background(), u.ID, nil, ptr("alice@example.com")); err != nil {
		t.Errorf("updating with same email must not fail: %v", err)
	}
}

func TestUpdate_NotFound_ReturnsErrNotFound(t *testing.T) {
	r := user.NewRepository()

	err := r.Update(context.Background(), "ghost", ptr("X"), nil)
	if !errors.Is(err, user.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestUpdate_BothNil_OnlyBumpsUpdatedAt(t *testing.T) {
	r := user.NewRepository()
	u := seedUser(t, r, "Alice", "alice@example.com")
	before := u.UpdatedAt

	_ = r.Update(context.Background(), u.ID, nil, nil)

	got, _ := r.GetByID(context.Background(), u.ID)
	if !got.UpdatedAt.After(before) {
		t.Error("UpdatedAt must be bumped even when name and email are nil")
	}
}

// ---------------------------------------------------------------------------
// UpdatePasswordHash
// ---------------------------------------------------------------------------

func TestUpdatePasswordHash_ChangesHash(t *testing.T) {
	r := user.NewRepository()
	u := seedUser(t, r, "Alice", "alice@example.com")

	if err := r.UpdatePasswordHash(context.Background(), u.ID, "newhash"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, _ := r.GetByID(context.Background(), u.ID)
	if got.PasswordHash != "newhash" {
		t.Errorf("expected hash 'newhash', got '%s'", got.PasswordHash)
	}
}

func TestUpdatePasswordHash_BumpsUpdatedAt(t *testing.T) {
	r := user.NewRepository()
	u := seedUser(t, r, "Alice", "alice@example.com")
	before := u.UpdatedAt

	_ = r.UpdatePasswordHash(context.Background(), u.ID, "newhash")

	got, _ := r.GetByID(context.Background(), u.ID)
	if !got.UpdatedAt.After(before) {
		t.Error("UpdatedAt must be bumped after UpdatePasswordHash")
	}
}

func TestUpdatePasswordHash_NotFound_ReturnsErrNotFound(t *testing.T) {
	r := user.NewRepository()

	err := r.UpdatePasswordHash(context.Background(), "ghost", "hash")
	if !errors.Is(err, user.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

// ---------------------------------------------------------------------------
// Delete
// ---------------------------------------------------------------------------

func TestDelete_RemovesUser(t *testing.T) {
	r := user.NewRepository()
	u := seedUser(t, r, "Alice", "alice@example.com")

	if err := r.Delete(context.Background(), u.ID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err := r.GetByID(context.Background(), u.ID)
	if !errors.Is(err, user.ErrNotFound) {
		t.Error("GetByID must return ErrNotFound after Delete")
	}
}

func TestDelete_RemovesEmailIndex(t *testing.T) {
	r := user.NewRepository()
	u := seedUser(t, r, "Alice", "alice@example.com")

	_ = r.Delete(context.Background(), u.ID)

	_, err := r.GetByEmail(context.Background(), "alice@example.com")
	if !errors.Is(err, user.ErrNotFound) {
		t.Error("email index must be removed after Delete")
	}
}

func TestDelete_AllowsReuseOfEmail(t *testing.T) {
	r := user.NewRepository()
	u := seedUser(t, r, "Alice", "alice@example.com")
	_ = r.Delete(context.Background(), u.ID)

	if err := r.Create(context.Background(), newUser("Alice2", "alice@example.com")); err != nil {
		t.Errorf("email must be reusable after Delete: %v", err)
	}
}

func TestDelete_NotFound_ReturnsErrNotFound(t *testing.T) {
	r := user.NewRepository()

	err := r.Delete(context.Background(), "ghost")
	if !errors.Is(err, user.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}
