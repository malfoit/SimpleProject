package user

import (
	"context"
	"errors"
	"sync"

	"github.com/malfoit/SimpleProject/internal/model"
	"github.com/malfoit/SimpleProject/internal/repository"
)

var (
	ErrNotFound      = errors.New("user not found")
	ErrAlreadyExists = errors.New("user already exists")
)

type repo struct {
	mu       sync.RWMutex
	byID     map[string]*model.User
	emailIdx map[string]string // email -> id
}

func NewRepository() repository.UserRepo {
	return &repo{
		byID:     make(map[string]*model.User),
		emailIdx: make(map[string]string),
	}
}

// newID генерирует уникальный строковый идентификатор.
//
// Подсказка: используй пакет "crypto/rand" для генерации 16 случайных байт,
// затем fmt.Sprintf для форматирования в UUID-подобную строку вида:
// "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
// Лучше использовать пакет google/uuid :)
func newID() string {
	// TODO: сгенерируй 16 случайных байт через rand.Read
	// TODO: отформатируй их в строку через fmt.Sprintf с глаголом %x
	panic("TODO: implement newID")
}

// Create сохраняет нового пользователя в хранилище.
//
// Шаги:
//  1. Захвати write-lock: r.mu.Lock() / defer r.mu.Unlock()
//  2. Проверь r.emailIdx — если email уже есть, верни ErrAlreadyExists
//  3. Сгенерируй ID через newID() и запиши в user.ID
//  4. Проставь user.CreatedAt = time.Now(), user.UpdatedAt = user.CreatedAt
//  5. Сохрани КОПИЮ структуры (cp := *user), чтобы внешний код
//     не мог изменить внутреннее состояние через указатель
//  6. Добавь в r.byID[user.ID] и r.emailIdx[email] = user.ID
func (r *repo) Create(ctx context.Context, user *model.User) error {
	// TODO: реализуй метод
	panic("TODO: implement Create")
}

// GetByID возвращает пользователя по его ID.
//
// Шаги:
//  1. Захвати read-lock: r.mu.RLock() / defer r.mu.RUnlock()
//  2. Найди пользователя в r.byID — если нет, верни ErrNotFound
//  3. Верни КОПИЮ (cp := *u), а не оригинальный указатель
func (r *repo) GetByID(ctx context.Context, id string) (*model.User, error) {
	// TODO: реализуй метод
	panic("TODO: implement GetByID")
}

// GetByEmail возвращает пользователя по email.
//
// Шаги:
//  1. Захвати read-lock
//  2. Найди id в r.emailIdx по email — если нет, верни ErrNotFound
//  3. Получи пользователя из r.byID[id] и верни его КОПИЮ
func (r *repo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	// TODO: реализуй метод
	panic("TODO: implement GetByEmail")
}

// Update обновляет имя и/или email пользователя.
//
// Шаги:
//  1. Захвати write-lock
//  2. Найди пользователя в r.byID — если нет, верни ErrNotFound
//  3. Если email != nil И новый email отличается от текущего:
//     - проверь, не занят ли новый email (ErrAlreadyExists)
//     - удали старый ключ из r.emailIdx
//     - обнови u.UserInfo.Email и добавь новый ключ в r.emailIdx
//  4. Если name != nil — обнови u.UserInfo.Name
//  5. Обнови u.UpdatedAt = time.Now()
func (r *repo) Update(ctx context.Context, id string, name, email *string) error {
	// TODO: реализуй метод
	panic("TODO: implement Update")
}

// UpdatePasswordHash заменяет хэш пароля пользователя.
//
// Шаги:
//  1. Захвати write-lock
//  2. Найди пользователя в r.byID — если нет, верни ErrNotFound
//  3. Запиши новый passwordHash в u.PasswordHash
//  4. Обнови u.UpdatedAt = time.Now()
func (r *repo) UpdatePasswordHash(ctx context.Context, id, passwordHash string) error {
	// TODO: реализуй метод
	panic("TODO: implement UpdatePasswordHash")
}

// Delete удаляет пользователя из хранилища.
//
// Шаги:
//  1. Захвати write-lock
//  2. Найди пользователя в r.byID — если нет, верни ErrNotFound
//  3. Удали запись из r.emailIdx по u.UserInfo.Email
//  4. Удали запись из r.byID по id
func (r *repo) Delete(ctx context.Context, id string) error {
	// TODO: реализуй метод
	panic("TODO: implement Delete")
}
