import type { User } from "@/types";
import { isTelegramMiniApp } from "@/lib/runtime";
import { getTelegramUser } from "@/lib/telegram";
import { readStoredTelegramProfile } from "@/lib/tg-profile";

/** Display name for the user card — matches Mini App Telegram identity in both modes. */
export function userDisplayName(user: Pick<User, "firstName" | "username">): string {
  const tg = getTelegramUser();
  const stored = readStoredTelegramProfile();

  if (isTelegramMiniApp() && tg?.first_name?.trim()) {
    return tg.first_name.trim();
  }

  const first = user.firstName?.trim();
  if (first && first !== "Пользователь") return first;

  const fromStored = stored?.first_name?.trim() || stored?.username?.trim();
  if (fromStored) return fromStored;

  const username = user.username?.trim();
  if (username) return username;

  return first || tg?.first_name?.trim() || "Пользователь";
}

export function userAvatarLetter(user: Pick<User, "firstName" | "username">): string {
  const name = userDisplayName(user);
  return name ? name[0].toUpperCase() : "?";
}
