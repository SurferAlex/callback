import { Ic } from '@/components/surf/icons';
import type { IconName } from '@/components/surf/icons';

const PLATFORMS: { key: string; name: string; desc: string; url: string; Icon: IconName }[] = [
  { key: 'ios',     name: 'iOS',     desc: 'iPhone и iPad',        url: 'https://apps.apple.com/ru/app/happ-proxy-utility-plus/id6746188973',     Icon: 'Phone'   },
  { key: 'android', name: 'Android', desc: 'Телефон и планшет',    url: 'https://play.google.com/store/apps/details?id=com.happproxy&hl=ru', Icon: 'Droid'   },
  { key: 'macos',   name: 'macOS',   desc: 'Mac на Apple Silicon', url: 'https://apps.apple.com/ru/app/happ-proxy-utility-plus/id6746188973',   Icon: 'Laptop'  },
  { key: 'windows', name: 'Windows', desc: 'ПК и ноутбук',         url: 'https://github.com/Happ-proxy/happ-desktop/releases/latest/download/setup-Happ.x64.exe', Icon: 'Monitor' },
];

export function InstallGrid() {
  return (
    <section className="install">
      <div className="section-head">
        <h2>Установить приложение</h2>
        <p>Выбери свою платформу</p>
      </div>
      <div className="install-grid">
        {PLATFORMS.map((p) => {
          const Glyph = Ic[p.Icon];
          return (
            <a key={p.key} className="install-card" href={p.url} target="_blank" rel="noreferrer">
              <span className="install-ic"><Glyph size={26} /></span>
              <span className="install-name">{p.name}</span>
              <span className="install-desc">{p.desc}</span>
              <span className="install-btn"><Ic.Download size={16} />Установить</span>
            </a>
          );
        })}
      </div>
    </section>
  );
}
