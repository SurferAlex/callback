# Surfer VPN — Design Port Contract (READ FIRST)

We are **porting the provided Claude design 1:1** into this React 19 + TypeScript
+ Vite + Tailwind project. The design is the source of truth — **not** the old
components in `src/components/*` (those are being replaced).

## The design source (READ THESE — they define exact look & behavior)

All under `/Users/hicktox/Surfer VPN/Design Claude/Surf VPN/`:

- `Surfer VPN.html` — the **full CSS** (inside `<style>`) and structure. This is
  the visual source of truth: every class, color, gradient, shadow, animation.
- `src/app.jsx` — the React component tree & behavior (splash 2.1s, copy→toast, etc).
- `src/icons.jsx` — the icon set (`Ic.*`).
- `src/logos.jsx` — inline brand SVGs (`SymbolLogo`, `FullLogo`).
- `screenshots/app-final.png`, `screenshots/splash2.png` — rendered reference.

Reproduce it faithfully. Keep all Russian copy, gradients, animations
(waves sway, button shine, splash dots, cloud drift), shadows and radii exactly.

## Stack rules (unchanged)

- Alias `@/` → `src/`. Strict TS (`noUnusedLocals`/`noUnusedParameters`/`isolatedModules`):
  no unused imports/vars, use `import type`. **Named** function-component exports.
- No new npm deps. Icons/logos are inline SVG (ported), NOT lucide/shadcn here.
- Use `cn` from `@/lib/utils` only if helpful; the design is class-driven CSS.

## Adaptations from the design (IMPORTANT)

1. **Drop the iOS device frame.** The design HTML scales a 402×874 frame
   (`#stage`/`#scaler`/`#app`, `IOSDevice`, `ios-frame.jsx`). DO NOT port that.
   Instead: `.screen` fills the viewport and is centered with a max width for
   desktop preview. See the styles task below.
2. **Assets** live in `public/` and resolve at absolute URLs:
   - `/images/surfer-hero.png` (hero + splash art)
   - `/images/logo-symbol.svg`, `/images/logo-full.svg`
   In CSS use `url('/images/surfer-hero.png')`. In TSX use `<img src="/images/logo-full.svg">`.
3. **Data** comes from the existing data layer, not an inline literal:
   - `useUser()` from `@/hooks` → `{ user, loading, error, refetch }`, `user: User`.
   - `User` (see `@/types`): `telegramId`, `firstName`, `vpnKey`,
     `subscription { status, plan, expiresAt, daysLeft }`, `server { city, country, countryCode }`.
   - Date: `formatDate(iso)` from `@/lib/utils`. Days left: `user.subscription.daysLeft`.
   - Open Happ URL: `buildHappUrl(vpnKey)` from `@/lib/constants`; open via
     `openLink(...)` and fire `haptic("medium")` (both from `@/lib/telegram`).

## File ownership (build ONLY your files)

### Agent `styles` — `src/styles/surfer.css` + `src/index.css`
- Create `src/styles/surfer.css`: copy the ENTIRE `<style>` block from
  `Surfer VPN.html` VERBATIM, with these edits only:
  - Remove the device-frame rules: `#stage`, `#scaler`, `#app` drop-shadow,
    and the `body{display:flex;align-items:center;...overflow:hidden}` centering
    that assumes a scaled frame.
  - Keep the body background gradient + font. Set `body{min-height:100dvh}`.
  - `.screen`: `width:100%; max-width:440px; margin:0 auto; min-height:100dvh;
    position:relative; background:#fff; box-shadow:0 30px 60px rgba(20,70,120,.14)`.
    On `@media(max-width:440px)` drop the max-width/shadow so it's edge-to-edge.
  - `.scroll`: `min-height:100dvh` (was `height:100%`).
  - Fix asset URLs to `/images/surfer-hero.png`.
  - Keep `.tg-header` absolute with its `54px` top padding (status-bar safe area);
    you may reduce to `max(env(safe-area-inset-top),18px)` + 8px.
- Edit `src/index.css`: keep the `@tailwind base/components/utilities` lines and
  set the base `font-family` to `'Plus Jakarta Sans'`. Remove the old ocean theme
  custom CSS if it conflicts; keep it minimal. (Tailwind stays available.)
- `surfer.css` is imported by the page/wiring agent via `main.tsx` — just create it.

### Agent `icons` — `src/components/surf/icons.tsx`
Port `src/icons.jsx` to TSX. Export a typed `Ic` object with the SAME keys:
`Wave, ShieldCheck, Copy, Pin, Clock, Hash, Arrow, Close, Dots, Download,
Phone, Droid, Laptop, Monitor`. Each is a component `({ size?, color? }: { size?: number; color?: string }) => JSX`.
Copy the SVG paths exactly. Convert SVG attrs to JSX (`stroke-width`→`strokeWidth`,
`fill-opacity`→`fillOpacity`, `stroke-linecap`→`strokeLinecap`, etc.).
`export const Ic = {...}` and also `export type IconName = keyof typeof Ic`.

### Agent `logos` — `src/components/surf/logos.tsx`
Export `SymbolLogo({ className?: string })` and `FullLogo({ className?: string })`.
Render the brand SVGs via `<img>` pointing at `/images/logo-symbol.svg` and
`/images/logo-full.svg` (alt "Surfer VPN"), forwarding `className`. (The SVGs use
gradient ids; `<img>` avoids id collisions.)

### Agent `shell` — `src/components/surf/TgHeader.tsx`, `Hero.tsx`, `Waves.tsx`
- `TgHeader` → the `.tg-header`: Close icon btn (left), centered title
  (`.tg-title-main` "Surfer VPN" + `.tg-title-sub` "mini app"), Dots icon btn (right).
- `Waves` → the `.waves` SVG (3 layered paths `wv3/wv2/wv1`, exact `d` & fills).
- `Hero` → the `.hero`: `.hero-sky` (clouds c1/c2 + sun), `.hero-art` (CSS bg, just
  the empty div), `.hero-inner` → `.hero-copy` with `.brandmark`
  (`.brandmark-badge` wrapping `<SymbolLogo className="brandmark-symbol" />` +
  `.brandmark-name` "Surfer VPN"), `.hero-slogan` "Свобода<br/><span class='accent'>без границ</span>",
  `.hero-sub` "Быстрый и безопасный<br/>интернет для тебя" + `<Ic.Wave size={15} color="#2b97e6"/>`.
  Render `<Waves/>` at the hero bottom.
Import `Ic` from `@/components/surf/icons`, logos from `@/components/surf/logos`.

### Agent `usercard` — `src/components/surf/UserCard.tsx`
Props `{ user: User }`. Port `UserCard` + `StatusPill` + `InfoRow` + days banner
from `app.jsx`/CSS:
- `.user-top`: `.user-avatar` (first letter of `user.firstName`), `.user-name`
  (firstName), `.user-plan` (`${user.subscription.plan} · подписка`), `<StatusPill status={user.subscription.status}/>`.
- `.days-banner`: `.days-num` = `user.subscription.daysLeft`, text "дней до окончания",
  `.days-bar span` width = `min(100,(daysLeft/90)*100)%`.
- `.info-list`: rows — Статус (`<StatusPill/>` value), Действует до (`formatDate(expiresAt)`),
  Сервер (`${country}, ${city}`, `accent`), Telegram ID (`telegramId`). Icons:
  `Ic.ShieldCheck, Ic.Clock, Ic.Pin, Ic.Hash`.
`StatusPill` maps active→{Активна, ok}, trial→{Пробная, trial}, expired→{Истекла, off}.
Export `UserCard` (and you may export `StatusPill`). Use `formatDate` from `@/lib/utils`.

### Agent `actions` — `src/components/surf/Actions.tsx`, `Toast.tsx`, `InstallGrid.tsx`
- `Actions` props `{ vpnKey: string; onCopy: () => void }`: `.btn-primary`
  (`.btn-primary-shine`, label "Открыть Happ" / desc "Подключиться в один тап",
  `.btn-primary-arrow` `<Ic.Arrow/>`) → onClick opens `buildHappUrl(vpnKey)` via
  `openLink` + `haptic("medium")`. `.btn-ghost` (`<Ic.Copy/>` + "Скопировать ключ") → `onCopy`.
- `Toast` props `{ msg: string; show: boolean }`: the `.toast` with `.toast-check`
  `<Ic.ShieldCheck size={18} color="#fff"/>`.
- `InstallGrid`: `.install` head ("Установить приложение" / "Выбери свою платформу")
  + `.install-grid` 2×2 over PLATFORMS = `[ios(Phone,'iPhone и iPad'), android(Droid,'Телефон и планшет'), macos(Laptop,'Mac на Apple Silicon'), windows(Monitor,'ПК и ноутбук')]`
  using `<a class="install-card" href target=_blank>`: `.install-ic` glyph,
  `.install-name`, `.install-desc`, `.install-btn` (`<Ic.Download size={16}/>` + "Установить").
  (Use the local PLATFORMS array per the design, not `@/lib/constants`.)

### Agent `wiring` — `src/pages/HomePage.tsx`, `src/pages/NotFoundPage.tsx`, `src/main.tsx`, `src/App.tsx`
- `main.tsx`: keep `BrowserRouter` + `<App/>`. Import `"@/styles/surfer.css"` (after `./index.css`).
  Remove the sonner `<Toaster>` (design has its own toast). Keep `initTelegram()`/`applyTelegramTheme()`.
- `App.tsx`: routes `/` → `HomePage`, `*` → `NotFoundPage`.
- `HomePage.tsx`: reproduce `app.jsx`'s `App`: `.screen` → `<TgHeader/>` + `.scroll`
  ( `<Hero/>` + `.page` ( `<UserCard user={user}/>`, `<Actions vpnKey onCopy/>`,
  `<InstallGrid/>`, `.foot` "Surfer VPN · быстрый и свободный интернет" ) ) +
  `<Toast/>` + `<Splash hidden={!loading}/>`. State: splash `loading` true for 2100ms
  (also clears once `useUser` resolves), toast {show,msg}, copy handler →
  `navigator.clipboard.writeText(user.vpnKey)` then toast "Ключ успешно скопирован".
  Use `useUser()` for data; while user is null show nothing under the hero (splash covers it).
- `NotFoundPage.tsx`: simple on-brand 404 inside `.screen` using `<FullLogo/>` +
  "Страница не найдена" + a link back to `/` (react-router `Link`). No deps on old components.
- Also import `<Splash/>` from `@/components/surf/Splash` (built by the `splash` agent).

### Agent `splash` — `src/components/surf/Splash.tsx`
Props `{ hidden: boolean }`. Port `.splash`: `.splash-art` (CSS bg surfer image),
`.splash-veil`, `.splash-center` (`<FullLogo className="splash-logo"/>` +
`.splash-tag` "Свобода без границ"), `.splash-loader` (3 `<span>`). Toggle `.hide` when `hidden`.

## Integration assumptions

- Icons at `@/components/surf/icons` exporting `Ic`.
- Logos at `@/components/surf/logos` exporting `SymbolLogo`, `FullLogo`.
- All surf components under `@/components/surf/*`, named exports.
- All class names come from `src/styles/surfer.css` (the `styles` agent owns it).
- Strict TS must pass: `npm run build` (tsc -b + vite build) green.
