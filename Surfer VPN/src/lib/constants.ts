import type { Platform } from "@/types";

/** Brand strings — single source of truth for copy. */
export const BRAND = {
  name: "Surfer VPN",
  slogan: "Свобода без границ",
  subtitle: "Быстрый и безопасный интернет для тебя",
} as const;

/** Base URL for the subscription deep-link (Happ client). */
export const SUB_BASE_URL =
  import.meta.env.VITE_SUB_BASE_URL ?? "https://sub.surfervpn.com";

/** Build the "Open Happ" deep link for a given VPN key. */
export function buildHappUrl(vpnKey: string): string {
  return `${SUB_BASE_URL}/open?key=${encodeURIComponent(vpnKey)}`;
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
