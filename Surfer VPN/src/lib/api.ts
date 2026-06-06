import type { User } from "@/types";
import { MOCK_USER, daysUntil } from "@/lib/mock-data";
import { clearAccessToken, getAccessToken } from "@/lib/auth-store";
import { isTelegramMiniApp } from "@/lib/runtime";
import { readStoredTelegramProfile } from "@/lib/tg-profile";
import { getTelegramUser, getWebApp } from "@/lib/telegram";

export const API_BASE_URL =
  import.meta.env.VITE_API_BASE_URL ?? "";

const USE_MOCK =
  import.meta.env.VITE_USE_MOCK === "true" ||
  import.meta.env.VITE_USE_MOCK === "1";

type ApiUserMe = {
  telegramId: number;
  firstName: string;
  lastName?: string;
  username?: string;
  vpnKey: string;
  subscription: {
    status: string;
    plan: string;
    expiresAt: string;
    daysLeft: number;
    autoRenew: boolean;
  };
  server: {
    id: string;
    city: string;
    countryCode: string;
    country: string;
    pingMs?: number;
  };
};

function buildAuthHeaders(forceTma = false): HeadersInit {
  const headers: Record<string, string> = {};
  if (isTelegramMiniApp() || forceTma) {
    const initData = getWebApp()?.initData;
    if (initData) {
      headers.Authorization = `tma ${initData}`;
      return headers;
    }
  }
  const token = getAccessToken();
  if (token) {
    headers.Authorization = `Bearer ${token}`;
  }
  return headers;
}

let refreshPromise: Promise<boolean> | null = null;

async function tryRefreshAccess(): Promise<boolean> {
  if (isTelegramMiniApp()) return false;
  if (!refreshPromise) {
    refreshPromise = (async () => {
      const { refreshSession } = await import("@/lib/auth-api");
      const t = await refreshSession();
      return t !== null;
    })().finally(() => {
      refreshPromise = null;
    });
  }
  return refreshPromise;
}

export async function apiFetch<T>(path: string, init?: RequestInit): Promise<T> {
  const base = API_BASE_URL || "";
  const doFetch = (forceTma = false) =>
    fetch(`${base}${path}`, {
      ...init,
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
        ...buildAuthHeaders(forceTma),
        ...init?.headers,
      },
    });

  let res = await doFetch();
  if (res.status === 401 && isTelegramMiniApp()) {
    clearAccessToken();
    res = await doFetch(true);
  } else if (res.status === 401 && !isTelegramMiniApp()) {
    const ok = await tryRefreshAccess();
    if (ok) {
      res = await doFetch();
    }
  }

  const body = await res.json().catch(() => ({}));
  if (!res.ok) {
    const msg =
      typeof body === "object" && body && "error" in body
        ? String((body as { error: string }).error)
        : res.statusText;
    throw new Error(msg || `HTTP ${res.status}`);
  }
  return body as T;
}

function resolveFirstName(data: ApiUserMe): string {
  const tg = getTelegramUser();
  const stored = readStoredTelegramProfile();

  if (isTelegramMiniApp() && tg?.first_name?.trim()) {
    return tg.first_name.trim();
  }
  const apiFirst = data.firstName?.trim();
  if (apiFirst && apiFirst !== "Пользователь") return apiFirst;
  if (stored?.first_name?.trim()) return stored.first_name.trim();
  if (stored?.username?.trim()) return stored.username.trim();
  if (data.username?.trim()) return data.username.trim();
  return apiFirst || tg?.first_name?.trim() || "Пользователь";
}

function mapApiUser(data: ApiUserMe): User {
  const tg = getTelegramUser();
  const stored = readStoredTelegramProfile();
  const expiresAt = data.subscription.expiresAt || new Date().toISOString();
  return {
    telegramId: data.telegramId,
    firstName: resolveFirstName(data),
    lastName: data.lastName ?? tg?.last_name ?? stored?.last_name,
    username: data.username ?? tg?.username ?? stored?.username,
    photoUrl: tg?.photo_url,
    vpnKey: data.vpnKey || "",
    subscription: {
      status:
        data.subscription.status === "active"
          ? "active"
          : data.subscription.status === "trial"
            ? "trial"
            : data.subscription.status === "none"
              ? "none"
              : "expired",
      plan: data.subscription.plan || "—",
      expiresAt,
      daysLeft: data.subscription.daysLeft ?? daysUntil(expiresAt),
      autoRenew: data.subscription.autoRenew ?? false,
    },
    server: {
      id: data.server.id,
      city: data.server.city,
      countryCode: data.server.countryCode,
      country: data.server.country,
      pingMs: data.server.pingMs,
    },
  };
}

export async function getCurrentUser(): Promise<User> {
  if (USE_MOCK) {
    const tg = getTelegramUser();
    const base = MOCK_USER;
    return {
      ...base,
      telegramId: tg?.id ?? base.telegramId,
      firstName: tg?.first_name ?? base.firstName,
      lastName: tg?.last_name ?? base.lastName,
      username: tg?.username ?? base.username,
      photoUrl: tg?.photo_url ?? base.photoUrl,
      subscription: {
        ...base.subscription,
        daysLeft: daysUntil(base.subscription.expiresAt),
      },
    };
  }

  const data = await apiFetch<ApiUserMe>("/api/v1/user/me");
  return mapApiUser(data);
}

export async function getVpnKey(): Promise<string> {
  if (USE_MOCK) {
    return (await getCurrentUser()).vpnKey;
  }
  const data = await apiFetch<{ vlessUri: string }>("/api/v1/user/config");
  return data.vlessUri;
}

export async function refreshVpnKey(): Promise<string> {
  if (USE_MOCK) {
    return getVpnKey();
  }
  const data = await apiFetch<{ vlessUri: string }>("/api/v1/user/config/refresh", {
    method: "POST",
  });
  return data.vlessUri;
}
