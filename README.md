# Geo Events Service (Go + NATS JetStream)

Сервис для приёма географических координат через HTTP API и асинхронной доставки событий через NATS JetStream.

Проект демонстрирует базовую, но приближенную к production архитектуру:

* REST API (входная точка)
* Message broker (NATS JetStream)
* Возможность подключения нескольких consumers

---

## Что делает сервис

1. Принимает HTTP POST запрос с координатами:

```json
{
  "lat": 55.7558,
  "lon": 37.6173
}
```

2. Публикует событие в **NATS JetStream**

3. Дальше:
    * любое количество consumers может подписаться
    * можно обрабатывать события независимо (логирование, аналитика, трекинг и т.д.)

---

## Зачем это нужно

Этот проект - учебный и демонстрационный, чтобы изучить:

* как работает NATS JetStream
* разницу между HTTP и event-driven архитектурой
* как Go сервис интегрируется с брокером сообщений

---

## Архитектура

```
[Client]
   |
   v
[HTTP API (Go)]
   |
   v
[NATS JetStream Stream]
   |
   +--> Consumer #1
   +--> Consumer #2
   +--> Consumer #N
```

---

## Стек технологий

* Go
* NATS
* NATS JetStream
* Docker Compose

---

## Конфигурация

Используется `.env` (опционально):

```env
NATS_URL=nats://localhost:4222
NATS_STREAM=GEO_STREAM
NATS_SUBJECT=geo.coordinates
HTTP_PORT=8080
```

---

## Запуск NATS

```bash
docker-compose up -d
```

Доступно:

* NATS: `localhost:4222`
* Monitoring UI: [http://localhost:8222](http://localhost:8222)

---

## Запуск сервиса

```bash
go run ./cmd/app
```

Сервис стартует на:

```
http://localhost:8080
```

---

## 📡 API

### POST `/coordinates`

Отправка координат:

```bash
curl -X POST http://localhost:8080/coordinates \
  -H "Content-Type: application/json" \
  -d '{"lat": 55.75, "lon": 37.61}'
```

Ответ:

```
202 Accepted
```