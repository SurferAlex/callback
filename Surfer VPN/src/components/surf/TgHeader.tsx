import { Ic } from "@/components/surf/icons";
import { BRAND } from "@/lib/constants";
import { useAuth } from "@/contexts/AuthContext";
import { isTelegramMiniApp } from "@/lib/runtime";
import { getWebApp } from "@/lib/telegram";

export function TgHeader() {
  const { mode, logout } = useAuth();
  const inTg = isTelegramMiniApp();
  const sub = inTg ? "mini app" : "web cabinet";

  return (
    <div className="tg-header">
      {inTg ? (
        <button
          type="button"
          className="tg-iconbtn"
          aria-label="Закрыть"
          onClick={() => getWebApp()?.close()}
        >
          <Ic.Close />
        </button>
      ) : (
        <span className="tg-iconbtn tg-iconbtn-spacer" />
      )}
      <div className="tg-title">
        <span className="tg-title-main">{BRAND.name}</span>
        <span className="tg-title-sub">{sub}</span>
      </div>
      {mode === "web" ? (
        <button
          type="button"
          className="tg-logout"
          onClick={() => void logout()}
        >
          Выйти
        </button>
      ) : (
        <button type="button" className="tg-iconbtn" aria-label="Меню">
          <Ic.Dots />
        </button>
      )}
    </div>
  );
}
