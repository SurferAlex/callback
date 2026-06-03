import { Ic } from "@/components/surf/icons";

export function TgHeader() {
  return (
    <div className="tg-header">
      <button className="tg-iconbtn" aria-label="Закрыть">
        <Ic.Close />
      </button>
      <div className="tg-title">
        <span className="tg-title-main">Surfer VPN</span>
        <span className="tg-title-sub">mini app</span>
      </div>
      <button className="tg-iconbtn" aria-label="Меню">
        <Ic.Dots />
      </button>
    </div>
  );
}
