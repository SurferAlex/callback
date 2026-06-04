# Nginx на VPS (surfwave.space)

Единая точка входа на хосте: порты **80** → редирект на **443**, TLS Let's Encrypt.

| Поддомен | Upstream | Сервис |
|----------|----------|--------|
| `https://app.surfwave.space` | `127.0.0.1:3000` | Surfer VPN (Mini App / Web Cabinet) |
| `https://api.surfwave.space` | `127.0.0.1:8080` | Go API (`vpnapi`) |
| `https://status.surfwave.space` | `127.0.0.1:3001` | Uptime Kuma |

## Установка

```bash
sudo cp deploy/nginx/sites-available/*.conf /etc/nginx/sites-available/
sudo ln -sf /etc/nginx/sites-available/app.surfwave.space.conf /etc/nginx/sites-enabled/
sudo ln -sf /etc/nginx/sites-available/api.surfwave.space.conf /etc/nginx/sites-enabled/
sudo ln -sf /etc/nginx/sites-available/status.surfwave.space.conf /etc/nginx/sites-enabled/
sudo nginx -t && sudo systemctl reload nginx
```

Пути к сертификатам в конфигах должны совпадать с выводом `certbot certificates` (часто один cert на несколько имён).

## Сборка кабинета

```bash
cd "Surfer VPN"
docker build --build-arg VITE_API_BASE_URL=https://api.surfwave.space -t surf-miniapp .
# или статика на :3000 через serve / docker с портом 3000:80
```

## Env на сервере

**backend/.env**

```env
CORS_ORIGINS=https://app.surfwave.space
COOKIE_DOMAIN=.surfwave.space
COOKIE_SECURE=1
```

**user_bot/.env**

```env
MINI_APP_URL=https://app.surfwave.space
```

Папка `nginx/conf.d/` в корне репозитория — для **локального** Docker Compose (без TLS).
