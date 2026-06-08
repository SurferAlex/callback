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
  import.meta.env.VITE_SUB_BASE_URL ?? "https://sub.surfwave.space";

/** Build Happ deep link for VLESS key or subscription HTTPS URL. */
export function buildHappUrl(configOrSubUrl: string): string {
  const value = configOrSubUrl.trim();
  if (!value) return "";
  if (value.startsWith("happ://")) return value;
  return `happ://add/${value}`;
}

/** HTTPS bridge on API → 302 to happ://add/… (for Telegram Mini App). */
export function buildHappApiOpenUrl(configOrSubUrl: string): string {
  const value = configOrSubUrl.trim();
  if (!value) return "";
  const base = API_BASE_URL.replace(/\/$/, "");
  return `${base}/api/v1/happ/open?key=${encodeURIComponent(value)}`;
}

/** Prefer subscription URL; fall back to legacy VLESS key. */
export function userConfigLink(user: {
  subscriptionUrl?: string;
  vpnKey: string;
}): string {
  const sub = user.subscriptionUrl?.trim();
  if (sub) return sub;
  return user.vpnKey.trim();
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
