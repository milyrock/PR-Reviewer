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


Результаты нагрузки, полученные с помощью k6:

 TOTAL RESULTS 

    checks_total.......: 305     5.015489/s
    checks_succeeded...: 100.00% 305 out of 305
    checks_failed......: 0.00%   0 out of 305

    ✓ TEAM ADD status 201 or 400
    ✓ PR CREATE status 201 or 404 or 409
    ✓ PR REASSIGN status 200 or 404 or 409
    ✓ PR MERGE status 200 or 404
    ✓ GET REVIEW status 200

    HTTP
    http_req_duration..............: avg=11.55ms  min=1.11ms  med=10.41ms max=26.9ms p(90)=21.21ms p(95)=22.26ms
      { expected_response:true }...: avg=13.46ms  min=2.04ms  med=12.01ms max=26.9ms p(90)=21.48ms p(95)=22.42ms
    http_req_failed................: 20.00% 61 out of 305
    http_reqs......................: 305    5.015489/s

    EXECUTION
    iteration_duration.............: avg=996.89ms min=815.7ms med=999.5ms max=1s     p(90)=1s      p(95)=1s     
    iterations.....................: 61     1.003098/s
    vus............................: 1      min=1         max=1
    vus_max........................: 1      min=1         max=1

    NETWORK
    data_received..................: 99 kB  1.6 kB/s
    data_sent......................: 71 kB  1.2 kB/s
