# Surfer VPN — control-plane

MVP: Telegram-бот, Mini App / Web Cabinet, Go API, PostgreSQL, 3x-ui.

## Production (VPS)

| URL | Порт | Сервис |
|-----|------|--------|
| https://app.surfwave.space | 3000 | miniapp |
| https://api.surfwave.space | 8080 | vpnapi |

Host nginx: [`deploy/nginx/`](deploy/nginx/README.md). **Не** поднимайте `nginx` из docker-compose на VPS.

```bash
git clone https://github.com/SurferAlex/callback.git /opt/project
cd /opt/project
cp backend/.env.example backend/.env
cp user_bot/.env.example user_bot/.env
cp compose.env.example .env   # опционально, для build miniapp

docker compose up -d postgres
docker compose run --rm migrate
docker compose up -d --build vpnapi miniapp user_bot
```

`VPNAPI_INTERNAL_TOKEN` (бот) = `INTERNAL_TOKEN` (backend).

## Локально

```bash
docker compose up -d postgres migrate vpnapi
cd "Surfer VPN" && cp .env.example .env.local && npm run dev
```

В `backend/.env` для dev: `CORS_ORIGINS=http://localhost:5173`, `COOKIE_SECURE=0`, `COOKIE_DOMAIN=`.

## API (user)

- `POST /api/v1/auth/session/webapp` — Mini App
- `POST /api/v1/auth/session/widget` — браузер
- `POST /api/v1/auth/refresh`, `POST /api/v1/auth/logout`
- `GET /api/v1/user/me`, `POST /api/v1/user/trial/activate`
- `GET /api/v1/user/config`, `POST /api/v1/user/config/refresh`

Admin (`X-Internal-Token`): `/api/v1/clients/*`, `/api/v1/servers`, `POST /api/v1/user/subscription/mock-activate`.

## Стек

`backend/`, `user_bot/`, `Surfer VPN/`, `docker-compose.yml`
