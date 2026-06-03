import { useEffect, useReducer } from "react";

import type {
  ColorScheme,
  TelegramThemeParams,
  TelegramUser,
  TelegramWebApp,
} from "@/types";
import {
  applyTelegramTheme,
  getColorScheme,
  getStartParam,
  getTelegramUser,
  getThemeParams,
  getWebApp,
  isTelegram as isTelegramRuntime,
  onThemeChanged,
} from "@/lib/telegram";

export interface UseTelegramResult {
  webApp: TelegramWebApp | null;
  isTelegram: boolean;
  colorScheme: ColorScheme;
  themeParams: TelegramThemeParams;
  user: TelegramUser | null;
  startParam: string | null;
}

/**
 * Reactive access to the Telegram Mini App runtime.
 *
 * On mount we apply the native theme and subscribe to `themeChanged`; each
 * change re-applies the theme and forces a re-render so consumers always read
 * fresh `colorScheme` / `themeParams`. Fully safe outside Telegram (all
 * helpers no-op and the hook returns sensible defaults).
 */
export function useTelegram(): UseTelegramResult {
  // Lightweight force-update used to re-read the runtime after theme changes.
  const [, forceRender] = useReducer((tick: number) => tick + 1, 0);

  useEffect(() => {
    // Reflect the current Telegram theme into our CSS variables immediately.
    applyTelegramTheme();

    const handleThemeChanged = () => {
      applyTelegramTheme();
      forceRender();
    };

    const unsubscribe = onThemeChanged(handleThemeChanged);
    return unsubscribe;
  }, []);

  return {
    webApp: getWebApp(),
    isTelegram: isTelegramRuntime(),
    colorScheme: getColorScheme(),
    themeParams: getThemeParams(),
    user: getTelegramUser(),
    startParam: getStartParam(),
  };
}
