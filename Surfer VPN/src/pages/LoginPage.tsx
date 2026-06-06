import { useEffect, useRef, useState } from "react";
import { useLocation, useNavigate } from "react-router-dom";

import { Waves } from "@/components/surf/Waves";
import { useAuth } from "@/contexts/AuthContext";
import { sessionFromTelegramWidget } from "@/lib/auth-api";
import { BRAND } from "@/lib/constants";

const BOT_USERNAME = import.meta.env.VITE_TELEGRAM_BOT_USERNAME ?? "";

declare global {
  interface Window {
    onTelegramAuth?: (user: Record<string, unknown>) => void;
  }
}

export function LoginPage() {
  const navigate = useNavigate();
  const location = useLocation();
  const { authenticated, bootstrap } = useAuth();
  const widgetRef = useRef<HTMLDivElement>(null);
  const [loginError, setLoginError] = useState<string | null>(null);

  const redirectTo = (location.state as { from?: string } | null)?.from ?? "/";

  useEffect(() => {
    if (authenticated) {
      navigate(redirectTo, { replace: true });
    }
  }, [authenticated, navigate, redirectTo]);

  useEffect(() => {
    if (!BOT_USERNAME || !widgetRef.current) return;

    window.onTelegramAuth = async (user) => {
      setLoginError(null);
      try {
        await sessionFromTelegramWidget(user);
        await bootstrap();
        navigate(redirectTo, { replace: true });
      } catch (err) {
        console.error("[auth] widget login failed:", err);
        setLoginError("Не удалось войти. Попробуйте ещё раз.");
      }
    };

    widgetRef.current.innerHTML = "";
    const script = document.createElement("script");
    script.src = "https://telegram.org/js/telegram-widget.js?22";
    script.async = true;
    script.setAttribute("data-telegram-login", BOT_USERNAME);
    script.setAttribute("data-size", "large");
    script.setAttribute("data-radius", "12");
    script.setAttribute("data-onauth", "onTelegramAuth(user)");
    script.setAttribute("data-request-access", "write");
    widgetRef.current.appendChild(script);

    return () => {
      delete window.onTelegramAuth;
    };
  }, [bootstrap, navigate, redirectTo]);

  return (
    <div className="screen login-screen">
      <Waves />
      <main className="login-card page">
        <h1 className="login-title">🏄 {BRAND.name}</h1>
        <p className="login-sub">
          Войдите через Telegram, чтобы открыть личный кабинет в браузере.
        </p>
        <p className="login-hint">
          После первого входа в Mini App сессия сохранится — повторный вход не
          понадобится.
        </p>
        {BOT_USERNAME ? (
          <div ref={widgetRef} className="login-widget" />
        ) : (
          <p className="login-error">
            Укажите VITE_TELEGRAM_BOT_USERNAME в настройках сборки.
          </p>
        )}
        {loginError && <p className="login-error">{loginError}</p>}
      </main>
    </div>
  );
}
