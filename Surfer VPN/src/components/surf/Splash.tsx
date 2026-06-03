import { FullLogo } from "@/components/surf/logos";

export function Splash({ hidden }: { hidden: boolean }) {
  return (
    <div className={"splash" + (hidden ? " hide" : "")} aria-hidden={hidden}>
      <div className="splash-art"></div>
      <div className="splash-veil"></div>
      <div className="splash-center">
        <FullLogo className="splash-logo" />
        <p className="splash-tag">Свобода без границ</p>
      </div>
      <div className="splash-loader">
        <span></span>
        <span></span>
        <span></span>
      </div>
    </div>
  );
}
