import type { User } from "@/types";

/**
 * Mock user used when running outside Telegram (browser dev) or before the
 * backend is connected. The Telegram identity fields are overwritten at
 * runtime with real `initDataUnsafe.user` data when available
 * (see `lib/api.ts` -> `getCurrentUser`).
 */
export const MOCK_USER: User = {
  telegramId: 482917365,
  firstName: "Алекс",
  lastName: undefined,
  username: "alex_surfer",
  photoUrl: undefined,
  vpnKey:
    "vless://7c1f9a2e-4b8d-46a1-9e3c-1d2f5a6b8c90@nl.surfervpn.com:443?type=ws#Surfer",
  subscription: {
    status: "active",
    plan: "Premium",
    expiresAt: "2026-07-15T00:00:00.000Z",
    daysLeft: 43,
    autoRenew: true,
  },
  server: {
    id: "nl-ams-1",
    city: "Амстердам",
    countryCode: "NL",
    country: "Нидерланды",
    pingMs: 24,
  },
};

/** Compute whole days remaining from now until an ISO date (never negative). */
export function daysUntil(iso: string): number {
  const ms = new Date(iso).getTime() - Date.now();
  return Math.max(0, Math.ceil(ms / 86_400_000));
}
