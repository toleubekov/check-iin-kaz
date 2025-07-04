# 🇰🇿 Сервис валидации ИИН Казахстана

[![CI/CD](https://github.com/toleubekov/check-iin-kaz/actions/workflows/ci.yml/badge.svg)](https://github.com/toleubekov/check-iin-kaz/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/toleubekov/check-iin-kaz)](https://goreportcard.com/report/github.com/toleubekov/check-iin-kaz)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.24.2+-blue.svg)](https://golang.org)

> Высокопроизводительный сервис для валидации Индивидуальных Идентификационных Номеров (ИИН) Республики Казахстан с извлечением персональной информации.

## ✨ Особенности

- 🔍 **Полная валидация ИИН** с проверкой контрольной суммы
- 📊 **Извлечение информации**: пол, дата рождения, регион
- 🚀 **REST API** для интеграции с любыми системами
- 🐳 **Docker support** для быстрого развертывания
- 📈 **Нагрузочное тестирование** в комплекте
- 🏛️ **Clean Architecture** для maintainability
- 🔐 **PostgreSQL** для надежного хранения данных

## 🎯 Применение

Идеально подходит для:
- 🏦 **Банковские системы** - KYC и верификация клиентов
- 🏢 **Госуслуги** - проверка документов граждан
- 💳 **Fintech** - онлайн-верификация пользователей
- 🏥 **Healthcare** - регистрация пациентов
- 🎓 **Образование** - учет студентов и сотрудников

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
- `1` — мужчина, 1800-1899 гг.
- `2` — женщина, 1800-1899 гг.
- `3` — мужчина, 1900-1999 гг.
- `4` — женщина, 1900-1999 гг.
- `5` — мужчина, 2000-2099 гг.
- `6` — женщина, 2000-2099 гг.
- `0` — иностранные граждане

## 🚀 Быстрый старт

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
curl http://localhost:8080/iin_check/031231500123
```

### Локальная разработка

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

## 📚 API Документация

### 🔍 Валидация ИИН

```http
GET /iin_check/{iin}
```

**Пример:**
```bash
curl http://localhost:8080/iin_check/031231500123
```

**Ответ:**
```json
{
  "correct": true,
  "sex": "male",
  "date_of_birth": "31.12.2003"
}
```

### 👤 Управление персонами

#### Создание записи

```http
POST /people/info
Content-Type: application/json

{
  "name": "Назарбаев Нурсултан Абишевич",
  "iin": "400701111111",
  "phone": "+77771234567"
}
```

#### Поиск по ИИН

```http
GET /people/info/iin/{iin}
```

#### Поиск по имени

```http
GET /people/info/name/{name_part}
```

## ⚙️ Конфигурация

Настройка через переменные среды или `.env` файл:

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

## 🧪 Тестирование

```bash
# Запуск unit тестов
go test ./...

# Тесты с покрытием
go test -cover ./...

# Нагрузочное тестирование
go run cmd/stress-test/main.go
```

## 🏗️ Архитектура

```
cmd/
├── server/         # HTTP сервер
└── stress-test/    # Нагрузочные тесты

internal/
├── api/            # HTTP handlers
├── model/          # Data models
├── repository/     # Data access layer
└── service/        # Business logic

schema/             # Database migrations
```

**Принципы:**
- 🏛️ **Clean Architecture** - разделение ответственности
- 🔄 **Dependency Injection** - тестируемость
- 📦 **Single Responsibility** - один модуль = одна задача
- 🛡️ **Error Handling** - graceful degradation

## 📈 Производительность

- ⚡ **Sub-millisecond** валидация ИИН
- 🔄 **Concurrent** обработка запросов
- 📊 **Connection pooling** для базы данных
- 🚀 **Horizontal scaling** ready

## 🤝 Участие в разработке

Мы приветствуем контрибьюции! Прочитайте [CONTRIBUTING.md](CONTRIBUTING.md) для деталей.

### Как помочь:

1. 🐛 Сообщите об ошибке
2. 💡 Предложите улучшение
3. 📝 Улучшите документацию
4. 🔧 Исправьте баг или добавьте фичу

## 📋 TODO

- [ ] GraphQL API
- [ ] Валидация ИИН соседних стран
- [ ] Кеширование результатов
- [ ] Metrics и мониторинг
- [ ] Rate limiting
- [ ] Batch validation endpoint

## 🛡️ Безопасность

- ✅ Валидация всех входных данных
- 🔒 SQL injection protection
- 🚫 No sensitive data logging
- 🔐 Environment-based configuration

## 📄 Лицензия

Этот проект лицензирован под [MIT License](LICENSE).

## 👨‍💻 Автор

**Zhandos Toleubekov** - [GitHub](https://github.com/toleubekov)

---

<div align="center">

**Сделано с ❤️ для казахстанского dev-сообщества**

[⭐ Star](https://github.com/toleubekov/check-iin-kaz) • [🐛 Report Bug](https://github.com/toleubekov/check-iin-kaz/issues) • [💡 Request Feature](https://github.com/toleubekov/check-iin-kaz/issues)

</div>
