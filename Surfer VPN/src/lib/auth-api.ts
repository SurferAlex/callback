import { API_BASE_URL } from "@/lib/api";
import { clearAccessToken, setAccessToken } from "@/lib/auth-store";
import { getWebApp } from "@/lib/telegram";

type TokenResponse = {
  accessToken: string;
  expiresIn: number;
};

async function parseTokenResponse(res: Response): Promise<TokenResponse> {
  const body = await res.json().catch(() => ({}));
  if (!res.ok) {
    const msg =
      typeof body === "object" && body && "error" in body
        ? String((body as { error: string }).error)
        : res.statusText;
    throw new Error(msg || `HTTP ${res.status}`);
  }
  return body as TokenResponse;
}

function applyToken(data: TokenResponse) {
  setAccessToken(data.accessToken);
  return data;
}

/** Mini App: exchange initData for JWT + HttpOnly refresh cookie. */
export async function sessionFromTelegramWebApp(): Promise<TokenResponse> {
  const initData = getWebApp()?.initData;
  if (!initData) {
    throw new Error("no telegram init data");
  }
  const res = await fetch(`${API_BASE_URL}/api/v1/auth/session/webapp`, {
    method: "POST",
    credentials: "include",
    headers: {
      Authorization: `tma ${initData}`,
    },
  });
  return applyToken(await parseTokenResponse(res));
}

/** Browser: Telegram Login Widget → JWT + refresh cookie. */
export async function sessionFromTelegramWidget(
  widgetUser: Record<string, unknown>
): Promise<TokenResponse> {
  const res = await fetch(`${API_BASE_URL}/api/v1/auth/session/widget`, {
    method: "POST",
    credentials: "include",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(widgetUser),
  });
  return applyToken(await parseTokenResponse(res));
}

/** Browser: restore session from HttpOnly refresh cookie. */
export async function refreshSession(): Promise<TokenResponse | null> {
  const res = await fetch(`${API_BASE_URL}/api/v1/auth/refresh`, {
    method: "POST",
    credentials: "include",
  });
  if (res.status === 401) {
    clearAccessToken();
    return null;
  }
  return applyToken(await parseTokenResponse(res));
}

export async function logoutSession(): Promise<void> {
  clearAccessToken();
  await fetch(`${API_BASE_URL}/api/v1/auth/logout`, {
    method: "POST",
    credentials: "include",
  }).catch(() => {});
}
