import { useCallback, useEffect, useRef, useState } from "react";

import type { User } from "@/types";
import { getCurrentUser } from "@/lib/api";

export interface UseUserResult {
  user: User | null;
  loading: boolean;
  error: Error | null;
  refetch: () => void;
}

/**
 * Loads the current user via `getCurrentUser()` on mount.
 *
 * Manages `loading` / `error` state and exposes `refetch` to re-run the
 * request. Guards against setting state after the component has unmounted (or
 * after a newer request has superseded an in-flight one).
 */
export function useUser(): UseUserResult {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<Error | null>(null);
  // Bumped by `refetch` to re-trigger the effect.
  const [requestId, setRequestId] = useState<number>(0);
  const silentRef = useRef(false);

  const refetch = useCallback(() => {
    silentRef.current = user !== null;
    setRequestId((id) => id + 1);
  }, [user]);

  useEffect(() => {
    let active = true;
    const silent = silentRef.current;
    silentRef.current = false;

    if (!silent) {
      setLoading(true);
      setError(null);
    }

    getCurrentUser()
      .then((result) => {
        if (!active) return;
        setUser(result);
        setError(null);
      })
      .catch((err: unknown) => {
        if (!active) return;
        if (!silent) {
          setError(err instanceof Error ? err : new Error(String(err)));
        }
      })
      .finally(() => {
        if (!active) return;
        if (!silent) {
          setLoading(false);
        }
      });

    return () => {
      active = false;
    };
  }, [requestId]);

  return { user, loading, error, refetch };
}
