## Запуск в Docker

### 1. Клонировать репозиторий
```bash
git clone https://github.com/w12qwi/MentorApiProject.git
cd MentorApiProject
```

### 2. Запустить через Docker Compose
```bash
docker-compose up --build
```

Сервис будет доступен на `http://localhost:8080`.

Миграции применяются автоматически при старте приложения.

---

## Запуск локально (без Docker)

### 1. Клонировать репозиторий
```bash
git clone https://github.com/w12qwi/MentorApiProject.git
cd MentorApiProject
```

### 2. Создать `.env` файл в корне проекта (пример ниже)

### 3. Запустить приложение
```bash
go run ./cmd/main.go
```

Сервис будет доступен на `http://localhost:9999`.

---

## Пример .env файла

```
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=mentor
POSTGRES_PASSWORD=mentor
POSTGRES_DB=calculations_project
POSTGRES_SSL=disable

GRPC_SERVER_HOST=localhost
GRPC_SERVER_PORT=50051
GRPC_TIMEOUT=10

KAFKA_BROKER_HOST=localhost
KAFKA_BROKER_PORT=29092
KAFKA_TOPIC=calculations-topic
KAFKA_CONSUMER_GROUP=calculations-consumer-group
KAFKA_DLQ_TOPIC=calculations-dlq

JAEGER_PORT=14268
JAEGER_HOST=localhost
JAEGER_TRACES_ENDPOINT=/api/traces
```

---

## API

### POST /calculate

**Body:**
```json
{
    "numA": 10,
    "numB": 5,
    "sign": "+"
}
```

**Response:**
```json
{
    "result": 15
}
```

Доступные знаки: `+`, `-`, `*`, `/`


### GET /calculations/{id}
Получить вычисление по ID.

---

### GET /calculations
Получить вычисления c использованием фильтров(если фильтров нет клиент получает все существующие вычисления).

**Body:**
```json
{
    "date": "2024-01-15",
    "dateFrom": "2024-01-15",
    "dateTo": "2024-01-15",
    "sign": "+"
}
```
