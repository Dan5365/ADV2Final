# Online Booking System

Микросервисная архитектура для системы онлайн-бронирования.

## Архитектура

| Сервис | Протокол | Порт | Описание |
|---|---|---|---|
| user-service | gRPC | 50051 | Управление пользователями |
| order-service | gRPC | 50052 | Управление заказами |
| booking-service | gRPC | 50053 | Управление бронированиями |
| notification-service | — | — | Обработка событий RabbitMQ |
| api-gateway | HTTP | 8080 | REST API для клиентов |

## Технологии

- Go 1.25
- gRPC + Protocol Buffers
- PostgreSQL (3 БД: users, orders, bookings)
- RabbitMQ (обмен событиями)
- Docker + Docker Compose

---

## Требования

- [Docker Desktop](https://www.docker.com/products/docker-desktop/)
- Git

---

## Запуск

### 1. Клонировать репозиторий

```bash
git clone https://github.com/Dan5365/ADV2Final.git
cd ADV2Final
```

### 2. Запустить все сервисы

```bash
docker-compose up --build
```

Первый запуск занимает 2-5 минут (скачивание образов + сборка).  
Подождите пока в логах появится:

```
api-gateway       | api-gateway listening on :8080
user-service      | user-service listening on :50051
order-service     | order-service listening on :50052
booking-service   | booking-service listening on :50053
notification-service | notification-service consuming messages...
```

### 3. Проверить работу

```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Test User","email":"test@example.com","password":"pass123"}'
```

---

## API Endpoints

### Пользователи

**POST /users** — создать пользователя
```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"john@example.com","password":"secret123"}'
```

**GET /users/{id}** — получить пользователя
```bash
curl http://localhost:8080/users/1
```

### Бронирования

**POST /bookings** — создать бронирование
```bash
curl -X POST http://localhost:8080/bookings \
  -H "Content-Type: application/json" \
  -d '{"user_id":1,"resource":"Conference Room A","start_time":"2026-06-01T10:00:00Z","end_time":"2026-06-01T12:00:00Z"}'
```

**GET /bookings/{id}** — получить бронирование
```bash
curl http://localhost:8080/bookings/1
```

**GET /bookings?user_id={id}** — список бронирований пользователя
```bash
curl "http://localhost:8080/bookings?user_id=1&page=1&page_size=10"
```

**PUT /bookings/{id}/status** — обновить статус бронирования
```bash
curl -X PUT http://localhost:8080/bookings/1/status \
  -H "Content-Type: application/json" \
  -d '{"status":"confirmed"}'
```

### Заказы

**POST /orders** — создать заказ
```bash
curl -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -d '{"user_id":1,"items":[{"name":"Room Service","quantity":2,"price":50.00}]}'
```

**GET /orders/{id}** — получить заказ
```bash
curl http://localhost:8080/orders/1
```

---

## RabbitMQ

Панель управления: http://localhost:15672  
Логин: `guest` / Пароль: `guest`

События которые публикуются:
- `booking.created` — создано бронирование
- `booking.status_updated` — обновлён статус бронирования
- `order.created` — создан заказ
- `order.status_updated` — обновлён статус заказа

---

## Остановка

```bash
docker-compose down
```

Полная очистка включая базы данных:
```bash
docker-compose down -v
```

---

## Структура проекта

```
.
├── api-gateway/              # HTTP REST → gRPC
│   ├── cmd/main.go
│   └── internal/handler/
├── user-service/             # gRPC сервис пользователей
│   ├── cmd/main.go
│   ├── internal/
│   │   ├── repository/
│   │   └── server/
│   └── migrations/
├── order-service/            # gRPC сервис заказов
├── booking-service/          # gRPC сервис бронирований
├── notification-service/     # RabbitMQ consumer
├── gen/                      # Сгенерированный protobuf код
├── proto/                    # .proto файлы
└── docker-compose.yml
```
