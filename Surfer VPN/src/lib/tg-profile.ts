const STORAGE_KEY = "surf_tg_profile";

export type StoredTelegramProfile = {
  first_name?: string;
  last_name?: string;
  username?: string;
};

function canUseStorage(): boolean {
  return typeof localStorage !== "undefined";
}

export function saveTelegramProfile(user: Record<string, unknown>): void {
  if (!canUseStorage()) return;
  const profile: StoredTelegramProfile = {
    first_name:
      typeof user.first_name === "string" ? user.first_name : undefined,
    last_name: typeof user.last_name === "string" ? user.last_name : undefined,
    username: typeof user.username === "string" ? user.username : undefined,
  };
  if (!profile.first_name && !profile.username) return;
  localStorage.setItem(STORAGE_KEY, JSON.stringify(profile));
}

export function saveTelegramProfileFields(profile: StoredTelegramProfile): void {
  if (!canUseStorage()) return;
  if (!profile.first_name?.trim() && !profile.username?.trim()) return;
  localStorage.setItem(STORAGE_KEY, JSON.stringify(profile));
}

export function readStoredTelegramProfile(): StoredTelegramProfile | null {
  if (!canUseStorage()) return null;
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    if (!raw) return null;
    return JSON.parse(raw) as StoredTelegramProfile;
  } catch {
    return null;
  }
}

export function clearTelegramProfile(): void {
  if (!canUseStorage()) return;
  localStorage.removeItem(STORAGE_KEY);
}
