# EventBooker

**EventBooker** — это сервис бронирования мероприятий с дедлайнами на подтверждение брони.

Проект позволяет:
- создавать мероприятия;
- бронировать места;
- подтверждать бронирование;
- автоматически отменять неподтвержденные брони по истечении времени;
- просматривать события и статистику по ним;
- просматривать список броней в админской части;
- отправлять email-уведомления об отмене брони через MailHog.

Сервис реализован как учебный MVP, но при этом включает важные для реальных систем вещи:
- транзакции;
- защиту от гонок при бронировании;
- фонового воркера для очистки просроченных броней;
- привязку бронирований к пользователям по email;
- ограничение: один пользователь не может иметь несколько активных броней на одно событие.

---

# Возможности

## Основной функционал
- создание мероприятий;
- получение списка мероприятий;
- получение информации о конкретном мероприятии;
- бронирование мест;
- подтверждение брони;
- автоматическая отмена истекших неподтвержденных броней.

## Дополнительно
- простой UI для пользователя и администратора;
- список броней по событию;
- уведомления по email через MailHog;
- поддержка пользователей через email;
- разные TTL брони для разных мероприятий.

---

# Архитектура

Проект разделен на несколько слоев:

- **handler** — HTTP-слой, принимает запросы и возвращает ответы;
- **service** — бизнес-логика;
- **repository** — работа с базой данных;
- **worker** — фоновые процессы;
- **notification** — отправка уведомлений;
- **domain** — доменные модели и ошибки.

Структура проекта:

```text
│   .env
│   .env.example
│   .gitignore
│   docker-compose.yml
│   Dockerfile
│   go.mod
│   go.sum
│
├───.vscode
│       settings.json
│
├───cmd
│   ├───eventbooker
│   │       main.go
│   │
│   └───eventbooker-worker
│           main.go
│
├───internal
│   ├───config
│   │       config.go
│   │
│   ├───domain
│   │       booking.go
│   │       errors.go
│   │       event.go
│   │       user.go
│   │
│   ├───handler
│   │   └───http
│   │           booking_handler.go
│   │           event_handler.go
│   │           middleware.go
│   │           respond.go
│   │           router.go
│   │
│   ├───notification
│   │       email.go
│   │       noop.go
│   │       notifier.go
│   │
│   ├───repository
│   │   │   repository.go
│   │   │   transactor.go
│   │   │
│   │   └───postgres
│   │           booking_repo.go
│   │           db.go
│   │           event_repo.go
│   │           tx.go
│   │           user_repo.go
│   │
│   ├───service
│   │       booking_service.go
│   │       event_service.go
│   │       service.go
│   │
│   ├───web
│   │       admin.html
│   │       admin.js
│   │       app.js
│   │       index.html
│   │       style.css
│   │
│   └───worker
│           booking_expirer.go
│
└───migrations
        001_create_event_bookings_table.down.sql
        001_create_event_bookings_table.up.sql
        002_users.down.sql
        002_users.up.sql│   .env
│   .env.example
│   .gitignore
│   docker-compose.yml
│   Dockerfile
│   go.mod
│   go.sum
│
├───.vscode
│       settings.json
│
├───cmd
│   ├───eventbooker
│   │       main.go
│   │
│   └───eventbooker-worker
│           main.go
│
├───internal
│   ├───config
│   │       config.go
│   │
│   ├───domain
│   │       booking.go
│   │       errors.go
│   │       event.go
│   │       user.go
│   │
│   ├───handler
│   │   └───http
│   │           booking_handler.go
│   │           event_handler.go
│   │           middleware.go
│   │           respond.go
│   │           router.go
│   │
│   ├───notification
│   │       email.go
│   │       noop.go
│   │       notifier.go
│   │
│   ├───repository
│   │   │   repository.go
│   │   │   transactor.go
│   │   │
│   │   └───postgres
│   │           booking_repo.go
│   │           db.go
│   │           event_repo.go
│   │           tx.go
│   │           user_repo.go
│   │
│   ├───service
│   │       booking_service.go
│   │       event_service.go
│   │       service.go
│   │
│   ├───web
│   │       admin.html
│   │       admin.js
│   │       app.js
│   │       index.html
│   │       style.css
│   │
│   └───worker
│           booking_expirer.go
│
└───migrations
        001_create_event_bookings_table.down.sql
        001_create_event_bookings_table.up.sql
        002_users.down.sql
        002_users.up.sql
```
# Стек
- Go
- PostgreSQL
- Docker / Docker Compose
- Chi
- MailHog
- HTML / CSS / JavaScript

# Доменные сущности
## Event

### Поля:
- id
- title
- starts_at
- capacity
- requires_payment
- booking_ttl_seconds
- created_at

## Booking

### Поля:
- id
- event_id
- user_id
- user_email
- status
- created_at
- expires_at
- confirmed_at
- cancelled_at

### Статусы:
- pending — бронь создана, но еще не подтверждена;
- confirmed — бронь подтверждена;
- cancelled — бронь отменена.

## User
### Поля:
- id
- email
- created_at

# Логика бронирования
## Если requires_payment = true
- при бронировании создается бронь со статусом pending;
- для брони рассчитывается expires_at;
- если бронь не подтверждена вовремя, worker автоматически переводит ее в cancelled;
- пользователю отправляется email об отмене.

## Если requires_payment = false
- бронь создается сразу со статусом confirmed;
- подтверждение через /confirm не требуется;
- попытка вызвать confirm для такого события приводит к ошибке.

# Ограничения
## Одно активное бронирование на событие для одного пользователя

На уровне базы данных введено ограничение:
- один и тот же пользователь не может иметь больше одной активной брони (pending или confirmed) на одно и то же событие.

Это реализовано через частичный уникальный индекс.

# Запуск проекта
## 1. Клонирование
```
git clone <repo-url>
cd event-booker
```
## 2. Подготовка .env
Создать .env по образцу .env.example.

Пример переменных:
```
HTTP_ADDR=:8080
PG_DSN = postgres://eventbooker:eventbooker@postgres:5432/eventbooker?sslmode=disable
DATABASE_URL = postgres://eventbooker:eventbooker@postgres:5432/eventbooker?sslmode=disable
POSTGRES_DB = eventbooker
POSTGRES_USER= eventbooker
POSTGRES_PASSWORD= eventbooker
SMTP_FROM= noreply@eventbooker.local
SMTP_HOST= mailhog
SMTP_PORT= 1025
```

## 3. Запуск
```
docker compose up --build
```

## 4. Полная пересборка с очисткой базы
```
docker compose down -v
docker compose up --build
```

# Сервисы в Docker Compose
## PostgreSQL
База данных проекта.

## migrate
Одноразовый контейнер для применения миграций.

## api
Основной HTTP-сервис.

## worker
Фоновый обработчик истекших броней.

## mailhog
Тестовый SMTP-сервер и веб-интерфейс для просмотра писем.

MailHog доступен по адресу:
```
http://localhost:8025
```

# HTTP API
## 1. Healthcheck
### GET /healthz
Проверка, что сервер запущен.

### Response
```
ok
```

## 2. Создание мероприятия
### POST /events
Создает новое мероприятие.

### Request JSON
```
{
  "title": "Backend Workshop",
  "starts_at": "2026-03-10T18:00:00Z",
  "capacity": 30,
  "requires_payment": true,
  "booking_ttl_seconds": 900
}
```
Поля
- title — название мероприятия;
- starts_at — дата и время начала в формате RFC3339;
- capacity — количество мест;
- requires_payment — требуется ли подтверждение/оплата;
- booking_ttl_seconds — срок жизни брони в секундах.

### Response
201 Created
```
{
  "id": "2bd519f9-eb33-4b16-b958-a9a563e5a46b"
}
```
## Возможные ошибки

#### 400 Bad Request
- некорректный JSON;
- невалидные поля.

## 3. Список мероприятий
### GET /events

Возвращает список мероприятий со статистикой.

### Query params
- limit — количество записей;
- offset — смещение.

Пример:
```
GET /events?limit=50&offset=0
```
### Response
```
{
  "items": [
    {
      "Event": {
        "ID": "2bd519f9-eb33-4b16-b958-a9a563e5a46b",
        "Title": "Backend Workshop",
        "StartsAt": "2026-03-10T18:00:00Z",
        "Capacity": 30,
        "RequiresPayment": true,
        "BookingTTLSeconds": 900,
        "CreatedAt": "2026-03-06T10:00:00Z"
      },
      "stats": {
        "pending": 2,
        "confirmed": 5,
        "free_seats": 23
      }
    }
  ],
  "limit": 50,
  "offset": 0
}
```
## 4. Получение события по ID
### GET /events/{id}
Возвращает информацию о конкретном мероприятии и его статистику.

### Response
```
{
  "Event": {
    "ID": "2bd519f9-eb33-4b16-b958-a9a563e5a46b",
    "Title": "Backend Workshop",
    "StartsAt": "2026-03-10T18:00:00Z",
    "Capacity": 30,
    "RequiresPayment": true,
    "BookingTTLSeconds": 900,
    "CreatedAt": "2026-03-06T10:00:00Z"
  },
  "stats": {
    "pending": 2,
    "confirmed": 5,
    "free_seats": 23
  }
}
```
### Возможные ошибки

#### 400 Bad Request
- некорректный UUID.

#### 404 Not Found
- событие не найдено.

## 5. Бронирование места
### POST /events/{id}/book
Создает бронь на мероприятие.

### Request JSON
```
{
  "user_email": "user@example.com"
}
```

### Логика
- если requires_payment = true, создается pending бронь;
- если requires_payment = false, создается confirmed бронь;
- если мест нет, возвращается ошибка;
- если у пользователя уже есть активная бронь на это событие, возвращается ошибка.

### Response для платного события
201 Created
```
{
  "booking_id": "0edb1db0-f875-431f-8518-a4c8cbc2a089",
  "status": "pending",
  "expires_at": "2026-03-06T11:15:00Z"
}
```

### Response для бесплатного события
201 Created
```
{
  "booking_id": "0edb1db0-f875-431f-8518-a4c8cbc2a089",
  "status": "confirmed",
  "expires_at": "2026-03-06T11:15:00Z"
}
```
### Возможные ошибки

#### 400 Bad Request
- некорректный UUID;
- некорректный JSON;
- пустой email.

#### 404 Not Found
- событие не найдено.

#### 409 Conflict
- нет свободных мест;
- пользователь уже имеет активную бронь на это событие.

## 6. Подтверждение брони
POST /events/{id}/confirm

Подтверждает ранее созданную бронь.

### Request JSON 
```
{
  "booking_id": "0edb1db0-f875-431f-8518-a4c8cbc2a089"
}
```
### Response
```
{
  "status": "confirmed"
}
```
### Возможные ошибки

#### 400 Bad Request
- некорректный UUID;
- некорректный JSON.

#### 404 Not Found
- бронь не найдена.

#### 409 Conflict
- бронь истекла;
- бронь уже в неподходящем статусе;
- подтверждение не требуется для этого события.

## 7. Список броней по событию
### GET /events/{id}/bookings

Возвращает список всех броней по конкретному событию.

Используется в административной части.

### Response
```
{
  "items": [
    {
      "ID": "0edb1db0-f875-431f-8518-a4c8cbc2a089",
      "EventID": "2bd519f9-eb33-4b16-b958-a9a563e5a46b",
      "UserID": "34886500-d0e9-4966-8ddb-19f228db417f",
      "UserEmail": "user@example.com",
      "Status": "pending",
      "CreatedAt": "2026-03-06T11:10:00Z",
      "ExpiresAt": "2026-03-06T11:15:00Z",
      "ConfirmedAt": null,
      "CancelledAt": null
    }
  ]
}
```
### Возможные ошибки

#### 400 Bad Request
- некорректный UUID.

# UI
## Пользовательская часть
Доступна по адресу:
<http://localhost:8080/ui/index.html>

### Возможности:
- просмотр списка мероприятий;
- бронирование;
- подтверждение брони;
- наблюдение за изменением статусов.

## Административная часть
Доступна по адресу:
<http://localhost:8080/ui/admin.html>

### Возможности:
- создание мероприятия;
- просмотр списка мероприятий;
- просмотр статистики по местам;
- просмотр списка броней по каждому событию.

# Worker
Фоновый worker запускается отдельным процессом и выполняет:
- поиск всех pending броней, у которых expires_at <= now;
- перевод их в статус cancelled;
- отправку email-уведомлений пользователям.
Интервал работы задается в коде и в текущей реализации составляет 5 секунд.

# Email-уведомления
Для уведомлений используется MailHog.
После отмены брони worker отправляет письмо на email пользователя.

Пример текста уведомления:

- бронирование было отменено;
- причина — бронь не была подтверждена вовремя;
- указывается название мероприятия.

Письма можно посмотреть в MailHog:
<http://localhost:8025>

# Безопасность и защита от гонок
## Транзакции
Бронирование выполняется в транзакции:

- событие блокируется через SELECT ... FOR UPDATE;
- считается количество активных мест;
- создается бронь.

Это защищает от ситуации, когда несколько пользователей одновременно бронируют последние места.

## Ограничение на пользователя
На уровне базы существует уникальный индекс, запрещающий одному пользователю иметь несколько активных броней на одно событие.

# Миграции
Проект использует SQL-миграции.

Основные миграции:
- создание таблицы events;
- создание таблицы bookings;
- создание таблицы users;
- добавление индексов и ограничений.

# Пример сценария использования
## Платное событие
- Администратор создает событие с requires_payment = true.
- Пользователь бронирует место.
- Бронь получает статус pending.
- Если пользователь подтверждает бронь — статус меняется на confirmed.
- Если не подтверждает — worker отменяет бронь и отправляет email.

## Бесплатное событие
- Администратор создает событие с requires_payment = false.
- Пользователь бронирует место.
- Бронь сразу получает статус confirmed.
- Подтверждение не требуется.

# Ограничения текущей версии
Это MVP-проект, поэтому:
- пользователь идентифицируется только по email;
- полноценная регистрация и авторизация отсутствуют;
- UI минималистичен и предназначен прежде всего для тестирования логики;
- уведомления реализованы только через email и MailHog;
- часть JSON-ответов использует стандартную сериализацию Go-структур, а не отдельные response DTO.
