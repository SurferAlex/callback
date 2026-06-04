# Surfer VPN — control-plane

MVP: пользовательский Telegram-бот, Mini App / Web Cabinet и backend API (PostgreSQL + 3x-ui на отдельном VPS).

## Production (VPS + Nginx)

| URL | Порт на хосте | Сервис |
|-----|---------------|--------|
| https://app.surfwave.space | 3000 | Surfer VPN (статика) |
| https://api.surfwave.space | 8080 | Go API |
| https://status.surfwave.space | 3001 | Uptime Kuma |

Конфиги Nginx для копирования на сервер: [`deploy/nginx/`](deploy/nginx/README.md).

**backend/.env** (боевой):

```env
CORS_ORIGINS=https://app.surfwave.space
COOKIE_DOMAIN=.surfwave.space
COOKIE_SECURE=1
```

**user_bot/.env**: `MINI_APP_URL=https://app.surfwave.space`

**Сборка Mini App**: `VITE_API_BASE_URL=https://api.surfwave.space` (см. `Surfer VPN/.env.production`).

В BotFather: Web App URL → `https://app.surfwave.space`, Login Widget domain → `app.surfwave.space`.

## Сервисы (Docker Compose)

| Сервис | Назначение |
|--------|------------|
| `vpnapi` | Backend (`backend/`), порт **8080** |
| `user_bot` | Пользовательский Telegram-бот |
| `miniapp` | Статика Surfer VPN, порт **3000** → 80 в контейнере |
| `nginx` | Локальный reverse proxy (опционально), `nginx/conf.d/` |
| `postgres` | БД |

Admin-бот (`bot/`), `connect_bot/`, `monitor/` в compose **не включены**.

## Быстрый старт (локально)

1. Скопируйте env-файлы:
   - `backend/.env` из `backend/.env.example`
   - `user_bot/.env` из `user_bot/.env.example`
2. `TELEGRAM_BOT_TOKEN` — один токен в `backend/.env` и `user_bot/.env`.
3. Для локальной разработки API добавьте в `CORS_ORIGINS`: `http://localhost:5173`.
4. `docker compose up -d --build`

Локальный compose-nginx: `api.surfwave.space` / `app.surfwave.space` (пропишите в `/etc/hosts` → `127.0.0.1`).

## Auth API

- `POST /api/v1/auth/session/webapp` — Mini App (`tma` initData → JWT + cookie)
- `POST /api/v1/auth/session/widget` — браузер (Telegram Login Widget)
- `POST /api/v1/auth/refresh` / `POST /api/v1/auth/logout`

## User API

- `GET /api/v1/user/me` — `Authorization: Bearer` / `tma` / internal
- `POST /api/v1/user/trial/activate`
- `GET /api/v1/user/config`, `POST /api/v1/user/config/refresh`

## Mini App (локально)

```bash
cd "Surfer VPN"
cp .env.example .env.local
# .env.local: VITE_API_BASE_URL= пусто, Vite proxy /api → :8080
npm install && npm run dev
```
