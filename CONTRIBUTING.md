# Руководство для контрибьюторов

Спасибо за интерес к проекту! Мы рады любому вкладу в развитие сервиса валидации ИИН Казахстана.

## Как внести вклад

### 🐛 Сообщить об ошибке

1. Проверьте, не была ли уже создана похожая issue
2. Создайте новую issue с:
   - Четким описанием проблемы
   - Шагами для воспроизведения
   - Ожидаемым и фактическим результатом
   - Информацией о среде (OS, Go version)

### 💡 Предложить улучшение

1. Создайте issue с:
   - Описанием предлагаемой функциональности
   - Объяснением, зачем это нужно
   - Примерами использования

### 🔧 Внести изменения в код

1. **Fork** репозитория
2. Создайте ветку для вашей фичи:
   ```bash
   git checkout -b feature/amazing-feature
   ```
3. Внесите изменения
4. Добавьте тесты для новой функциональности
5. Убедитесь, что все тесты проходят:
   ```bash
   go test ./...
   ```
6. Проверьте код линтером:
   ```bash
   golangci-lint run
   ```
7. Зафиксируйте изменения:
   ```bash
   git commit -m "Add amazing feature"
   ```
8. Отправьте изменения в ваш fork:
   ```bash
   git push origin feature/amazing-feature
   ```
9. Создайте Pull Request

## Стандарты кода

### Go Code Style

- Используйте `gofmt` для форматирования
- Следуйте [Effective Go](https://golang.org/doc/effective_go.html)
- Добавляйте комментарии к публичным функциям
- Используйте осмысленные имена переменных

### Commit Messages

Используйте конвенцию:
```
type(scope): description

- feat: новая функциональность
- fix: исправление ошибки
- docs: обновление документации
- test: добавление тестов
- refactor: рефакторинг без изменения функциональности
```

Примеры:
- `feat(api): add endpoint for IIN batch validation`
- `fix(service): correct checksum calculation for edge cases`
- `docs(readme): update API examples`

### Тестирование

- Все новые функции должны быть покрыты тестами
- Старайтесь поддерживать покрытие тестами выше 80%
- Включайте как позитивные, так и негативные тесты
- Используйте table-driven tests для множественных случаев

### Документация

- Обновляйте README.md при добавлении новых API endpoints
- Добавляйте комментарии к сложной бизнес-логике
- Включайте примеры использования

## Настройка среды разработки

### Требования

- Go 1.24.2+
- Docker и Docker Compose
- PostgreSQL (для локальной разработки)

### Локальная настройка

1. Клонируйте репозиторий:
   ```bash
   git clone https://github.com/toleubekov/check-iin-kaz.git
   cd check-iin-kaz
   ```

2. Установите зависимости:
   ```bash
   go mod download
   ```

3. Запустите PostgreSQL:
   ```bash
   docker run --name=iin-postgres -e POSTGRES_PASSWORD=qwerty -p 5432:5432 -d postgres:14-alpine
   ```

4. Примените миграции:
   ```bash
   migrate -path ./schema -database 'postgres://postgres:qwerty@localhost:5432/postgres?sslmode=disable' up
   ```

5. Запустите сервер:
   ```bash
   go run cmd/server/main.go
   ```

### Использование Docker

```bash
docker-compose up -d
```

## Структура проекта

```
├── cmd/                    # Точки входа приложений
│   ├── server/            # HTTP сервер
│   └── stress-test/       # Нагрузочное тестирование
├── internal/              # Приватный код приложения
│   ├── api/               # HTTP handlers и роутинг
│   ├── model/             # Модели данных
│   ├── repository/        # Слой доступа к данным
│   └── service/           # Бизнес-логика
├── schema/                # Миграции базы данных
└── docs/                  # Документация
```

## Code Review Process

1. Все Pull Requests проходят code review
2. CI pipeline должен пройти успешно
3. Требуется одобрение мейнтейнера
4. Изменения сквошатся при мердже

## Вопросы?

Если у вас есть вопросы:
- Создайте issue с меткой "question"
- Свяжитесь с мейнтейнером через GitHub

Спасибо за ваш вклад! 🚀
