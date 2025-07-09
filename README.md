# 🇰🇿 Сервис валидации ИИН Казахстана
[![Go Report Card](https://goreportcard.com/badge/github.com/toleubekov/check-iin-kaz)](https://goreportcard.com/report/github.com/toleubekov/check-iin-kaz)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.24.2+-blue.svg)](https://golang.org)

> Универсальное решение для валидации Индивидуальных Идентификационных Номеров (ИИН) Республики Казахстан. Используйте как **легковесную библиотеку** в вашем коде или как **готовый HTTP сервис**.

## 🎯 Два способа использования

### 📚 **Вариант 1: Легковесная библиотека** (рекомендуется)

Импортируйте только пакет валидации в ваш проект - без зависимостей, быстро и просто.

```bash
go get github.com/toleubekov/check-iin-kaz
```

```go
import "github.com/toleubekov/check-iin-kaz/iin"

// Быстрая проверка
if iin.IsValid("031231500126") {
    fmt.Println("ИИН корректный")
}

// Полная информация
info, err := iin.Validate("031231500126")
if err == nil {
    fmt.Printf("Пол: %s, Дата рождения: %s\n", info.Sex, info.DateOfBirth)
}
```

**Размер**: ~5-10 KB в итоговом бинарнике  
**Зависимости**: 0  
**Время старта**: мгновенно

### 🚀 **Вариант 2: HTTP сервис** 

Разверните готовый REST API сервер с базой данных для корпоративного использования.

```bash
git clone https://github.com/toleubekov/check-iin-kaz.git
cd check-iin-kaz
docker-compose up -d
```

```bash
# API для валидации
curl http://localhost:8080/iin_check/031231500126

# API для управления персонами
curl -X POST http://localhost:8080/people/info \
  -H "Content-Type: application/json" \
  -d '{"name":"Иван Иванов","iin":"031231500126","phone":"+77771234567"}'
```

**Включает**: PostgreSQL, REST API, админ панель  
**Применение**: корпоративные системы, микросервисы

---

## ✨ Особенности

- 🔍 **Полная валидация ИИН** с проверкой контрольной суммы
- 📊 **Извлечение информации**: пол, дата рождения, век рождения
- ⚡ **Sub-millisecond** производительность
- 🐳 **Docker support** для быстрого развертывания
- 📈 **Нагрузочное тестирование** в комплекте
- 🏛️ **Clean Architecture** для maintainability
- 🔐 **PostgreSQL** для надежного хранения данных
- 🚫 **Без зависимостей** при использовании как библиотека

## 📖 О формате ИИН

ИИН состоит из 12 цифр и кодирует следующую информацию:

```
[ГГ][ММ][ДД][V][NNNN][K]
 └─┘ └─┘ └─┘ │  └──┘ │
  │   │   │  │   │   └── Контрольная цифра
  │   │   │  │   └────── Порядковый номер (1000-9999)
  │   │   │  └────────── Век и пол (1-6)
  │   │   └───────────── День рождения (01-31)
  │   └───────────────── Месяц рождения (01-12)
  └───────────────────── Год рождения (00-99)
```

**Век и пол (7-я позиция):**
- `1` — мужчина, 1800-1899 гг. (XIX век)
- `2` — женщина, 1800-1899 гг. (XIX век)
- `3` — мужчина, 1900-1999 гг. (XX век)
- `4` — женщина, 1900-1999 гг. (XX век)
- `5` — мужчина, 2000-2099 гг. (XXI век)
- `6` — женщина, 2000-2099 гг. (XXI век)

---

## 📚 Использование как библиотеки

### 🚀 Быстрый старт

```bash
# Создать новый проект
mkdir my-app && cd my-app
go mod init my-app

# Установить библиотеку
go get github.com/toleubekov/check-iin-kaz

# Создать main.go
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "log"
    
    "github.com/toleubekov/check-iin-kaz/iin"
)

func main() {
    // Валидация с извлечением полной информации
    info, err := iin.Validate("031231500126")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Валидный: %t\n", info.Valid)
    fmt.Printf("Пол: %s\n", info.Sex)
    fmt.Printf("Дата рождения: %s\n", info.DateOfBirth)
    fmt.Printf("Век: %d\n", info.Century)
    fmt.Printf("Региональный код: %d\n", info.RegionCode)
}
EOF

# Запустить
go run main.go
```

### 📋 API библиотеки

#### Основные функции

```go
// Полная валидация с извлечением информации
info, err := iin.Validate("031231500126")

// Быстрая проверка валидности
isValid := iin.IsValid("031231500126")

// Совместимый API (для миграции с internal service)
valid, sex, dateOfBirth, err := iin.ValidateAndExtract("031231500126")

// Извлечение отдельных компонентов
sex, err := iin.ExtractSex("031231500126")
dateOfBirth, err := iin.ExtractDateOfBirth("031231500126")
```

#### Структура IINInfo

```go
type IINInfo struct {
    Valid       bool   `json:"valid"`            // Валидность ИИН
    Sex         string `json:"sex"`              // "male" или "female"
    DateOfBirth string `json:"date_of_birth"`    // DD.MM.YYYY
    Century     int    `json:"century"`          // Номер века рождения (19, 20, 21)
    RegionCode  int    `json:"region_code"`      // Региональный код (1000-9999)
}
```

### 🏗️ Примеры интеграции

#### HTTP сервер

```go
package main

import (
    "encoding/json"
    "net/http"
    
    "github.com/toleubekov/check-iin-kaz/iin"
)

func validateHandler(w http.ResponseWriter, r *http.Request) {
    iinStr := r.URL.Query().Get("iin")
    
    info, err := iin.Validate(iinStr)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(info)
}

func main() {
    http.HandleFunc("/validate", validateHandler)
    http.ListenAndServe(":8080", nil)
}
```

#### Middleware для валидации

```go
func IINValidationMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        iinStr := r.Header.Get("X-IIN")
        
        if iinStr != "" && !iin.IsValid(iinStr) {
            http.Error(w, "Invalid IIN", http.StatusBadRequest)
            return
        }
        
        next.ServeHTTP(w, r)
    })
}
```

#### Валидация в struct

```go
type User struct {
    Name string `json:"name"`
    IIN  string `json:"iin"`
}

func (u *User) Validate() error {
    if !iin.IsValid(u.IIN) {
        return errors.New("invalid IIN")
    }
    return nil
}
```

---

## 🚀 Использование как HTTP сервис

### Docker Compose (Рекомендуется)

```bash
# Клонирование репозитория
git clone https://github.com/toleubekov/check-iin-kaz.git
cd check-iin-kaz

# Запуск всех сервисов
docker-compose up -d

# Проверка логов
docker-compose logs -f server

# Тестовый запрос
curl http://localhost:8080/iin_check/031231500126
```

### 📚 API Документация сервиса

#### 🔍 Валидация ИИН

```http
GET /iin_check/{iin}
```

**Пример:**
```bash
curl http://localhost:8080/iin_check/031231500126
```

**Ответ:**
```json
{
  "correct": true,
  "sex": "male",
  "date_of_birth": "31.12.2003"
}
```

#### 👤 Управление персонами

**Создание записи:**
```http
POST /people/info
Content-Type: application/json

{
  "name": "Мыркымбаев Мыркымбай Мыркымбаевич",
  "iin": "031231500126",
  "phone": "+77771234567"
}
```

**Поиск по ИИН:**
```http
GET /people/info/iin/{iin}
```

**Поиск по имени:**
```http
GET /people/info/name/{name_part}
```

### ⚙️ Конфигурация сервиса

```env
# База данных
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=qwerty
DB_NAME=postgres
DB_SSLMODE=disable

# Сервер
SERVER_PORT=8080

# Нагрузочное тестирование
SERVER_URL=http://localhost:8080
NUM_GOROUTINES=10
NUM_REQUESTS=100
```

### Локальная разработка сервиса

```bash
# Установка зависимостей
go mod download

# Запуск PostgreSQL
docker run --name=iin-postgres -e POSTGRES_PASSWORD=qwerty -p 5432:5432 -d postgres:14-alpine

# Миграции
migrate -path ./schema -database 'postgres://postgres:qwerty@localhost:5432/postgres?sslmode=disable' up

# Запуск сервера
go run cmd/server/main.go
```

---

## 🧪 Тестирование

### Библиотека

```bash
# Тесты библиотеки
go test ./iin/...

# Тесты с покрытием
go test -cover ./iin/...

# Бенчмарки
go test -bench=. ./iin/...
```

### Сервис

```bash
# Запуск всех тестов
go test ./...

# Нагрузочное тестирование
go run cmd/stress-test/main.go
```

## 🎯 Применение

### 📚 Библиотека идеальна для:
- 🦋 **Легковесных приложений** - CLI утилиты, боты
- ⚡ **Быстрых валидаций** - формы регистрации, чекауты  
- 🔗 **Микросервисов** - без лишних зависимостей
- 📱 **Мобильных backend'ов** - минимальный footprint

### 🚀 HTTP сервис идеален для:
- 🏦 **Банковские системы** - KYC и верификация клиентов
- 🏢 **Госуслуги** - проверка документов граждан
- 💳 **Fintech платформы** - централизованная валидация
- 🏥 **Healthcare** - регистрация пациентов
- 🎓 **Корпоративные системы** - управление сотрудниками

## 📈 Производительность

### Библиотека:
- ⚡ **Sub-millisecond** валидация ИИН
- 🔄 **Zero allocation** для простых проверок
- 📦 **~5-10 KB** добавляется к бинарнику
- 🚫 **Без зависимостей**

### HTTP сервис:
- 🚀 **1000+ RPS** на обычном железе
- 🔄 **Concurrent** обработка запросов
- 📊 **Connection pooling** для базы данных
- 📈 **Horizontal scaling** ready

## 🏗️ Архитектура проекта

```
├── iin/                    # 📦 Публичная библиотека (без зависимостей)
│   ├── iin.go             # Основная функциональность
│   └── iin_test.go        # Тесты
├── examples/              # 📚 Примеры использования библиотеки
│   └── main.go           
├── cmd/                   # 🚀 HTTP сервис
│   ├── server/           # REST API сервер
│   └── stress-test/      # Нагрузочные тесты
├── internal/             # 🔒 Внутренние пакеты сервиса
│   ├── api/              # HTTP handlers
│   ├── model/            # Data models
│   ├── repository/       # Database layer
│   └── service/          # Бизнес-логика (использует iin/)
├── schema/               # 🗄️ Миграции базы данных
└── docs/                 # 📖 Документация
```

## 🤝 Участие в разработке

Мы приветствуем контрибьюции! Прочитайте [CONTRIBUTING.md](CONTRIBUTING.md) для деталей.

### Как помочь:

1. 🐛 Сообщите об ошибке
2. 💡 Предложите улучшение
3. 📝 Улучшите документацию
4. 🔧 Исправьте баг или добавьте фичу

## 📋 TODO

- [ ] GraphQL API для сервиса
- [ ] Валидация ИИН соседних стран  
- [ ] Кеширование результатов в сервисе
- [ ] Batch validation endpoint
- [ ] Metrics и мониторинг
- [ ] Rate limiting для API

## 🛡️ Безопасность

- ✅ Валидация всех входных данных
- 🔒 SQL injection protection (в сервисе)
- 🚫 No sensitive data logging
- 🔐 Environment-based configuration
- 🦋 Zero dependencies (в библиотеке)


[![pkg.go.dev](https://pkg.go.dev/badge/github.com/toleubekov/check-iin-kaz/iin.svg)](https://pkg.go.dev/github.com/toleubekov/check-iin-kaz/iin)

## 📄 Лицензия

Этот проект лицензирован под [MIT License](LICENSE).

## 👨‍💻 Автор

**Zhandos Toleubekov** - [GitHub](https://github.com/toleubekov)



---

<div align="center">

**Выберите то, что подходит вашему проекту:**

**📚 Нужна простая валидация?** → Используйте как библиотеку  
**🚀 Нужен полноценный сервис?** → Разворачивайте HTTP API

**Сделано с ❤️ для казахстанского dev-сообщества**

[⭐ Star](https://github.com/toleubekov/check-iin-kaz) • [🐛 Report Bug](https://github.com/toleubekov/check-iin-kaz/issues) • [💡 Request Feature](https://github.com/toleubekov/check-iin-kaz/issues)

</div>
