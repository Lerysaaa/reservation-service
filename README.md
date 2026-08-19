# reservation-service

Небольшой backend-сервис для резервирования товаров на Go и PostgreSQL.

Идея проекта — корректно уменьшать остаток товара при создании резерва и возвращать его при отмене. При одновременных запросах строка товара блокируется внутри транзакции, поэтому сервис не должен создавать резервов больше доступного остатка.

## Что умеет

- создавать товары и хранить текущий остаток;
- показывать список товаров;
- создавать резерв на нужное количество единиц;
- получать резерв по id;
- отменять резерв и возвращать количество в остаток;
- возвращать `409 Conflict`, если товара недостаточно или резерв уже отменён.

## Стек

Go, `net/http`, PostgreSQL, pgx, Docker Compose, GitHub Actions.

## Запуск

Нужен Docker с Docker Compose.

```bash
docker compose up --build
```

После запуска API доступен на `http://localhost:8080`.

Проверка:

```bash
curl http://localhost:8080/health
```

Создать товар:

```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{"name":"Mechanical keyboard","stock":5}'
```

Создать резерв:

```bash
curl -X POST http://localhost:8080/api/v1/reservations \
  -H "Content-Type: application/json" \
  -d '{"product_id":1,"quantity":2}'
```

Отменить резерв:

```bash
curl -X POST http://localhost:8080/api/v1/reservations/1/cancel
```

## API

| Method | Path | Description |
|---|---|---|
| `GET` | `/health` | health check |
| `POST` | `/api/v1/products` | create product |
| `GET` | `/api/v1/products` | list products |
| `POST` | `/api/v1/reservations` | create reservation |
| `GET` | `/api/v1/reservations/{id}` | get reservation |
| `POST` | `/api/v1/reservations/{id}/cancel` | cancel reservation |

## Почему используется транзакция

При резервировании сервис выполняет `SELECT ... FOR UPDATE` для строки товара. Пока транзакция не завершилась, второй конкурентный запрос не сможет одновременно прочитать и изменить тот же остаток. После проверки количества сервис уменьшает `stock` и создаёт запись резерва в рамках той же транзакции.

При отмене используется такой же подход: резерв блокируется, его количество возвращается товару, после чего статус меняется на `canceled`.

## Структура

```text
cmd/api/                    запуск HTTP-сервера
internal/httpapi/           HTTP handlers
internal/reservation/       модели и бизнес-логика
internal/storage/postgres/  работа с PostgreSQL
migrations/                 схема базы данных
```

## Тесты

```bash
go test ./...
```

Сейчас тестами покрыты базовая валидация бизнес-логики и поведение HTTP-слоя. Проверка также запускается в GitHub Actions.
