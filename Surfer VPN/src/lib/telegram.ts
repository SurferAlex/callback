import type {
  ColorScheme,
  TelegramThemeParams,
  TelegramUser,
  TelegramWebApp,
} from "@/types";

/**
 * Safe accessor for the Telegram WebApp runtime.
 * Returns `null` when running outside Telegram (e.g. local browser dev).
 */
export function getWebApp(): TelegramWebApp | null {
  if (typeof window === "undefined") return null;
  return window.Telegram?.WebApp ?? null;
}

/** True when the app is actually running inside the Telegram client. */
export function isTelegram(): boolean {
  const wa = getWebApp();
  return Boolean(wa && wa.initData && wa.initData.length > 0);
}

/**
 * Initialize the Mini App: mark ready, expand to full height, and set the
 * header/background colors to match our ocean theme. No-ops outside Telegram.
 */
export function initTelegram(): void {
  const wa = getWebApp();
  if (!wa) return;
  try {
    wa.ready();
    wa.expand();
    // Light ocean header to blend with the hero gradient.
    wa.setHeaderColor?.("#e0f2fe");
    wa.setBackgroundColor?.("#f0f9ff");
  } catch {
    /* defensive: older clients may lack some setters */
  }
}

/** Read the Telegram user from launch params (null if unavailable). */
export function getTelegramUser(): TelegramUser | null {
  return getWebApp()?.initDataUnsafe.user ?? null;
}

/** The `start_param` deep-link payload, if the app was opened with one. */
export function getStartParam(): string | null {
  return getWebApp()?.initDataUnsafe.start_param ?? null;
}

/** Current Telegram color scheme; defaults to "light". */
export function getColorScheme(): ColorScheme {
  return getWebApp()?.colorScheme ?? "light";
}

/** Current Telegram theme params (empty object outside Telegram). */
export function getThemeParams(): TelegramThemeParams {
  return getWebApp()?.themeParams ?? {};
}

/**
 * Bridge Telegram theme params into our CSS custom properties so native
 * theming can influence the UI. Called on init and on `themeChanged`.
 */
export function applyTelegramTheme(): void {
  const wa = getWebApp();
  if (!wa || typeof document === "undefined") return;
  const tp = wa.themeParams;
  const root = document.documentElement;
  const set = (name: string, value?: string) => {
    if (value) root.style.setProperty(name, value);
  };
  set("--tg-bg", tp.bg_color);
  set("--tg-text", tp.text_color);
  set("--tg-hint", tp.hint_color);
  set("--tg-link", tp.link_color);
  set("--tg-button", tp.button_color);
  set("--tg-button-text", tp.button_text_color);
}

/** Subscribe to Telegram theme changes. Returns an unsubscribe fn. */
export function onThemeChanged(handler: () => void): () => void {
  const wa = getWebApp();
  if (!wa) return () => {};
  wa.onEvent("themeChanged", handler);
  return () => wa.offEvent("themeChanged", handler);
}

/** Open an external link using Telegram's native handler when available. */
export function openLink(url: string): void {
  const wa = getWebApp();
  if (wa) {
    wa.openLink(url);
  } else if (typeof window !== "undefined") {
    window.open(url, "_blank", "noopener,noreferrer");
  }
}

/** Fire light haptic feedback (no-op outside Telegram). */
export function haptic(
  type: "light" | "medium" | "heavy" | "success" | "error" | "warning" = "light"
): void {
  const hf = getWebApp()?.HapticFeedback;
  if (!hf) return;
  try {
    if (type === "success" || type === "error" || type === "warning") {
      hf.notificationOccurred(type);
    } else {
      hf.impactOccurred(type);
    }
  } catch {
    /* ignore */
  }
}
