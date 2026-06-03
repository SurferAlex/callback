import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";

/** Tailwind-aware className combiner used by every shadcn/ui component. */
export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

/** Format an ISO date string as a localized Russian date, e.g. "5 июля 2026". */
export function formatDate(iso: string): string {
  try {
    return new Intl.DateTimeFormat("ru-RU", {
      day: "numeric",
      month: "long",
      year: "numeric",
    }).format(new Date(iso));
  } catch {
    return iso;
  }
}

/** Russian plural helper: pluralize(5, ["день","дня","дней"]) -> "дней". */
export function pluralize(
  count: number,
  forms: [one: string, few: string, many: string]
): string {
  const n = Math.abs(count) % 100;
  const n1 = n % 10;
  if (n > 10 && n < 20) return forms[2];
  if (n1 > 1 && n1 < 5) return forms[1];
  if (n1 === 1) return forms[0];
  return forms[2];
}

/** Mask a secret key for display, keeping the first/last few chars visible. */
export function maskKey(key: string, visible = 4): string {
  if (key.length <= visible * 2) return key;
  return `${key.slice(0, visible)}••••${key.slice(-visible)}`;
}
