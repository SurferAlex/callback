# Surfer VPN — control-plane

MVP: пользовательский Telegram-бот, Mini App и backend API (PostgreSQL + 3x-ui на отдельном VPS).

## Сервисы (Docker Compose)

| Сервис | Назначение |
|--------|------------|
| `vpnapi` | Backend (`backend/`), единственная точка работы с 3x-ui |
| `user_bot` | Пользовательский Telegram-бот |
| `miniapp` | Статика Surfer VPN (Telegram Mini App) |
| `nginx` | 80/443 → `app.*` и `api.*` |
| `postgres` | БД |

Admin-бот (`bot/`), `connect_bot/`, `monitor/` в compose **не включены**.

## Быстрый старт

1. Скопируйте env-файлы:
   - `backend/.env` из `backend/.env.example`
   - `user_bot/.env` из `.env.example`
2. Укажите `TELEGRAM_BOT_TOKEN` (один токен — и в `backend/.env`, и в `user_bot/.env` для Mini App auth).
3. Настройте `XUI_*` на URL панели на VPN VPS.
4. `docker compose up -d --build`

## Домены (nginx)

- `api.surfervpn.local` → backend API (`nginx/conf.d/api.conf`)
- `app.surfervpn.local` → Mini App (`nginx/conf.d/app.conf`)

Замените `server_name` на боевые `api.surfervpn`, `app.surfervpn` после DNS.

## User API

- `GET /api/v1/user/me` — профиль и подписка (`Authorization: tma <initData>` или internal + заголовки Telegram)
- `POST /api/v1/user/subscription/mock-activate` — `{"plan":"1m"|"3m"|"6m"|"12m"}`
- `GET /api/v1/user/config` — VLESS URI
- `POST /api/v1/user/config/refresh` — перевыпуск конфига

## Mini App (локально)

```bash
cd "Surfer VPN"
cp .env.example .env.local
npm install && npm run dev
```

В Telegram: `MINI_APP_URL` в `user_bot/.env` → URL приложения в BotFather.
