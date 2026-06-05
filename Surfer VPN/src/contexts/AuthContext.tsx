import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
  type ReactNode,
} from "react";

import {
  logoutSession,
  refreshSession,
  sessionFromTelegramWebApp,
} from "@/lib/auth-api";
import { isTelegramMiniApp, isWebCabinet } from "@/lib/runtime";

type AuthState = {
  ready: boolean;
  authenticated: boolean;
  mode: "telegram" | "web";
  bootstrap: () => Promise<boolean>;
  logout: () => Promise<void>;
};

const AuthContext = createContext<AuthState | null>(null);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [ready, setReady] = useState(false);
  const [authenticated, setAuthenticated] = useState(false);
  const mode: AuthState["mode"] = isTelegramMiniApp() ? "telegram" : "web";

  const bootstrap = useCallback(async (): Promise<boolean> => {
    if (import.meta.env.VITE_USE_MOCK === "true" || import.meta.env.VITE_USE_MOCK === "1") {
      setAuthenticated(true);
      setReady(true);
      return true;
    }

    try {
      if (isTelegramMiniApp()) {
        await sessionFromTelegramWebApp();
        setAuthenticated(true);
        return true;
      }
      if (isWebCabinet()) {
        const t = await refreshSession();
        const ok = t !== null;
        setAuthenticated(ok);
        return ok;
      }
      setAuthenticated(false);
      return false;
    } catch (err) {
      console.error("[auth] bootstrap failed:", err);
      setAuthenticated(false);
      return false;
    } finally {
      setReady(true);
    }
  }, []);

  const logout = useCallback(async () => {
    await logoutSession();
    setAuthenticated(false);
  }, []);

  useEffect(() => {
    void bootstrap();
  }, [bootstrap]);

  const value = useMemo(
    () => ({ ready, authenticated, mode, bootstrap, logout }),
    [ready, authenticated, mode, bootstrap, logout]
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth(): AuthState {
  const ctx = useContext(AuthContext);
  if (!ctx) {
    throw new Error("useAuth must be used within AuthProvider");
  }
  return ctx;
}
