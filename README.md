# Money tracker

Сервис для учета личных трат.

## Запуск сервиса

```bash
$ docker-compose up
```

## Конфигурация

Минимальная конфигурация уже описана в `docker-compose.yml`. Пример полного конфига:
```yaml
database:
  url: "postgres://money_tracker:money_tracker@postgres:5432/money_tracker?sslmode=disable"

server:
  port: "8080"

jwt:
  expiration: 24h
  secret: "secret-key"

worker:
  poll_interval: 5s

log:
  level: debug
```

## Миграции

[Тут](internal/infrastructure/repository/postgres/migrations/00001_init_schema.sql).

## Создание пользователя

Обычного пользователя можно создать через endpoint `/users`.
```bash
$ curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

Администратора создать через endpoint **нельзя**. Это намеренное ограничение, так как аккаунт(ы) для администраторов должны заводиться вручную теми, у кого есть полный доступ к сервису. Проще всего сделать так:

```bash
$ curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "password123"
  }'

$ docker-compose exec postgres psql -U money_tracker -d money_tracker -c "UPDATE users SET role = 'admin' WHERE email = 'admin@example.com';"
```

## Описание эндпоинтов

На эндпоинте /swagger или [тут](docs/swagger.yaml).
