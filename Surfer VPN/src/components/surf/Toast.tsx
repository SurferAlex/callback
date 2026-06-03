import { Ic } from '@/components/surf/icons';

export function Toast({ msg, show }: { msg: string; show: boolean }) {
  return (
    <div className={'toast' + (show ? ' show' : '')} role="status">
      <span className="toast-check"><Ic.ShieldCheck size={18} color="#fff" /></span>
      {msg}
    </div>
  );
}
