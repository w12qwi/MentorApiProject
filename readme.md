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

Сервис будет доступен на `http://localhost:8080`.

---

## Пример .env файла

```
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=mentor
POSTGRES_PASSWORD=mentor
POSTGRES_DB=calculations_project
POSTGRES_SSL=disable
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

---

### GET /calculations
Получить все вычисления.

---

### GET /calculations/{id}
Получить вычисление по ID.

---

### GET /calculations/by-date
Получить вычисления за определённую дату.

**Body:**
```json
{
    "date": "2024-01-15"
}
```

---

### GET /calculations/by-date-range
Получить вычисления за диапазон дат.

**Body:**
```json
{
    "fromDate": "2024-01-01",
    "toDate": "2024-01-31"
}
```