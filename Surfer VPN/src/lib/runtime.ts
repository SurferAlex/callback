import { getWebApp } from "@/lib/telegram";

/** True when running inside Telegram Mini App with valid initData. */
export function isTelegramMiniApp(): boolean {
  const wa = getWebApp();
  return Boolean(wa?.initData && wa.initData.length > 0);
}

/** Web cabinet in a normal browser (JWT + refresh cookie). */
export function isWebCabinet(): boolean {
  return !isTelegramMiniApp();
}
