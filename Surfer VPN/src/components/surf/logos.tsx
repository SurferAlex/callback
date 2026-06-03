// Surfer VPN — brand logos rendered as <img> to avoid SVG gradient-id collisions.

export function SymbolLogo({ className }: { className?: string }) {
  return <img src="/images/logo-symbol.svg" alt="Surfer VPN" className={className} />;
}

export function FullLogo({ className }: { className?: string }) {
  return <img src="/images/logo-full.svg" alt="Surfer VPN" className={className} />;
}
