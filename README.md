# FAQ Backend API (Accordion)

Backend для FAQ-аккордеона. Отдаёт REST API на `net/http`, хранит данные
в PostgreSQL, поддерживает CRUD, порядок и видимость.

## Возможности

- CRUD для FAQ
- UUID идентификаторы
- PostgreSQL
- JSON API
- CORS + recovery + логирование

## Стек

- Go (net/http)
- PostgreSQL
- Docker / Docker Compose

## Get Started

### Запуск через Docker

```
cp .env.example deployments/.env
make docker-up
make migrate-up
```

Корневой обработчик `/` не задан — используйте:

- `http://localhost:8080/healthz`
- `http://localhost:8080/api/v1/faqs`

### Локальный запуск

Убедитесь, что PostgreSQL запущен (можно использовать docker Postgres на `localhost:5432`).

```
POSTGRES_HOST=localhost \
POSTGRES_PORT=5432 \
POSTGRES_USER=postgres \
POSTGRES_PASSWORD=postgres \
POSTGRES_DB=faq \
POSTGRES_SSLMODE=disable \
SERVER_HOST=0.0.0.0 \
SERVER_PORT=8080 \
make run
```
## Swagger

Генерация документации:

```
make swagger
```

Swagger UI:

- http://localhost:8080/swagger/index.html

## API

Базовый URL: `/api/v1`

| Метод  | URL         | Описание             |
| ------ | ----------- | -------------------- |
| GET    | /faqs       | Список активных FAQ  |
| GET    | /faqs/{id}  | Получить один FAQ    |
| POST   | /faqs       | Создать FAQ          |
| PUT    | /faqs/{id}  | Обновить FAQ         |
| DELETE | /faqs/{id}  | Удалить FAQ          |


## Линтер

Используется `golangci-lint`.

Запуск линтера:

```
make lint
```

Установка (если `golangci-lint` не установлен):

```
make lint-install
```

## Полезные команды

```
make help
make run
make build
make fmt
make vet
make lint
make lint-install
make docker-logs
```
