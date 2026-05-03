# Аудит проекта SimpleProject

**Дата**: 2026-05-03  
**Ревьюер**: Claude Code (claude-sonnet-4-6)  
**Репозиторий**: `github.com/malfoit/SimpleProject`

---

## 1. Обзор

**SimpleProject** — gRPC-микросервис на Go для управления пользователями. Использует Protocol Buffers и in-memory хранилище.

| Параметр | Значение |
|---|---|
| Язык | Go 1.25.0 |
| Транспорт | gRPC |
| Хранилище | In-memory (map + sync.RWMutex) |
| Генерация кода | protoc + protoc-gen-go + protoc-gen-go-grpc |
| Сборка | Makefile |
| Тесты | Отсутствуют |

**Текущее состояние**: проект **не компилируется**.

---

## 2. Статус компиляции

### Критические ошибки компиляции

#### 2.1 Импорт несуществующих пакетов (`cmd/server/di.go`)

```go
import (
    githubRepo "github.com/malfoit/SimpleProject/internal/repository/github"
    githubService "github.com/malfoit/SimpleProject/internal/service/github"
)
```

Директории `internal/repository/github/` и `internal/service/github/` **не существуют**. Проект не соберётся.

#### 2.2 Сигнатура функции не совпадает с использованием (`internal/repository/user/create.go`)

```go
// Сигнатура:
func (r *repo) Create(ctx context.Context, name, email string) error {
    // Использование внутри тела:
    Password: password, // ← переменная 'password' не объявлена
}
```

Параметр `password` отсутствует в сигнатуре, но используется внутри функции — **compilation error**.

#### 2.3 Дублирование метода `Create`

Метод `Create` объявлен одновременно в `internal/repository/user/repo.go` (с правильной сигнатурой) и в `internal/repository/user/create.go` (с неправильной). Один из них лишний.

---

## 3. Архитектура

### 3.1 Структура проекта

```
SimpleProject/
├── api/user/v1/user.proto          # Определение gRPC-сервиса
├── bin/                            # Бинарники protoc-плагинов
├── cmd/server/
│   ├── main.go                     # Точка входа
│   └── di.go                       # DI-контейнер (BROKEN)
├── internal/
│   ├── config/config.go            # Конфигурация
│   ├── handler/user/handler.go     # gRPC-хендлер (1 из 6 методов)
│   ├── model/user.go               # Доменная модель
│   ├── repository/
│   │   ├── repo.go                 # Интерфейс репозитория
│   │   └── user/
│   │       ├── repo.go             # Реализация (OK)
│   │       └── create.go           # ORPHANED / BROKEN
│   └── service/
│       ├── service.go              # Интерфейс сервиса
│       ├── user/
│       │   ├── service.go          # Инициализация
│       │   └── create.go           # Логика создания пользователя
│       └── sync/service.go         # Сервис синхронизации
├── pkg/user/v1/                    # Сгенерированный код protobuf
├── github.com/                     # STALE: старые сгенерированные файлы
├── go.mod
├── go.sum
└── Makefile
```

### 3.2 gRPC API (user.proto)

Определено 6 RPC-методов:

| Метод | Реализован | Примечание |
|---|---|---|
| `Create` | ✅ | Частично (нет хэширования пароля) |
| `Get` | ❌ | Только заглушка (UnimplementedUserV1Server) |
| `Update` | ❌ | Только заглушка |
| `UpdatePassword` | ❌ | Только заглушка |
| `Delete` | ❌ | Только заглушка |
| `ValidateCredentials` | ❌ | Только заглушка |

### 3.3 Слои

Проект имеет правильную слоистую структуру (handler → service → repository), но она реализована лишь частично.

### 3.4 DI-контейнер

`cmd/server/di.go` создаёт зависимости вручную. Подход приемлем для небольшого сервиса, но содержит ссылки на несуществующие пакеты (см. §2.1).

---

## 4. Безопасность

### CRITICAL

#### 4.1 Пароли хранятся в открытом виде

В `internal/model/user.go`:
```go
type User struct {
    Name            string
    Email           string
    Password        string        // ← plaintext
    PasswordConfirm string        // ← не должен быть в модели хранилища
}
```

Пароль передаётся через все слои и сохраняется в in-memory хранилище без хэширования. Необходимо использовать `bcrypt` или `argon2id`.

#### 4.2 Отсутствует аутентификация

Все 6 gRPC-методов доступны без какой-либо аутентификации. Нет JWT, OAuth, API-ключей, нет interceptor для проверки токенов.

#### 4.3 Отсутствует авторизация

Нет ролей, нет RBAC, нет проверки прав. Любой клиент может вызвать любой метод.

### HIGH

#### 4.4 Нет TLS

gRPC-сервер запускается без TLS (`grpc.NewServer()` без `credentials.NewServerTLSFromFile`). Данные, включая пароли, передаются в открытом виде.

#### 4.5 Отсутствует rate limiting

Нет ограничений на количество запросов. Метод `ValidateCredentials` (когда будет реализован) будет уязвим к брутфорсу.

### MEDIUM

#### 4.6 Небезопасная валидация email

В `internal/service/user/create.go` валидация email проверяет только наличие `@` и `.`. Такая проверка пропустит большинство некорректных адресов.

#### 4.7 GitHub-токен в конфигурации

`GITHUB_TOKEN` передаётся как переменная среды, что правильно, но значение по умолчанию — пустая строка без проверки обязательности.

---

## 5. Качество кода

### 5.1 Мёртвый код и мусор

| Файл/директория | Проблема |
|---|---|
| `internal/repository/user/create.go` | Дублирует и ломает уже определённый метод `Create` |
| `github.com/WithSoull/SimpleProject/` | Старые сгенерированные файлы, не используются |
| `.golangci.pipeline.yaml` | Указан в Makefile, но файл отсутствует |

### 5.2 Неиспользуемые зависимости в `go.mod`

```
go.opentelemetry.io/otel          v1.39.0  — нигде не импортируется
go.opentelemetry.io/otel/sdk      v1.39.0  — нигде не импортируется
go.opentelemetry.io/otel/trace    v1.39.0  — нигде не импортируется
gonum.org/v1/gonum                v0.17.0  — нигде не импортируется
```

### 5.3 Именование

В `cmd/server/di.go` используются нечитаемые сокращения:
```go
gRepo, gSvc, uRepo, uSvc, uHnd
```

### 5.4 Поле `PasswordConfirm` в доменной модели

`PasswordConfirm` — это поле для валидации при создании, оно не должно попадать в `model.User`. После проверки совпадения паролей оно не нужно.

### 5.5 `service.go` — незавершённый интерфейс

`internal/service/service.go` содержит только один метод `Create`. Если появятся `Get`, `Update` и т.д., их нужно будет добавить сюда.

### 5.6 Логирование

Используется стандартный пакет `log` вместо структурированного логгера (`slog` из стандартной библиотеки Go 1.21+ или `zap`/`zerolog`).

---

## 6. Обработка ошибок

- Ошибки возвращаются через `errors.New()` без gRPC status codes.  
  Клиент получит `Unknown` вместо `InvalidArgument`, `AlreadyExists`, `NotFound` и т.д.
- Нет middleware для перехвата паник.
- Нет единого типа ошибок — нельзя разграничить бизнес-ошибки и инфраструктурные.

**Правильный подход для gRPC-хендлера:**
```go
return nil, status.Errorf(codes.InvalidArgument, "invalid email: %s", email)
```

---

## 7. Тестирование

**Покрытие тестами: 0%**

Не найдено ни одного `_test.go` файла. Для production-ready сервиса необходимы:

- Юнит-тесты сервисного слоя (валидация, бизнес-логика)
- Тесты репозитория (in-memory реализация легко тестируется)
- Интеграционные тесты хендлера (с реальным gRPC-клиентом)

---

## 8. Конфигурация

`internal/config/config.go` читает переменные среды с дефолтами:

| Переменная | Дефолт | Проблема |
|---|---|---|
| `GRPC_PORT` | `50051` | OK |
| `GITHUB_TOKEN` | `""` | Нет проверки на пустоту |
| `GITHUB_OWNER` | `"malfoit"` | Захардкожено |
| `GITHUB_REPO` | `"SimpleProject"` | Захардкожено |

Нет файла `.env.example`. Нет валидации обязательных параметров при старте.

---

## 9. DevOps / эксплуатация

| Аспект | Статус |
|---|---|
| CI/CD | Отсутствует |
| Health check endpoint | Отсутствует |
| Graceful shutdown | Отсутствует |
| Метрики / трейсинг | Зависимости есть, код отсутствует |
| Персистентность данных | Нет (in-memory) |
| Docker / docker-compose | Отсутствует |

---

## 10. Итоговая таблица

| Область | Оценка | Приоритет исправления |
|---|---|---|
| Компиляция | ❌ Не работает | CRITICAL |
| Безопасность паролей | ❌ Plaintext | CRITICAL |
| Аутентификация | ❌ Отсутствует | CRITICAL |
| Авторизация | ❌ Отсутствует | CRITICAL |
| Реализация API | ⚠️ 1/6 методов | HIGH |
| TLS | ❌ Отсутствует | HIGH |
| Тестирование | ❌ 0% | HIGH |
| Обработка ошибок | ⚠️ Базовая | MEDIUM |
| Качество кода | ⚠️ Есть проблемы | MEDIUM |
| Конфигурация | ⚠️ Приемлемо | MEDIUM |
| CI/CD | ❌ Отсутствует | LOW |
| Документация | ❌ Отсутствует | LOW |

---

## 11. Рекомендации (по приоритету)

### CRITICAL — исправить до любого деплоя

1. **Починить компиляцию**
   - Удалить `internal/repository/user/create.go` (файл-дубликат с ошибкой)
   - Реализовать или убрать ссылки на `internal/repository/github` и `internal/service/github` в `cmd/server/di.go`

2. **Хэшировать пароли**
   - Добавить `golang.org/x/crypto/bcrypt` или `github.com/matthewhartstonge/argon2`
   - Хэшировать в сервисном слое перед передачей в репозиторий
   - Убрать `PasswordConfirm` из `model.User` — проверять только в сервисе

3. **Добавить аутентификацию**
   - Минимум: JWT-токены + unary interceptor
   - Хранить `user_id` в контексте после верификации токена

### HIGH — перед первым публичным использованием

4. **Реализовать оставшиеся RPC-методы** (`Get`, `Update`, `UpdatePassword`, `Delete`, `ValidateCredentials`)

5. **Добавить TLS**
   - Self-signed сертификат для разработки, реальный для production
   ```go
   creds, _ := credentials.NewServerTLSFromFile("cert.pem", "key.pem")
   grpc.NewServer(grpc.Creds(creds))
   ```

6. **Написать тесты** — минимум для `internal/service/user/create.go`

7. **Использовать gRPC status codes** в хендлере

8. **Добавить graceful shutdown** в `cmd/server/main.go`

### MEDIUM — улучшение качества

9. **Удалить мёртвый код**: `github.com/WithSoull/SimpleProject/`, `internal/repository/user/create.go`

10. **Убрать неиспользуемые зависимости**: `go mod tidy`

11. **Добавить структурированное логирование** (`log/slog`)

12. **Добавить `.golangci.pipeline.yaml`** или переименовать цель в Makefile

### LOW — хорошая практика

13. Добавить `README.md` с инструкцией по запуску

14. Настроить CI (GitHub Actions): lint + build + test

15. Добавить `Dockerfile` + `docker-compose.yml`

16. Добавить `.env.example`

---

*Аудит проведён статически на основании чтения исходного кода. Для полной проверки необходим запуск тестов и динамический анализ.*
