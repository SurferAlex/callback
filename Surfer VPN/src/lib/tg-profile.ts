const STORAGE_KEY = "surf_tg_profile";

export type StoredTelegramProfile = {
  first_name?: string;
  last_name?: string;
  username?: string;
};

export function saveTelegramProfile(user: Record<string, unknown>): void {
  if (typeof sessionStorage === "undefined") return;
  const profile: StoredTelegramProfile = {
    first_name:
      typeof user.first_name === "string" ? user.first_name : undefined,
    last_name: typeof user.last_name === "string" ? user.last_name : undefined,
    username: typeof user.username === "string" ? user.username : undefined,
  };
  sessionStorage.setItem(STORAGE_KEY, JSON.stringify(profile));
}

export function readStoredTelegramProfile(): StoredTelegramProfile | null {
  if (typeof sessionStorage === "undefined") return null;
  try {
    const raw = sessionStorage.getItem(STORAGE_KEY);
    if (!raw) return null;
    return JSON.parse(raw) as StoredTelegramProfile;
  } catch {
    return null;
  }
}
