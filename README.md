# me-ai: Chat с собственной LLM (Ollama) и современным интерфейсом

---

## Системные требования

- **Go**: >= 1.21 (рекомендуется 1.24+)
- **Node.js**: >= 18.x (рекомендуется LTS)
- **npm**: >= 9.x
- **PostgreSQL**: >= 14.x
- **Ollama**: >= 0.1.32
- **Docker** (для быстрой локальной БД)
- **Linux/Mac/Windows** (тестировалось на MacOS и Ubuntu)

---

## Переменные окружения

Создайте файл `.env` в корне проекта:

```

# Строка подключения к PostgreSQL
DSN="host=localhost user=your_user password=your_password dbname=your_dbname port=5432 sslmode=disable"
# URL Ollama API (обычно локально)
URL="http://localhost:11434"

# JWT-секрет для подписи токенов
TOKEN="your_jwt_secret"
```

**Пояснения:**
- `DSN` — строка подключения к вашей базе данных PostgreSQL.
- `URL` — адрес Ollama API (порт по умолчанию 11434).
- `TOKEN` — секрет для подписи JWT (любой длинный случайный текст).

---

## Миграции базы данных

1. **Запустите PostgreSQL через Docker:**
   ```bash
   docker-compose up -d
   ```
2. **Примените миграции (пример с goose):**
   ```bash
   go install github.com/pressly/goose/v3/cmd/goose@latest
   goose -dir ./migrations postgres "$DSN" up
   ```
   Или примените SQL-файлы вручную через psql/pgAdmin.
3. **Проверьте структуру:**
   ```sql
   \dt
   \d users
   \d conversations
   \d messages
   ```

---

## API (Backend Endpoints)

### Auth
- `POST /api/auth/register` — регистрация
  - body: `{ "email": string, "password": string, "name": string }`
  - response: `{ "token": string }`
- `POST /api/auth/login` — вход
  - body: `{ "email": string, "password": string }`
  - response: `{ "token": string }`

### Чаты
- `GET /api/conversations` — список чатов пользователя
- `POST /api/conversations/create` — создать чат
  - body: `{ "title": string }`
- `POST /api/conversations/delete` — удалить чат
  - body: `{ "id": number }`
- `POST /api/conversations/rename` — переименовать чат
  - body: `{ "id": number, "title": string }`

### Сообщения
- `GET /api/messages?conversation_id=ID` — получить сообщения чата
- `POST /api/messages/delete` — удалить сообщение
  - body: `{ "id": number }`

### Общение с LLM
- `POST /api/chat` — отправить сообщение в чат (и получить ответ LLM)
  - body: `{ "conversation_id": number, "message": string }`
  - response: `{ "message": string, "timestamp": string }`
- `WS /api/ws` — WebSocket для real-time общения

---

## Структура проекта

```
me_ai/
  cmd/                # main.go — точка входа backend
  internal/
    auth/             # JWT, регистрация, вход, middleware
    llm/              # Интеграция с Ollama, обработчики чата, WebSocket
    models/           # Модели и репозитории (User, Conversation, Message)
    middleware/       # CORS, JWT и др. middleware
  pkg/
    db/               # Инициализация и подключение к БД
    jwt/              # Работа с JWT
    req/, res/        # Утилиты для обработки запросов/ответов
  configs/            # Загрузка переменных окружения
  migrations/         # SQL-миграции для PostgreSQL
  frontend/
    src/              # Исходники React-приложения
    public/           # Статика
    package.json      # Зависимости и скрипты
  docker-compose.yml  # Быстрый запуск PostgreSQL
  go.mod, go.sum      # Зависимости Go
  README.md           # Документация
```

**Пояснения:**
- `cmd/` — точка входа, запуск сервера
- `internal/` — вся бизнес-логика и API
- `frontend/` — современный UI на React + Vite + MUI
- `migrations/` — миграции для БД
- `pkg/` — вспомогательные пакеты

---

## Настройка Ollama

1. **Установите Ollama:** https://ollama.com/download
2. **Импортируйте или создайте свою модель:**
   ```bash
   ollama create my-model -f Modelfile
   # или загрузите готовую
   ollama pull llama2
   ```
3. **Запустите Ollama:**
   ```bash
   ollama run my-model
   # или
   ollama serve
   ```
4. **Проверьте, что API Ollama доступен:**
   ```bash
   curl http://localhost:11434/api/tags
   ```
5. **Пример запроса к Ollama (Python):**
   ```python
   import requests
   r = requests.post('http://localhost:11434/api/chat', json={
       'model': 'my-model',
       'messages': [
           {'role': 'user', 'content': 'Привет!'}
       ]
   })
   print(r.json())
   ```

---

## Примеры конфигурации

**.env:**
```
DSN=postgres://postgres17:meAI17@localhost:5432/postgres?sslmode=disable
URL=http://localhost:11434
TOKEN=your_jwt_secret
```

**docker-compose.yml:**
```
services:
  postgres:
    image: postgres:16.4
    environment:
      POSTGRES_USER: postgres17
      POSTGRES_PASSWORD: meAI17
      PGDATA: /data/postgres
    volumes:
      - ./postgres-data:/data/postgres
    ports:
      - "5432:5432"
```

---

## Troubleshooting

- **Ошибка: pq: column "login" of relation "users" does not exist**
  - Проверьте, что структура таблицы users соответствует ожидаемой (есть email, name, password).
  - Выполните миграции или вручную добавьте нужные столбцы.
- **Ошибка: Field validation for 'Email' failed on the 'required' tag**
  - Проверьте, что frontend отправляет поле email, а не login.
- **Ошибка: unsupported protocol scheme ""**
  - Проверьте, что переменная URL в .env задана и Ollama запущен.
- **CORS/401/403**
  - Проверьте настройки CORS и JWT-секреты на frontend и backend.
- **Frontend не видит backend**
  - Проверьте proxy в vite.config.ts:
    ```js
    proxy: {
      '/api': 'http://localhost:8081'
    }
    ```
- **Ollama не отвечает**
  - Проверьте, что Ollama запущен и порт совпадает с URL в .env.

---

## Разработка и полезные команды

- `npm run dev` — запуск frontend (Vite)
- `npm run build` — сборка frontend
- `npm run lint` — линтинг frontend
- `go run main.go` (в папке cmd) — запуск backend
- `docker-compose up -d` — запуск PostgreSQL
- `goose -dir ./migrations postgres "$DSN" up` — миграции
- `ollama run my-model` — запуск Ollama с вашей моделью

---

## Deployment (Production)

1. **Соберите frontend:**
   ```bash
   cd frontend
   npm run build
   ```
2. **Соберите backend (Go):**
   ```bash
   cd cmd
   go build -o meai-server main.go
   ```
3. **Настройте переменные окружения и БД.**
4. **Запустите Ollama с вашей моделью.**
5. **Запустите backend и раздачу статики (например, через nginx или встроенный сервер).**
6. **(Опционально) Используйте Docker для деплоя всех сервисов.**

---

## Дообучение LLM (Fine-tuning)

Код для дообучения вашей LLM-модели  рекомендуется хранить в отдельной папке.

**Что включать:**
- Jupyter/Colab-ноутбук с процессом обучения (пример: `finetune.ipynb`)
- Python-скрипты для подготовки датасета, запуска обучения, конвертации модели
- Пример датасета (или ссылку на него)
- Инструкцию по запуску обучения и сохранению модели для Ollama


**Интеграция с Ollama:**
- После дообучения сохраните модель в Ollama:
  ```
  ollama create my-model -f Modelfile
  ```
- Подробнее — в `training/README.md` (если добавите).


---

## Контакты

Автор: Степан Коротеев  
Telegram: [@Stepankoroteev](https://t.me/Stepankoroteev)



---

## Важно: настройте system prompt и имя модели

Перед запуском убедитесь, что в файле `internal/llm/llm.go`:
- В поле `System` (system prompt) указан ваш собственный промт, который будет использоваться для общения с LLM.
- В поле `Model` указано имя вашей модели в Ollama (например, `model9`, `my-llm`, `llama2:custom` и т.д.).

**Пример:**
```go
reqBody := OllamaRequest{
    Model:   "my-llm", // <-- ваше имя модели
    System:  "Вы — ассистент, который ...", // <-- ваш system prompt
    ...
}
```

# ME-AI
