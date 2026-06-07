import { useEffect, useState } from "react";

import {
  computeSubscriptionRemaining,
  type SubscriptionRemaining,
} from "@/lib/subscription-time";

export function useSubscriptionCountdown(expiresAt: string): SubscriptionRemaining {
  const [remaining, setRemaining] = useState(() =>
    computeSubscriptionRemaining(expiresAt)
  );

  useEffect(() => {
    setRemaining(computeSubscriptionRemaining(expiresAt));
    const id = window.setInterval(() => {
      setRemaining(computeSubscriptionRemaining(expiresAt));
    }, 1000);
    return () => window.clearInterval(id);
  }, [expiresAt]);

  return remaining;
}
