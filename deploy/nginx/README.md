# Nginx на VPS (surfwave.space)

На production **не запускайте** docker-сервис `nginx` из compose — используйте host nginx.

| Поддомен | Upstream | Сервис |
|----------|----------|--------|
| `https://app.surfwave.space` | `127.0.0.1:3000` | Surfer VPN (miniapp) |
| `https://api.surfwave.space` | `127.0.0.1:8080` | Go API (vpnapi) |
| `https://status.surfwave.space` | `127.0.0.1:3001` | Uptime Kuma (отдельно) |

## Деплой

```bash
cd /opt/project
cp backend/.env.example backend/.env
cp user_bot/.env.example user_bot/.env
# Заполнить токены, XUI, JWT_SECRET
# VPNAPI_INTERNAL_TOKEN == INTERNAL_TOKEN

docker compose up -d postgres
docker compose run --rm migrate
docker compose up -d --build vpnapi miniapp user_bot
```

## Env (обязательно)

**backend/.env:** `DATABASE_URL`, `INTERNAL_TOKEN`, `TELEGRAM_BOT_TOKEN`, `JWT_SECRET`, `XUI_*`, `CORS_ORIGINS`, `COOKIE_DOMAIN`, `COOKIE_SECURE`

**user_bot/.env:** `TELEGRAM_BOT_TOKEN`, `VPNAPI_BASE_URL=http://vpnapi:8080`, `VPNAPI_INTERNAL_TOKEN`, `MINI_APP_URL`

**compose build (miniapp):** см. `compose.env.example` — `VITE_TELEGRAM_BOT_USERNAME` обязателен для Login Widget.

## Nginx

```bash
sudo cp deploy/nginx/sites-available/*.conf /etc/nginx/sites-available/
sudo ln -sf /etc/nginx/sites-available/app.surfwave.space.conf /etc/nginx/sites-enabled/
sudo ln -sf /etc/nginx/sites-available/api.surfwave.space.conf /etc/nginx/sites-enabled/
sudo nginx -t && sudo systemctl reload nginx
```
