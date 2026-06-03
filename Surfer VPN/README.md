# 🏄 Surfer VPN — Telegram Mini App

> Свобода без границ. Быстрый и безопасный интернет для тебя.

Telegram Mini App для VPN-сервиса **Surfer VPN**. Светлая океаническая тема,
аниме-маскот, splash-экран, кастомные волны и анимации. Реализован 1:1 по
готовому дизайну.

## Стек

- **React 19** + **TypeScript** (strict)
- **Vite 6**
- **TailwindCSS** (база) + кастомный CSS дизайна (`src/styles/surfer.css`)
- **Telegram Mini Apps SDK** (официальный `telegram-web-app.js`)
- **React Router 6**

## Запуск

```bash
npm install
npm run dev      # дев-сервер (http://localhost:5173)
npm run build    # tsc -b + vite build -> dist/
npm run preview  # предпросмотр прод-сборки
npm run lint     # tsc --noEmit
```

Открыть в Telegram: захостить `dist/` (статик-хостинг с HTTPS) и указать URL в
@BotFather → Bot Settings → Menu Button / Web App.

## Структура

```
public/
└── images/            # арт и логотипы (surfer-hero.png, logo-symbol.svg, logo-full.svg)
src/
├── components/surf/   # компоненты экрана (по дизайну)
│   ├── icons.tsx      # набор иконок (Ic.*)
│   ├── logos.tsx      # SymbolLogo, FullLogo
│   ├── TgHeader.tsx   # верхний бар Telegram
│   ├── Hero.tsx       # hero с маскотом + Waves
│   ├── Waves.tsx      # анимированные волны
│   ├── UserCard.tsx   # карточка: статус, дни, инфо-строки
│   ├── Actions.tsx    # «Открыть Happ» + «Скопировать ключ»
│   ├── InstallGrid.tsx# сетка установки 2×2
│   ├── Toast.tsx      # тост
│   └── Splash.tsx     # splash-экран загрузки
├── pages/             # HomePage, NotFoundPage
├── hooks/             # useTelegram, useUser, useClipboard
├── lib/               # utils, constants, telegram (SDK), api, mock-data
├── types/             # User, Subscription, VpnServer, Platform, Telegram*
├── styles/surfer.css  # весь визуальный CSS дизайна
├── App.tsx            # роутинг
├── main.tsx           # точка входа + init Telegram
└── index.css          # tailwind + шрифт Plus Jakarta Sans
```

## Данные и backend

Всё на **mock-данных** (`src/lib/mock-data.ts`). Слой `src/lib/api.ts`
(`getCurrentUser` / `getVpnKey`) изолирует UI от источника: при подключении
бэкенда меняется только тело этих функций (`fetch` к `VITE_API_BASE_URL` с
авторизацией через `initData`) — компоненты и хуки не трогаются, они зависят
только от типа `User`. Личность (id, имя) подтягивается из
`Telegram.WebApp.initDataUnsafe.user`, если приложение открыто внутри Telegram.

## Ключевые сценарии

- **Открыть Happ** → `https://sub.surfervpn.com/open?key={vpnKey}`
  (`buildHappUrl`, `src/lib/constants.ts`), нативный `openLink` + хаптика.
- **Скопировать ключ** → `navigator.clipboard` + тост «Ключ успешно скопирован».
- **Splash** показывается ~2.1 c и пока грузятся данные.

> `DESIGN-PORT.md` — заметки о переносе дизайна (можно удалить).
