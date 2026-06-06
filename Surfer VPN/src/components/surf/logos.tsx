import { BRAND } from "@/lib/constants";

export function SymbolLogo({ className }: { className?: string }) {
  return <img src="/images/logo-symbol.svg" alt={BRAND.name} className={className} />;
}

export function FullLogo({ className }: { className?: string }) {
  return <img src="/images/logo-full.svg" alt={BRAND.name} className={className} />;
}
