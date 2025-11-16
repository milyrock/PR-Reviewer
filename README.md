Структура проекта
-----------------

```
root/
|- cmd/
|- config/
|- db/
|- docker/
|- e2e/
|- internal/
|  ├─ app/
|  ├─ config/
|  ├─ handlers/v1/
|  ├─ models/
|  ├─ repository/             
|  ├─ service/     
|  └─ test/             
|- docker-compose.yaml       
└─ makefile                  
```
Инструменты
-----------

Микросервис работает на Go версии >= 1.24

Для разворачивания проекта есть два способа:

Способ 1 (Docker Compose через Makefile)
--------------------------------------------------------

Установить Docker и Docker Compose. Затем в корне проекта выполнить:

```bash
make up
```

После этого через curl проверяем работоспособность проекта:

```bash
curl http://localhost:8080/health
```

Для остановки сервиса:

```bash
make down
```

Миграции применяются автоматически при старте приложения. 
Сервис и его зависимости поднимаются одной командой `make up`.
Сервис доступен на порту 8080.

Способ 2 (Docker Compose напрямую)
-----------------------------------

Альтернативный способ - запуск через docker-compose напрямую:

```bash
docker compose up -d --build
```

После запуска проверяем работоспособность:

```bash
curl http://localhost:8080/health
```

Для остановки:

```bash
docker compose down
```

Миграции применяются автоматически при старте приложения через `app.InitDB()`.
Файл `db/pr.sql` выполняется автоматически при подключении к базе данных.

Дополнительные команды
-----------------------

Для запуска E2E тестов:

```bash
make test-e2e
```

Для проверки кода линтером:

```bash
make lint
```

Взаимодействие
--------------

Основные endpoints:

#### Health Check
```bash
GET /health
```

#### Команды
```bash
POST /team/add          # Создать команду с участниками
GET  /team/get?team_name=<name>  # Получить команду
```

#### Пользователи
```bash
POST /users/setIsActive  # Установить флаг активности пользователя
GET  /users/getReview?user_id=<id>  # Получить PR'ы пользователя
```

#### Pull Request'ы
```bash
POST /pullRequest/create    # Создать PR и назначить ревьюверов
POST /pullRequest/merge     # Пометить PR как MERGED
POST /pullRequest/reassign  # Переназначить ревьювера
```

#### Статистика
```bash
GET /statistics  # Получить статистику по пользователям и PR'ам
```

Примеры запросов:

```bash
# Health Check
curl http://localhost:8080/health

# Создать команду
curl -X POST http://localhost:8080/team/add \
  -H "Content-Type: application/json" \
  -d '{
    "team_name": "backend",
    "members": [
      {"user_id": "u1", "username": "Alice", "is_active": true},
      {"user_id": "u2", "username": "Bob", "is_active": true}
    ]
  }'

# Получить PR'ы пользователя
curl "http://localhost:8080/users/getReview?user_id=u1"

# Создать PR
curl -X POST http://localhost:8080/pullRequest/create \
  -H "Content-Type: application/json" \
  -d '{
    "pull_request_id": "pr-1",
    "pull_request_name": "Add feature",
    "author_id": "u1"
  }'

# Пометить PR как MERGED
curl -X POST http://localhost:8080/pullRequest/merge \
  -H "Content-Type: application/json" \
  -d '{
    "pull_request_id": "pr-1"
  }'

# Переназначить ревьювера
curl -X POST http://localhost:8080/pullRequest/reassign \
  -H "Content-Type: application/json" \
  -d '{
    "pull_request_id": "pr-1",
    "old_reviewer_id": "u2"
  }'

```

Для удобного тестирования API можно использовать [Postman](https://www.postman.com/).

Примеры для Postman:

- `GET http://localhost:8080/health`

- `POST http://localhost:8080/team/add`
  ```json
  {
    "team_name": "backend",
    "members": [
      {"user_id": "u1", "username": "Alice", "is_active": true},
      {"user_id": "u2", "username": "Bob", "is_active": true}
    ]
  }
  ```

- `POST http://localhost:8080/users/setIsActive`
  ```json
  {
    "user_id": "u1",
    "is_active": false
  }
  ```


