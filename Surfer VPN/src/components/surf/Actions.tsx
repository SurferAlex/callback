import { openLink, haptic } from '@/lib/telegram';
import { buildHappUrl } from '@/lib/constants';
import { Ic } from '@/components/surf/icons';

export function Actions({
  vpnKey,
  onCopy,
  onRefresh,
  refreshing,
}: {
  vpnKey: string;
  onCopy: () => void;
  onRefresh?: () => void;
  refreshing?: boolean;
}) {
  const openHapp = () => {
    if (!vpnKey.trim()) return;
    haptic('medium');
    openLink(buildHappUrl(vpnKey));
  };
  return (
    <section className="actions">
      <button className="btn-primary" onClick={openHapp}>
        <span className="btn-primary-shine"></span>
        <span className="btn-primary-label">
          <span className="btn-primary-title">Открыть Happ</span>
          <span className="btn-primary-desc">Подключиться в один тап</span>
        </span>
        <span className="btn-primary-arrow"><Ic.Arrow /></span>
      </button>
      <button className="btn-ghost" onClick={onCopy}>
        <Ic.Copy /><span>Скопировать ключ</span>
      </button>
      {onRefresh && (
        <button className="btn-ghost" onClick={onRefresh} disabled={refreshing}>
          <Ic.Refresh /><span>{refreshing ? "Обновляем…" : "Обновить конфиг"}</span>
        </button>
      )}
    </section>
  );
}
