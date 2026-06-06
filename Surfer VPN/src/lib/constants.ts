import type { Platform } from "@/types";

/** Brand strings — single source of truth for copy. */
export const BRAND = {
  name: "Surf VPN",
  slogan: "Свобода без границ",
  subtitle: "Быстрый и безопасный интернет для тебя",
} as const;

/** API base URL (Happ bridge for Mini App). */
export const API_BASE_URL =
  import.meta.env.VITE_API_BASE_URL ?? "https://api.surfwave.space";

/** Base URL for subscription pages (optional web fallback). */
export const SUB_BASE_URL =
  import.meta.env.VITE_SUB_BASE_URL ?? "https://sub.surfervpn.com";

/** Build Happ deep link: happ://add/vless://… (see 3x-ui / Happ docs). */
export function buildHappUrl(vpnKey: string): string {
  const key = vpnKey.trim();
  if (!key) return "";
  if (key.startsWith("happ://")) return key;
  return `happ://add/${key}`;
}

/** HTTPS bridge on API → 302 to happ://add/… (for Telegram Mini App). */
export function buildHappApiOpenUrl(vpnKey: string): string {
  const key = vpnKey.trim();
  if (!key) return "";
  const base = API_BASE_URL.replace(/\/$/, "");
  return `${base}/api/v1/happ/open?key=${encodeURIComponent(key)}`;
}

/** Install grid data — mock links for now. */
export const PLATFORMS: Platform[] = [
  {
    id: "macos",
    name: "macOS",
    description: "Приложение для Mac",
    url: "https://example.com/macos",
  },
  {
    id: "ios",
    name: "iOS",
    description: "Приложение для iPhone",
    url: "https://example.com/ios",
  },
  {
    id: "android",
    name: "Android",
    description: "Приложение для Android",
    url: "https://example.com/android",
  },
  {
    id: "windows",
    name: "Windows",
    description: "Приложение для Windows",
    url: "https://example.com/windows",
  },
];
