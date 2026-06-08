// Surf VPN — icon set. Simple geometric glyphs only (no brand logos).

type IconProps = { size?: number; color?: string };

export const Ic = {
  // small wave / surf mark used in the logo
  Wave: ({ size = 28, color = '#fff' }: IconProps) => (
    <svg width={size} height={size} viewBox="0 0 32 32" fill="none">
      <path
        d="M3 21c3.2 0 3.2-4 6.4-4s3.2 4 6.4 4 3.2-4 6.4-4 3.2 4 6.4 4"
        stroke={color}
        strokeWidth="2.6"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
      <path
        d="M3 26c3.2 0 3.2-4 6.4-4s3.2 4 6.4 4 3.2-4 6.4-4 3.2 4 6.4 4"
        stroke={color}
        strokeWidth="2.6"
        strokeLinecap="round"
        strokeLinejoin="round"
        opacity="0.5"
      />
      <circle cx="24" cy="9" r="4" fill={color} opacity="0.9" />
    </svg>
  ),
  ShieldCheck: ({ size = 18, color = '#16a34a' }: IconProps) => (
    <svg width={size} height={size} viewBox="0 0 24 24" fill="none">
      <path
        d="M12 2.5l7 2.7v5.6c0 4.6-3 8.4-7 10.2-4-1.8-7-5.6-7-10.2V5.2l7-2.7z"
        fill={color}
        fillOpacity="0.14"
        stroke={color}
        strokeWidth="1.7"
        strokeLinejoin="round"
      />
      <path
        d="M8.6 12.2l2.3 2.3 4.4-4.7"
        stroke={color}
        strokeWidth="1.9"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
    </svg>
  ),
  Copy: ({ size = 20, color = '#2563EB' }: IconProps) => (
    <svg width={size} height={size} viewBox="0 0 24 24" fill="none">
      <rect x="8" y="8" width="12" height="13" rx="3" stroke={color} strokeWidth="1.9" />
      <path
        d="M16 8V6a3 3 0 00-3-3H7a3 3 0 00-3 3v8a3 3 0 003 3h1"
        stroke={color}
        strokeWidth="1.9"
        strokeLinecap="round"
      />
    </svg>
  ),
  Pin: ({ size = 16, color = '#2563EB' }: IconProps) => (
    <svg width={size} height={size} viewBox="0 0 24 24" fill="none">
      <path
        d="M12 21s7-5.6 7-11a7 7 0 10-14 0c0 5.4 7 11 7 11z"
        fill={color}
        fillOpacity="0.14"
        stroke={color}
        strokeWidth="1.8"
        strokeLinejoin="round"
      />
      <circle cx="12" cy="10" r="2.6" fill={color} />
    </svg>
  ),
  Clock: ({ size = 16, color = '#0ea5b7' }: IconProps) => (
    <svg width={size} height={size} viewBox="0 0 24 24" fill="none">
      <circle cx="12" cy="12" r="8.5" stroke={color} strokeWidth="1.8" />
      <path
        d="M12 7.5V12l3 1.8"
        stroke={color}
        strokeWidth="1.8"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
    </svg>
  ),
  Hash: ({ size = 16, color = '#64748b' }: IconProps) => (
    <svg width={size} height={size} viewBox="0 0 24 24" fill="none">
      <path
        d="M9 4L7 20M17 4l-2 16M4 9h16M3 15h16"
        stroke={color}
        strokeWidth="1.8"
        strokeLinecap="round"
      />
    </svg>
  ),
  Arrow: ({ size = 22, color = '#fff' }: IconProps) => (
    <svg width={size} height={size} viewBox="0 0 24 24" fill="none">
      <path
        d="M5 12h13M13 6l6 6-6 6"
        stroke={color}
        strokeWidth="2.4"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
    </svg>
  ),
  Close: ({ size = 18, color = '#0B2A45' }: IconProps) => (
    <svg width={size} height={size} viewBox="0 0 24 24" fill="none">
      <path d="M6 6l12 12M18 6L6 18" stroke={color} strokeWidth="2.2" strokeLinecap="round" />
    </svg>
  ),
  Dots: ({ size = 18, color = '#0B2A45' }: IconProps) => (
    <svg width={size} height={size} viewBox="0 0 24 24">
      <circle cx="5" cy="12" r="2" fill={color} />
      <circle cx="12" cy="12" r="2" fill={color} />
      <circle cx="19" cy="12" r="2" fill={color} />
    </svg>
  ),
  Download: ({ size = 18, color = '#2563EB' }: IconProps) => (
    <svg width={size} height={size} viewBox="0 0 24 24" fill="none">
      <path
        d="M12 3v11m0 0l-4-4m4 4l4-4"
        stroke={color}
        strokeWidth="2"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
      <path d="M5 18.5h14" stroke={color} strokeWidth="2" strokeLinecap="round" />
    </svg>
  ),
  // platform glyphs — neutral device silhouettes (no brand marks)
  Phone: ({ size = 26, color = '#2563EB' }: IconProps) => (
    <svg width={size} height={size} viewBox="0 0 24 24" fill="none">
      <rect x="6.5" y="2.5" width="11" height="19" rx="3" stroke={color} strokeWidth="1.8" />
      <path d="M10.5 5.2h3" stroke={color} strokeWidth="1.6" strokeLinecap="round" />
      <path d="M11 18.6h2" stroke={color} strokeWidth="1.8" strokeLinecap="round" />
    </svg>
  ),
  Droid: ({ size = 26, color = '#2563EB' }: IconProps) => (
    <svg width={size} height={size} viewBox="0 0 24 24" fill="none">
      <rect x="6" y="9" width="12" height="9" rx="2.4" stroke={color} strokeWidth="1.8" />
      <path
        d="M7.5 9c0-2.5 2-4.5 4.5-4.5S16.5 6.5 16.5 9"
        stroke={color}
        strokeWidth="1.8"
        strokeLinejoin="round"
      />
      <path d="M8 4.5l1.4 2M16 4.5l-1.4 2" stroke={color} strokeWidth="1.6" strokeLinecap="round" />
      <circle cx="10" cy="12" r="0.9" fill={color} />
      <circle cx="14" cy="12" r="0.9" fill={color} />
      <path d="M8 18v2M16 18v2" stroke={color} strokeWidth="1.8" strokeLinecap="round" />
    </svg>
  ),
  Laptop: ({ size = 26, color = '#2563EB' }: IconProps) => (
    <svg width={size} height={size} viewBox="0 0 24 24" fill="none">
      <rect x="5" y="5" width="14" height="9.5" rx="2" stroke={color} strokeWidth="1.8" />
      <path d="M3 18.5h18l-1-2.5H4l-1 2.5z" stroke={color} strokeWidth="1.8" strokeLinejoin="round" />
    </svg>
  ),
  Monitor: ({ size = 26, color = '#2563EB' }: IconProps) => (
    <svg width={size} height={size} viewBox="0 0 24 24" fill="none">
      <rect x="3.5" y="4.5" width="9" height="9" rx="1" stroke={color} strokeWidth="1.7" />
      <rect x="13" y="4.5" width="7.5" height="6" rx="1" stroke={color} strokeWidth="1.7" />
      <rect x="3.5" y="14.5" width="9" height="5" rx="1" stroke={color} strokeWidth="1.7" />
      <rect x="14.5" y="12" width="6" height="7.5" rx="1" stroke={color} strokeWidth="1.7" />
    </svg>
  ),
  Refresh: ({ size = 20, color = '#2563EB' }: IconProps) => (
    <svg width={size} height={size} viewBox="0 0 24 24" fill="none" aria-hidden="true">
      <path
        d="M20 8v-2a8 8 0 10-8 8"
        stroke={color}
        strokeWidth="1.85"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
      <path
        d="M4 16v2a8 8 0 008-8"
        stroke={color}
        strokeWidth="1.85"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
      <path
        d="M20 4h-4M20 4l-2.5 2.5M4 20h4M4 20l2.5-2.5"
        stroke={color}
        strokeWidth="1.85"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
    </svg>
  ),
};

export type IconName = keyof typeof Ic;
