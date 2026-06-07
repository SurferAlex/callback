export type SubscriptionRemaining = {
  totalMs: number;
  days: number;
  hours: number;
  minutes: number;
  seconds: number;
  expired: boolean;
};

export function computeSubscriptionRemaining(expiresAt: string): SubscriptionRemaining {
  const end = new Date(expiresAt).getTime();
  if (Number.isNaN(end)) {
    return { totalMs: 0, days: 0, hours: 0, minutes: 0, seconds: 0, expired: true };
  }
  const totalMs = end - Date.now();
  if (totalMs <= 0) {
    return { totalMs: 0, days: 0, hours: 0, minutes: 0, seconds: 0, expired: true };
  }
  const totalSeconds = Math.floor(totalMs / 1000);
  const days = Math.floor(totalSeconds / 86400);
  const hours = Math.floor((totalSeconds % 86400) / 3600);
  const minutes = Math.floor((totalSeconds % 3600) / 60);
  const seconds = totalSeconds % 60;
  return { totalMs, days, hours, minutes, seconds, expired: false };
}

export function pad2(n: number): string {
  return String(n).padStart(2, "0");
}

export function subscriptionProgressPercent(
  startsAt: string | undefined,
  expiresAt: string,
  remainingMs: number
): number {
  const end = new Date(expiresAt).getTime();
  const start = startsAt ? new Date(startsAt).getTime() : end - 90 * 86400000;
  const total = end - start;
  if (!Number.isFinite(total) || total <= 0) return 0;
  return Math.min(100, Math.max(0, (remainingMs / total) * 100));
}
