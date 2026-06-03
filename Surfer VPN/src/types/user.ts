/** Subscription lifecycle status. */
export type SubscriptionStatus = "active" | "expired" | "trial" | "none";

/** A user's VPN subscription details. */
export interface Subscription {
  status: SubscriptionStatus;
  /** Human-readable plan name, e.g. "Premium". */
  plan: string;
  /** ISO-8601 date string of when the subscription ends. */
  expiresAt: string;
  /** Number of whole days remaining until expiration (>= 0). */
  daysLeft: number;
  /** Whether the subscription auto-renews. */
  autoRenew: boolean;
}

/** Currently connected/assigned VPN server. */
export interface VpnServer {
  /** Stable server id, e.g. "nl-ams-1". */
  id: string;
  /** Display name, e.g. "Amsterdam". */
  city: string;
  /** ISO country code, e.g. "NL". */
  countryCode: string;
  /** Country display name, e.g. "Нидерланды". */
  country: string;
  /** Optional measured latency in ms. */
  pingMs?: number;
}

/**
 * Application-level user model.
 * Combines Telegram identity with VPN/subscription data.
 */
export interface User {
  /** Telegram user id. */
  telegramId: number;
  /** Telegram first name. */
  firstName: string;
  /** Telegram last name (optional). */
  lastName?: string;
  /** Telegram @username (without the @, optional). */
  username?: string;
  /** Telegram avatar URL (optional). */
  photoUrl?: string;
  /** Personal VPN access key used to open the client / copy to clipboard. */
  vpnKey: string;
  /** Subscription details. */
  subscription: Subscription;
  /** Currently assigned server. */
  server: VpnServer;
}
