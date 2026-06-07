import type { ReactNode } from "react";
import type { SubscriptionStatus, User } from "@/types";
import { useSubscriptionCountdown } from "@/hooks/use-subscription-countdown";
import { userAvatarLetter, userDisplayName } from "@/lib/user-display";
import { pad2, subscriptionProgressPercent } from "@/lib/subscription-time";
import { formatDate, pluralize } from "@/lib/utils";
import { Ic } from "@/components/surf/icons";

interface StatusPillProps {
  status: SubscriptionStatus;
}

const STATUS_MAP: Record<string, { label: string; cls: string }> = {
  active: { label: "Активна", cls: "ok" },
  trial: { label: "Пробная", cls: "trial" },
  expired: { label: "Истекла", cls: "off" },
  none: { label: "Нет подписки", cls: "off" },
};

export function StatusPill({ status }: StatusPillProps) {
  const s = STATUS_MAP[status] ?? STATUS_MAP.active;
  return (
    <span className={"status-pill " + s.cls}>
      <span className="status-dot"></span>
      {s.label}
    </span>
  );
}

interface InfoRowProps {
  icon: ReactNode;
  label: string;
  value: ReactNode;
  accent?: boolean;
}

function InfoRow({ icon, label, value, accent }: InfoRowProps) {
  return (
    <div className="info-row">
      <span className="info-ic">{icon}</span>
      <span className="info-label">{label}</span>
      <span className={"info-value" + (accent ? " accent" : "")}>{value}</span>
    </div>
  );
}

interface UserCardProps {
  user: User;
}

export function UserCard({ user }: UserCardProps) {
  const displayName = userDisplayName(user);
  const countdown = useSubscriptionCountdown(user.subscription.expiresAt);
  const progress = subscriptionProgressPercent(
    user.subscription.startsAt,
    user.subscription.expiresAt,
    countdown.totalMs
  );

  const countdownLine = countdown.expired
    ? "Подписка завершена"
    : `${pluralize(countdown.days, ["день", "дня", "дней"])} · ${pad2(countdown.hours)}:${pad2(countdown.minutes)}:${pad2(countdown.seconds)}`;

  return (
    <section className="card user-card">
      <div className="user-top">
        <div className="user-id">
          <span className="user-avatar">{userAvatarLetter(user)}</span>
          <div>
            <div className="user-name">{displayName}</div>
            <div className="user-plan">{user.subscription.plan} · подписка</div>
          </div>
        </div>
        <StatusPill status={user.subscription.status} />
      </div>

      <div className="days-banner">
        <div className="days-num">{countdown.expired ? 0 : countdown.days}</div>
        <div className="days-text">
          <div className="days-countdown">{countdownLine}</div>
          <div className="days-sub">до окончания подписки</div>
          <div className="days-bar">
            <span style={{ width: progress + "%" }}></span>
          </div>
        </div>
      </div>

      <div className="info-list">
        <InfoRow
          icon={<Ic.ShieldCheck />}
          label="Статус"
          value={<StatusPill status={user.subscription.status} />}
        />
        <InfoRow
          icon={<Ic.Clock />}
          label="Действует до"
          value={formatDate(user.subscription.expiresAt)}
        />
        <InfoRow
          icon={<Ic.Pin />}
          label="Сервер"
          value={`${user.server.country}, ${user.server.city}`}
          accent
        />
        <InfoRow
          icon={<Ic.Hash />}
          label="Telegram ID"
          value={user.telegramId}
        />
      </div>
    </section>
  );
}
