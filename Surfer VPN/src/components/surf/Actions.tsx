import { openHappWithKey, haptic } from '@/lib/telegram';
import { Ic } from '@/components/surf/icons';

export function Actions({
  vpnKey,
  onCopy,
  onRefresh,
  onHappCopied,
  refreshing,
}: {
  vpnKey: string;
  onCopy: () => void;
  onRefresh?: () => void;
  onHappCopied?: () => void;
  refreshing?: boolean;
}) {
  const openHapp = async () => {
    if (!vpnKey.trim()) return;
    haptic('medium');
    const copied = await openHappWithKey(vpnKey);
    if (copied) onHappCopied?.();
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
