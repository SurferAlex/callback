# Аудит Surf VPN — отчёт

Дата: 2026-06-05

## Найденные проблемы

### Безопасность
- `POST /api/v1/user/subscription/mock-activate` был доступен любому пользователю с JWT/tma/internal impersonation
- `VPNAPI_INTERNAL_TOKEN` и `INTERNAL_TOKEN` легко рассинхронизировать → `unauthorized` у бота
- `.env` с секретами мог попасть в git (не было корневого `.gitignore`)

### Docker / деплой
- Сервис `nginx` в compose конфликтовал с host nginx (80/443)
- `migrate` без healthcheck postgres → гонка при старте
- `miniapp` build без `VITE_TELEGRAM_BOT_USERNAME` → Login Widget не работал в Docker
- Hardcoded production API URL в compose

### API / backend
- Мёртвый handler `MeJWT` (не зарегистрирован)
- Маршрут `/api/v1/monitor/targets` без потребителя в compose
- Слабое логирование ошибок (без `telegram_id`, без HTTP status)
- Статус подписки `trial` не отдавался в API → UI ломался

### user_bot
- Ошибки get/refresh config не логировались
- Мёртвый код: `formatDate`, `CBMockPay`, `GetMe`, `MockActivate` в vpnapi client

### Surfer VPN (frontend)
- Ошибки auth/login глотались без UI
- После login не восстанавливался исходный URL
- JWT после Mini App session не использовался в API (только tma)
- `useUser().error` не показывался на HomePage

## Что исправлено

- **mock-activate** перенесён на internal auth (`InternalUserAuth` + те же заголовки Telegram)
- Удалены: `MeJWT`, `monitor_http.go`, маршрут `/monitor/targets`
- Добавлены: `RequestLog` middleware, структурированные логи с `telegram_id`
- **docker-compose**: убран nginx, healthcheck postgres, env-based build args для miniapp
- **Dockerfile** miniapp: `VITE_TELEGRAM_BOT_USERNAME`, `VITE_SUB_BASE_URL`
- **Frontend**: ошибки auth/login/profile, redirect после login, Bearer приоритетнее tma, статус `trial`
- **user_bot**: логирование config ошибок, удалён мёртвый код
- Корневой `.gitignore`, обновлены `.env.example`, README, deploy docs

## Какие файлы изменены

| Область | Файлы |
|---------|-------|
| Compose | `docker-compose.yml`, `compose.env.example`, `.gitignore` |
| Backend | `internal/api/api.go`, `internal/handlers/user_http.go`, `internal/handlers/auth_http.go`, `internal/middleware/internal_user.go`, `internal/middleware/request_log.go` |
| Backend удалено | `internal/handlers/monitor_http.go` |
| user_bot | `internal/bot/router.go`, `internal/botapp/callbacks.go`, `internal/botapp/markup.go`, `vpnapi/client.go`, `.env.example` |
| Frontend | `src/lib/api.ts`, `src/contexts/AuthContext.tsx`, `src/pages/LoginPage.tsx`, `src/pages/HomePage.tsx`, `Dockerfile`, `.env.example`, `src/vite-env.d.ts` |
| Docs | `README.md`, `deploy/nginx/README.md`, `docs/AUDIT_REPORT.md` |

## Какие миграции добавлены

**Нет новых миграций** в этом аудите. Существующие: `000001`–`000006` (clients, xui, servers, users/subscriptions, trial, auth refresh).

## Какие API изменены

| Было | Стало |
|------|-------|
| `POST /api/v1/user/subscription/mock-activate` (user JWT/tma) | Тот же путь, но **только** `X-Internal-Token` + `X-Telegram-*` |
| `GET /api/v1/monitor/targets` | **Удалён** |
| `GET /api/v1/user/me` response | `subscription.status` может быть `"trial"` |

Остальные user/auth endpoints без изменений контракта.

## Что ещё рекомендуется улучшить

1. **Транзакции** при trial + provision 3x-ui (откат при ошибке панели)
2. **Атомарный refresh** JWT (revoke + issue в одной DB-транзакции)
3. **Revoke в 3x-ui** при deactivate клиента
4. **Rate limiting** на `/auth/*` и `/trial/activate`
5. **Ротация** утёкших токенов бота (BotFather)
6. `git filter-repo` если `.env` был в истории git
7. Uptime Kuma — отдельный compose или документ установки
8. Удалить неиспользуемые npm-пакеты в Surfer VPN (radix, sonner) — отдельный PR
9. `GIN_MODE=release` на production в `backend/.env`
