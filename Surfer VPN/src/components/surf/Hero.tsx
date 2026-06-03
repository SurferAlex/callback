import { Ic } from "@/components/surf/icons";
import { SymbolLogo } from "@/components/surf/logos";
import { Waves } from "@/components/surf/Waves";

export function Hero() {
  return (
    <header className="hero">
      <div className="hero-sky">
        <span className="cloud c1"></span>
        <span className="cloud c2"></span>
        <span className="sun"></span>
      </div>
      <div className="hero-art" aria-hidden="true"></div>
      <div className="hero-inner">
        <div className="hero-copy">
          <div className="brandmark">
            <span className="brandmark-badge">
              <SymbolLogo className="brandmark-symbol" />
            </span>
            <span className="brandmark-name">Surfer&nbsp;VPN</span>
          </div>
          <h1 className="hero-slogan">
            Свобода
            <br />
            <span className="accent">без границ</span>
          </h1>
          <p className="hero-sub">
            Быстрый и безопасный
            <br />
            интернет для тебя
            <span className="hero-sub-ic">
              <Ic.Wave size={15} color="#2b97e6" />
            </span>
          </p>
        </div>
      </div>
      <Waves />
    </header>
  );
}
