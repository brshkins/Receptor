# Receptor

Веб-сервис и Telegram-бот для подбора рецептов по продуктам на руках: ручной список ингредиентов или фото с распознаванием (YOLO). Есть каталог рецептов, карточка с шагами, избранное и JWT-авторизация.

**Для кого:** домашние повара, учебные/MVP-команды, интеграция «каталог + ML-детекция + матчинг по БД».

---

## Архитектура

```text
┌─────────────┐     ┌─────────────┐
│  Frontend   │     │ Telegram    │
│  React/TS   │     │ Bot (Go)    │
└──────┬──────┘     └──────┬──────┘
       │ HTTP (JSON)        │
       └─────────┬──────────┘
                 ▼
         ┌───────────────┐
         │ Backend (Go)  │
         │ Gin + JWT +   │
         │ PostgreSQL    │
         └───────┬───────┘
                 │ POST multipart /detect
                 ▼
         ┌───────────────┐
         │ ML (FastAPI)  │
         │ YOLO → names  │
         └───────────────┘
```

- **Backend** — REST API, нормализация ингредиентов, матчинг рецептов, загрузка фото и вызов ML.
- **Bot** — long polling, FSM для логина/регистрации и подбора, клиент к тем же эндпоинтам.
- **ML** — `POST /detect`: файл поля `file`, ответ `{"ingredients": ["...", ...]}`.
- **Frontend** — SPA; часть сценариев должна быть синхронизирована с контрактом API (см. раздел «Известные проблемы»).

---

## API (кратко)

Базовый URL по умолчанию: `http://localhost:8080`.

| Метод | Путь | Auth | Описание |
|--------|------|------|----------|
| POST | `/auth/register` | — | Регистрация, ответ `{ "data": { "token" } }` |
| POST | `/auth/login` | — | Вход |
| GET | `/auth/me` | JWT | Профиль |
| GET | `/recipes` | — | Список (query: `search`, `category`, `max_time`, `sort`) |
| GET | `/recipes/:id` | — | Краткая карточка (без шагов/описания для полного UI) |
| GET | `/recipes/:id/details` | — | Детали: `description`, `steps[]` |
| POST | `/favorites/:id` | JWT | Добавить |
| DELETE | `/favorites/:id` | JWT | Удалить |
| GET | `/favorites` | JWT | Список |
| POST | `/match/by-ingredients` | — | Тело `{ "ingredients": ["..."] }`, ответ массив матчей |
| POST | `/upload` | — | `multipart/form-data`, поле **`file`** → детект + матчинг |

Успех: `{ "data": ... }`. Ошибка: `{ "error": "..." }`.

Подробнее: `docs/API.md` (частично устарел — нет `GET /recipes/:id/details`).

---

## Требования

- Go 1.22+
- Node.js 18+ (frontend)
- Python 3.10+ (ML)
- PostgreSQL
- Веса YOLO (`MODEL_PATH` для ML-сервиса)

---

## Запуск

### 1. PostgreSQL и переменные backend

Создайте БД и задайте в `.env` в каталоге `backend` (или в корне, см. `godotenv`):

```env
DB_URL=postgres://user:pass@localhost:5432/receptor
JWT_SECRET=your-secret
PORT=8080
ML_URL=http://localhost:8000
```

Примените миграции/схему согласно вашему процессу (в репозитории есть `backend/seeds/seed.go` для данных).

```bash
cd backend
go run ./cmd
```

### 2. ML-сервис

```bash
cd model/service
# MODEL_PATH — путь к .pt модели Ultralytics
set MODEL_PATH=C:\path\to\model.pt
python -m uvicorn app.main:app --host 0.0.0.0 --port 8000
```

Эндпоинт детекции: `POST http://localhost:8000/detect` (поле формы `file`).

### 3. Bot

```env
TELEGRAM_BOT_TOKEN=...
BACKEND_URL=http://localhost:8080
# опционально: файл сессий
SESSION_STORE_PATH=./sessions.json
```

```bash
cd bot
go run ./cmd
```

### 4. Frontend

```env
# frontend/.env
REACT_APP_API_URL=http://localhost:8080
```

```bash
cd frontend
npm install
npm start
```

CORS на backend разрешён для `http://localhost:3000`.

---

## Основные сценарии

1. **Регистрация / вход** — JWT в `localStorage` (web) или сессия бота с токеном; `/auth/me` для профиля.
2. **Каталог** — `/recipes` с фильтрами; детальный экран должен запрашивать **`/recipes/:id/details`** для описания и шагов.
3. **Подбор по продуктам** — `POST /match/by-ingredients`; пустой или нераспознанный набор даёт пустой массив (не ошибка).
4. **Подбор по фото** — `POST /upload` с полем `file`; backend вызывает ML и затем тот же матчинг, что и для списка ингредиентов.
5. **Избранное** — только с валидным JWT; конфликт при повторном добавлении — `409`.

---

## Известные проблемы (кратко)

| Область | Проблема |
|---------|----------|
| Frontend | Подбор: ожидается несуществующий формат ответа и эндпоинт `/match/by-image`; загрузка использует имя поля `image` вместо `file`. |
| Frontend | Страница рецепта запрашивает только `GET /recipes/:id` — нет шагов и описания из `/details`. |
| Docs | `docs/API.md` без `GET /recipes/:id/details`. |
| Bot | Пагинация `page`/`limit` не поддержана backend — при большом каталоге видны только первые N рецептов. |
| Конфиг | `ML_URL` не обязателен при старте backend, но нужен для `/upload`. |

Актуальный приоритетный план исправлений — в аудите проекта / issue tracker.

---

## Статус проекта

**MVP / beta, нестабильный UX на frontend-подборе и деталях рецепта** до выравнивания контрактов с backend. Backend, бот и ML-цепочка для фото в целом согласованы.

---

## Лицензия

Укажите лицензию при публикации репозитория.
