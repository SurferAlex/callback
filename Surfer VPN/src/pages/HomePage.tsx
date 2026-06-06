import { useEffect, useRef, useState } from "react";

import { TgHeader } from "@/components/surf/TgHeader";
import { Hero } from "@/components/surf/Hero";
import { UserCard } from "@/components/surf/UserCard";
import { Actions } from "@/components/surf/Actions";
import { InstallGrid } from "@/components/surf/InstallGrid";
import { Toast } from "@/components/surf/Toast";
import { Splash } from "@/components/surf/Splash";
import { useUser } from "@/hooks";
import { BRAND } from "@/lib/constants";
import { refreshVpnKey } from "@/lib/api";
import { isTelegramMiniApp } from "@/lib/runtime";

interface ToastState {
  show: boolean;
  msg: string;
}

/** Surf VPN home screen. Mirrors the design's `App` component. */
export function HomePage() {
  const { user, loading: dataLoading, error: userError, refetch } = useUser();
  const [vpnKeyOverride, setVpnKeyOverride] = useState<string | null>(null);
  const [refreshing, setRefreshing] = useState(false);

  const [toast, setToast] = useState<ToastState>({ show: false, msg: "" });
  // Splash stays up until BOTH the 2.1s timer fired AND the data resolved.
  const [minElapsed, setMinElapsed] = useState(false);
  const toastTimer = useRef<ReturnType<typeof setTimeout> | null>(null);

  useEffect(() => {
    if (isTelegramMiniApp()) {
      setMinElapsed(true);
      return;
    }
    const t = setTimeout(() => setMinElapsed(true), 2100);
    return () => clearTimeout(t);
  }, []);

  useEffect(() => {
    return () => {
      if (toastTimer.current) clearTimeout(toastTimer.current);
    };
  }, []);

  const fireToast = (msg: string) => {
    setToast({ show: true, msg });
    if (toastTimer.current) clearTimeout(toastTimer.current);
    toastTimer.current = setTimeout(
      () => setToast((t) => ({ ...t, show: false })),
      2200
    );
  };

  const displayKey = vpnKeyOverride ?? user?.vpnKey ?? "";

  const copyKey = async () => {
    if (!displayKey) return;
    try {
      await navigator.clipboard.writeText(displayKey);
    } catch {
      /* clipboard may be blocked in sandbox */
    }
    fireToast("Ключ успешно скопирован");
  };

  // Splash only on first load (not on silent refetch after «Обновить конфиг»).
  const loading = (!minElapsed || (dataLoading && !user));

  return (
    <div className="screen">
      <TgHeader />
      <div className="scroll">
        <Hero />
        <main className="page">
          {userError && !user && (
            <p className="login-error" style={{ margin: "0 0 16px" }}>
              Не удалось загрузить профиль. Попробуйте обновить страницу.
            </p>
          )}
          {user && <UserCard user={user} />}
          {user && (
            <Actions
              vpnKey={displayKey}
              onCopy={copyKey}
              refreshing={refreshing}
              onRefresh={async () => {
                if (refreshing) return;
                setRefreshing(true);
                try {
                  const key = await refreshVpnKey();
                  setVpnKeyOverride(key);
                  refetch();
                  fireToast("Конфиг обновлён");
                } catch (err) {
                  const msg =
                    err instanceof Error && err.message
                      ? err.message
                      : "Не удалось обновить конфиг";
                  fireToast(msg);
                } finally {
                  setRefreshing(false);
                }
              }}
            />
          )}
          <InstallGrid />
          <footer className="foot">
            {BRAND.name} · быстрый и свободный интернет
          </footer>
        </main>
      </div>
      <Toast msg={toast.msg} show={toast.show} />
      <Splash hidden={!loading} />
    </div>
  );
}
