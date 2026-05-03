package user_test

import (
	"context"
	"errors"
	"testing"

	"golang.org/x/crypto/bcrypt"

	"github.com/malfoit/SimpleProject/internal/model"
	repouser "github.com/malfoit/SimpleProject/internal/repository/user"
	"github.com/malfoit/SimpleProject/internal/service/user"
)

// ---------------------------------------------------------------------------
// Mock repository
// ---------------------------------------------------------------------------

// mockRepo реализует repository.UserRepo через function fields.
// Если поле не установлено — метод возвращает нули/nil.
// Это позволяет каждому тесту задать только те вызовы, которые ему нужны,
// и провалить тест, если неожиданно был вызван лишний метод.
type mockRepo struct {
	createFn             func(ctx context.Context, u *model.User) error
	getByIDFn            func(ctx context.Context, id string) (*model.User, error)
	getByEmailFn         func(ctx context.Context, email string) (*model.User, error)
	updateFn             func(ctx context.Context, id string, name, email *string) error
	updatePasswordHashFn func(ctx context.Context, id, hash string) error
	deleteFn             func(ctx context.Context, id string) error
}

func (m *mockRepo) Create(ctx context.Context, u *model.User) error {
	if m.createFn != nil {
		return m.createFn(ctx, u)
	}
	return nil
}
func (m *mockRepo) GetByID(ctx context.Context, id string) (*model.User, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, nil
}
func (m *mockRepo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	if m.getByEmailFn != nil {
		return m.getByEmailFn(ctx, email)
	}
	return nil, nil
}
func (m *mockRepo) Update(ctx context.Context, id string, name, email *string) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, id, name, email)
	}
	return nil
}
func (m *mockRepo) UpdatePasswordHash(ctx context.Context, id, hash string) error {
	if m.updatePasswordHashFn != nil {
		return m.updatePasswordHashFn(ctx, id, hash)
	}
	return nil
}
func (m *mockRepo) Delete(ctx context.Context, id string) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, id)
	}
	return nil
}

// repoMustNotBeCalled возвращает mockRepo, все методы которого зафейлят тест —
// используется там, где сервис должен вернуть ошибку до обращения к репо.
func repoMustNotBeCalled(t *testing.T) *mockRepo {
	t.Helper()
	fail := func() { t.Helper(); t.Fatal("repo must not be called when input is invalid") }
	return &mockRepo{
		createFn:             func(_ context.Context, _ *model.User) error { fail(); return nil },
		getByIDFn:            func(_ context.Context, _ string) (*model.User, error) { fail(); return nil, nil },
		getByEmailFn:         func(_ context.Context, _ string) (*model.User, error) { fail(); return nil, nil },
		updateFn:             func(_ context.Context, _ string, _, _ *string) error { fail(); return nil },
		updatePasswordHashFn: func(_ context.Context, _, _ string) error { fail(); return nil },
		deleteFn:             func(_ context.Context, _ string) error { fail(); return nil },
	}
}

func ptr(s string) *string { return &s }

func bcryptHash(t *testing.T, password string) string {
	t.Helper()
	h, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("bcryptHash: %v", err)
	}
	return string(h)
}

// ---------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------

func TestCreate_RepoCalledWithHashedPassword(t *testing.T) {
	const plaintext = "securepass1"
	var capturedHash string

	repo := &mockRepo{
		createFn: func(_ context.Context, u *model.User) error {
			capturedHash = u.PasswordHash
			u.ID = "new-id"
			return nil
		},
	}
	svc := user.NewService(repo)

	id, err := svc.Create(context.Background(), "Alice", "alice@example.com", plaintext, plaintext)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != "new-id" {
		t.Errorf("expected id 'new-id', got %q", id)
	}
	if capturedHash == plaintext {
		t.Error("password must not be stored in plaintext")
	}
	if err = bcrypt.CompareHashAndPassword([]byte(capturedHash), []byte(plaintext)); err != nil {
		t.Errorf("stored hash does not match plaintext: %v", err)
	}
}

func TestCreate_RepoCalledWithTrimmedNameAndEmail(t *testing.T) {
	var gotName, gotEmail string
	repo := &mockRepo{
		createFn: func(_ context.Context, u *model.User) error {
			gotName = u.UserInfo.Name
			gotEmail = u.UserInfo.Email
			u.ID = "id"
			return nil
		},
	}
	svc := user.NewService(repo)

	_, _ = svc.Create(context.Background(), "  Alice  ", "  alice@example.com  ", "password123", "password123")

	if gotName != "Alice" {
		t.Errorf("expected trimmed name 'Alice', got %q", gotName)
	}
	if gotEmail != "alice@example.com" {
		t.Errorf("expected trimmed email 'alice@example.com', got %q", gotEmail)
	}
}

func TestCreate_RepoReturnsErrAlreadyExists_WrappedError(t *testing.T) {
	repo := &mockRepo{
		createFn: func(_ context.Context, _ *model.User) error {
			return repouser.ErrAlreadyExists
		},
	}
	svc := user.NewService(repo)

	_, err := svc.Create(context.Background(), "Alice", "alice@example.com", "password123", "password123")
	if err == nil {
		t.Fatal("expected error")
	}
	if errors.Is(err, repouser.ErrAlreadyExists) {
		t.Error("raw repo sentinel must not leak out of service")
	}
}

func TestCreate_RepoReturnsUnexpectedError_Propagated(t *testing.T) {
	boom := errors.New("db exploded")
	repo := &mockRepo{
		createFn: func(_ context.Context, _ *model.User) error { return boom },
	}
	svc := user.NewService(repo)

	_, err := svc.Create(context.Background(), "Alice", "alice@example.com", "password123", "password123")
	if !errors.Is(err, boom) {
		t.Errorf("expected boom error, got %v", err)
	}
}

func TestCreate_NameTooShort_RepoNotCalled(t *testing.T) {
	svc := user.NewService(repoMustNotBeCalled(t))
	_, err := svc.Create(context.Background(), "Al", "alice@example.com", "password123", "password123")
	if err == nil {
		t.Fatal("expected validation error for short name")
	}
}

func TestCreate_NameTooLong_RepoNotCalled(t *testing.T) {
	svc := user.NewService(repoMustNotBeCalled(t))
	longName := "A123456789012345678901234567890123456789012345678901" // 51 chars
	_, err := svc.Create(context.Background(), longName, "alice@example.com", "password123", "password123")
	if err == nil {
		t.Fatal("expected validation error for long name")
	}
}

func TestCreate_InvalidEmail_RepoNotCalled(t *testing.T) {
	svc := user.NewService(repoMustNotBeCalled(t))
	_, err := svc.Create(context.Background(), "Alice", "not-an-email", "password123", "password123")
	if err == nil {
		t.Fatal("expected validation error for invalid email")
	}
}

func TestCreate_PasswordTooShort_RepoNotCalled(t *testing.T) {
	svc := user.NewService(repoMustNotBeCalled(t))
	_, err := svc.Create(context.Background(), "Alice", "alice@example.com", "short", "short")
	if err == nil {
		t.Fatal("expected validation error for short password")
	}
}

func TestCreate_PasswordTooLong_RepoNotCalled(t *testing.T) {
	svc := user.NewService(repoMustNotBeCalled(t))
	long := string(make([]byte, 73)) // 73 chars
	_, err := svc.Create(context.Background(), "Alice", "alice@example.com", long, long)
	if err == nil {
		t.Fatal("expected validation error for long password")
	}
}

func TestCreate_PasswordsMismatch_RepoNotCalled(t *testing.T) {
	svc := user.NewService(repoMustNotBeCalled(t))
	_, err := svc.Create(context.Background(), "Alice", "alice@example.com", "password123", "different1")
	if err == nil {
		t.Fatal("expected validation error for mismatched passwords")
	}
}

// ---------------------------------------------------------------------------
// Get
// ---------------------------------------------------------------------------

func TestGet_ReturnsUserFromRepo(t *testing.T) {
	want := &model.User{ID: "u1", UserInfo: model.UserInfo{Name: "Alice", Email: "alice@example.com"}}
	repo := &mockRepo{
		getByIDFn: func(_ context.Context, id string) (*model.User, error) {
			if id != "u1" {
				t.Errorf("unexpected id passed to repo: %s", id)
			}
			return want, nil
		},
	}
	svc := user.NewService(repo)

	got, err := svc.Get(context.Background(), "u1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != want.ID {
		t.Errorf("expected id %s, got %s", want.ID, got.ID)
	}
}

func TestGet_RepoReturnsErrNotFound_WrappedError(t *testing.T) {
	repo := &mockRepo{
		getByIDFn: func(_ context.Context, _ string) (*model.User, error) {
			return nil, repouser.ErrNotFound
		},
	}
	svc := user.NewService(repo)

	_, err := svc.Get(context.Background(), "ghost")
	if err == nil {
		t.Fatal("expected error")
	}
	if errors.Is(err, repouser.ErrNotFound) {
		t.Error("raw repo sentinel must not leak out of service")
	}
}

func TestGet_RepoReturnsUnexpectedError_Propagated(t *testing.T) {
	boom := errors.New("storage failure")
	repo := &mockRepo{
		getByIDFn: func(_ context.Context, _ string) (*model.User, error) { return nil, boom },
	}
	svc := user.NewService(repo)

	_, err := svc.Get(context.Background(), "u1")
	if !errors.Is(err, boom) {
		t.Errorf("expected boom error, got %v", err)
	}
}

// ---------------------------------------------------------------------------
// Update
// ---------------------------------------------------------------------------

func TestUpdate_PassesCorrectArgsToRepo(t *testing.T) {
	var gotID, gotName, gotEmail string
	repo := &mockRepo{
		updateFn: func(_ context.Context, id string, name, email *string) error {
			gotID = id
			if name != nil {
				gotName = *name
			}
			if email != nil {
				gotEmail = *email
			}
			return nil
		},
	}
	svc := user.NewService(repo)

	err := svc.Update(context.Background(), "u1", ptr("  Bob  "), ptr("  bob@example.com  "))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotID != "u1" {
		t.Errorf("expected id 'u1', got %q", gotID)
	}
	if gotName != "Bob" {
		t.Errorf("expected trimmed name 'Bob', got %q", gotName)
	}
	if gotEmail != "bob@example.com" {
		t.Errorf("expected trimmed email 'bob@example.com', got %q", gotEmail)
	}
}

func TestUpdate_BothNil_RepoCalledWithNilNil(t *testing.T) {
	called := false
	repo := &mockRepo{
		updateFn: func(_ context.Context, _ string, name, email *string) error {
			called = true
			if name != nil || email != nil {
				t.Error("expected both name and email to be nil")
			}
			return nil
		},
	}
	svc := user.NewService(repo)
	_ = svc.Update(context.Background(), "u1", nil, nil)
	if !called {
		t.Error("repo.Update must be called even when both fields are nil")
	}
}

func TestUpdate_NameTooShort_RepoNotCalled(t *testing.T) {
	svc := user.NewService(repoMustNotBeCalled(t))
	err := svc.Update(context.Background(), "u1", ptr("Al"), nil)
	if err == nil {
		t.Fatal("expected validation error for short name")
	}
}

func TestUpdate_NameTooLong_RepoNotCalled(t *testing.T) {
	svc := user.NewService(repoMustNotBeCalled(t))
	long := string(make([]byte, 51))
	err := svc.Update(context.Background(), "u1", ptr(long), nil)
	if err == nil {
		t.Fatal("expected validation error for long name")
	}
}

func TestUpdate_InvalidEmail_RepoNotCalled(t *testing.T) {
	svc := user.NewService(repoMustNotBeCalled(t))
	err := svc.Update(context.Background(), "u1", nil, ptr("bad-email"))
	if err == nil {
		t.Fatal("expected validation error for invalid email")
	}
}

func TestUpdate_RepoReturnsErrNotFound_WrappedError(t *testing.T) {
	repo := &mockRepo{
		updateFn: func(_ context.Context, _ string, _, _ *string) error {
			return repouser.ErrNotFound
		},
	}
	svc := user.NewService(repo)

	err := svc.Update(context.Background(), "ghost", ptr("Alice"), nil)
	if err == nil {
		t.Fatal("expected error")
	}
	if errors.Is(err, repouser.ErrNotFound) {
		t.Error("raw repo sentinel must not leak out of service")
	}
}

func TestUpdate_RepoReturnsErrAlreadyExists_WrappedError(t *testing.T) {
	repo := &mockRepo{
		updateFn: func(_ context.Context, _ string, _, _ *string) error {
			return repouser.ErrAlreadyExists
		},
	}
	svc := user.NewService(repo)

	err := svc.Update(context.Background(), "u1", nil, ptr("taken@example.com"))
	if err == nil {
		t.Fatal("expected error")
	}
	if errors.Is(err, repouser.ErrAlreadyExists) {
		t.Error("raw repo sentinel must not leak out of service")
	}
}

func TestUpdate_RepoReturnsUnexpectedError_Propagated(t *testing.T) {
	boom := errors.New("db error")
	repo := &mockRepo{
		updateFn: func(_ context.Context, _ string, _, _ *string) error { return boom },
	}
	svc := user.NewService(repo)

	err := svc.Update(context.Background(), "u1", ptr("Alice"), nil)
	if !errors.Is(err, boom) {
		t.Errorf("expected boom error, got %v", err)
	}
}

// ---------------------------------------------------------------------------
// UpdatePassword
// ---------------------------------------------------------------------------

func TestUpdatePassword_RepoCalledWithHashNotPlaintext(t *testing.T) {
	const plaintext = "newpassword1"
	var capturedHash string

	repo := &mockRepo{
		updatePasswordHashFn: func(_ context.Context, id, hash string) error {
			if id != "u1" {
				t.Errorf("unexpected id: %s", id)
			}
			capturedHash = hash
			return nil
		},
	}
	svc := user.NewService(repo)

	if err := svc.UpdatePassword(context.Background(), "u1", plaintext, plaintext); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if capturedHash == plaintext {
		t.Error("password hash must not equal plaintext")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(capturedHash), []byte(plaintext)); err != nil {
		t.Errorf("stored hash does not match plaintext: %v", err)
	}
}

func TestUpdatePassword_PasswordTooShort_RepoNotCalled(t *testing.T) {
	svc := user.NewService(repoMustNotBeCalled(t))
	err := svc.UpdatePassword(context.Background(), "u1", "short", "short")
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestUpdatePassword_PasswordTooLong_RepoNotCalled(t *testing.T) {
	svc := user.NewService(repoMustNotBeCalled(t))
	long := string(make([]byte, 73))
	err := svc.UpdatePassword(context.Background(), "u1", long, long)
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestUpdatePassword_Mismatch_RepoNotCalled(t *testing.T) {
	svc := user.NewService(repoMustNotBeCalled(t))
	err := svc.UpdatePassword(context.Background(), "u1", "password123", "different1")
	if err == nil {
		t.Fatal("expected validation error for mismatched passwords")
	}
}

func TestUpdatePassword_RepoReturnsErrNotFound_WrappedError(t *testing.T) {
	repo := &mockRepo{
		updatePasswordHashFn: func(_ context.Context, _, _ string) error {
			return repouser.ErrNotFound
		},
	}
	svc := user.NewService(repo)

	err := svc.UpdatePassword(context.Background(), "ghost", "password123", "password123")
	if err == nil {
		t.Fatal("expected error")
	}
	if errors.Is(err, repouser.ErrNotFound) {
		t.Error("raw repo sentinel must not leak out of service")
	}
}

func TestUpdatePassword_RepoReturnsUnexpectedError_Propagated(t *testing.T) {
	boom := errors.New("storage failure")
	repo := &mockRepo{
		updatePasswordHashFn: func(_ context.Context, _, _ string) error { return boom },
	}
	svc := user.NewService(repo)

	err := svc.UpdatePassword(context.Background(), "u1", "password123", "password123")
	if !errors.Is(err, boom) {
		t.Errorf("expected boom error, got %v", err)
	}
}

// ---------------------------------------------------------------------------
// Delete
// ---------------------------------------------------------------------------

func TestDelete_RepoCalledWithCorrectID(t *testing.T) {
	var gotID string
	repo := &mockRepo{
		deleteFn: func(_ context.Context, id string) error {
			gotID = id
			return nil
		},
	}
	svc := user.NewService(repo)

	if err := svc.Delete(context.Background(), "u1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotID != "u1" {
		t.Errorf("expected id 'u1', got %q", gotID)
	}
}

func TestDelete_RepoReturnsErrNotFound_WrappedError(t *testing.T) {
	repo := &mockRepo{
		deleteFn: func(_ context.Context, _ string) error { return repouser.ErrNotFound },
	}
	svc := user.NewService(repo)

	err := svc.Delete(context.Background(), "ghost")
	if err == nil {
		t.Fatal("expected error")
	}
	if errors.Is(err, repouser.ErrNotFound) {
		t.Error("raw repo sentinel must not leak out of service")
	}
}

func TestDelete_RepoReturnsUnexpectedError_Propagated(t *testing.T) {
	boom := errors.New("storage failure")
	repo := &mockRepo{
		deleteFn: func(_ context.Context, _ string) error { return boom },
	}
	svc := user.NewService(repo)

	err := svc.Delete(context.Background(), "u1")
	if !errors.Is(err, boom) {
		t.Errorf("expected boom error, got %v", err)
	}
}

// ---------------------------------------------------------------------------
// ValidateCredentials
// ---------------------------------------------------------------------------

func TestValidateCredentials_CorrectPassword_ReturnsValidTrue(t *testing.T) {
	const plaintext = "correctpass"
	repo := &mockRepo{
		getByEmailFn: func(_ context.Context, email string) (*model.User, error) {
			if email != "alice@example.com" {
				t.Errorf("unexpected email: %s", email)
			}
			return &model.User{ID: "u1", PasswordHash: bcryptHash(t, plaintext)}, nil
		},
	}
	svc := user.NewService(repo)

	id, valid, err := svc.ValidateCredentials(context.Background(), "alice@example.com", plaintext)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !valid {
		t.Error("expected valid=true for correct password")
	}
	if id != "u1" {
		t.Errorf("expected id 'u1', got %q", id)
	}
}

func TestValidateCredentials_WrongPassword_ReturnsValidFalse(t *testing.T) {
	repo := &mockRepo{
		getByEmailFn: func(_ context.Context, _ string) (*model.User, error) {
			return &model.User{ID: "u1", PasswordHash: bcryptHash(t, "correctpass")}, nil
		},
	}
	svc := user.NewService(repo)

	id, valid, err := svc.ValidateCredentials(context.Background(), "alice@example.com", "wrongpass1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if valid {
		t.Error("expected valid=false for wrong password")
	}
	if id != "" {
		t.Errorf("expected empty id, got %q", id)
	}
}

func TestValidateCredentials_UserNotFound_ReturnsFalseNoError(t *testing.T) {
	repo := &mockRepo{
		getByEmailFn: func(_ context.Context, _ string) (*model.User, error) {
			return nil, repouser.ErrNotFound
		},
	}
	svc := user.NewService(repo)

	id, valid, err := svc.ValidateCredentials(context.Background(), "nobody@example.com", "password123")
	if err != nil {
		t.Fatalf("expected nil error for unknown user, got: %v", err)
	}
	if valid {
		t.Error("expected valid=false for unknown user")
	}
	if id != "" {
		t.Errorf("expected empty id, got %q", id)
	}
}

func TestValidateCredentials_RepoReturnsUnexpectedError_Propagated(t *testing.T) {
	boom := errors.New("storage failure")
	repo := &mockRepo{
		getByEmailFn: func(_ context.Context, _ string) (*model.User, error) { return nil, boom },
	}
	svc := user.NewService(repo)

	_, _, err := svc.ValidateCredentials(context.Background(), "alice@example.com", "password123")
	if !errors.Is(err, boom) {
		t.Errorf("expected boom error, got %v", err)
	}
}
