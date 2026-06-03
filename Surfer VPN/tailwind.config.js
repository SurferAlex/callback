/** @type {import('tailwindcss').Config} */
export default {
  darkMode: ["class"],
  content: ["./index.html", "./src/**/*.{ts,tsx}"],
  theme: {
    container: {
      center: true,
      padding: "1rem",
      screens: {
        "2xl": "480px",
      },
    },
    extend: {
      colors: {
        border: "hsl(var(--border))",
        input: "hsl(var(--input))",
        ring: "hsl(var(--ring))",
        background: "hsl(var(--background))",
        foreground: "hsl(var(--foreground))",
        primary: {
          DEFAULT: "hsl(var(--primary))",
          foreground: "hsl(var(--primary-foreground))",
        },
        secondary: {
          DEFAULT: "hsl(var(--secondary))",
          foreground: "hsl(var(--secondary-foreground))",
        },
        destructive: {
          DEFAULT: "hsl(var(--destructive))",
          foreground: "hsl(var(--destructive-foreground))",
        },
        muted: {
          DEFAULT: "hsl(var(--muted))",
          foreground: "hsl(var(--muted-foreground))",
        },
        accent: {
          DEFAULT: "hsl(var(--accent))",
          foreground: "hsl(var(--accent-foreground))",
        },
        popover: {
          DEFAULT: "hsl(var(--popover))",
          foreground: "hsl(var(--popover-foreground))",
        },
        card: {
          DEFAULT: "hsl(var(--card))",
          foreground: "hsl(var(--card-foreground))",
        },
        // Ocean brand palette
        ocean: {
          50: "#ecfeff",
          100: "#cffafe",
          200: "#a5f3fc",
          300: "#67e8f9",
          400: "#22d3ee",
          500: "#06b6d4",
          600: "#0891b2",
          700: "#0e7490",
          800: "#155e75",
          900: "#164e63",
        },
        sky: {
          50: "#f0f9ff",
          100: "#e0f2fe",
          200: "#bae6fd",
          300: "#7dd3fc",
          400: "#38bdf8",
          500: "#0ea5e9",
          600: "#0284c7",
          700: "#0369a1",
        },
        turquoise: {
          DEFAULT: "#2dd4bf",
          light: "#5eead4",
          dark: "#14b8a6",
        },
        foam: "#f8fdff",
      },
      borderRadius: {
        lg: "var(--radius)",
        md: "calc(var(--radius) - 4px)",
        sm: "calc(var(--radius) - 8px)",
        "2xl": "calc(var(--radius) + 8px)",
        "3xl": "calc(var(--radius) + 16px)",
      },
      fontFamily: {
        sans: [
          "Inter",
          "-apple-system",
          "BlinkMacSystemFont",
          "Segoe UI",
          "Roboto",
          "Helvetica Neue",
          "sans-serif",
        ],
      },
      backgroundImage: {
        "gradient-ocean":
          "linear-gradient(135deg, #38bdf8 0%, #22d3ee 45%, #2dd4bf 100%)",
        "gradient-ocean-soft":
          "linear-gradient(135deg, #e0f2fe 0%, #cffafe 50%, #ccfbf1 100%)",
        "gradient-sky":
          "linear-gradient(180deg, #e0f2fe 0%, #f0f9ff 60%, #ffffff 100%)",
        "gradient-cta":
          "linear-gradient(135deg, #0a6cff 0%, #1e90ff 52%, #3bb0ff 100%)",
        "gradient-cta-deep":
          "linear-gradient(135deg, #075fe4 0%, #1685ff 100%)",
        "gradient-aurora":
          "radial-gradient(60% 80% at 20% 10%, rgba(56,189,248,0.35) 0%, rgba(56,189,248,0) 60%), radial-gradient(70% 90% at 90% 20%, rgba(45,212,191,0.30) 0%, rgba(45,212,191,0) 55%)",
      },
      boxShadow: {
        glass: "0 8px 32px rgba(14, 165, 233, 0.12)",
        "glass-lg": "0 16px 48px rgba(6, 182, 212, 0.18)",
        soft: "0 4px 20px rgba(14, 116, 144, 0.08)",
        card: "0 12px 34px rgba(30, 144, 255, 0.12)",
        "card-hover": "0 18px 44px rgba(30, 144, 255, 0.18)",
        cta: "0 12px 30px rgba(30, 144, 255, 0.42)",
        "cta-hover": "0 18px 44px rgba(30, 144, 255, 0.52)",
        glow: "0 0 0 1px rgba(255,255,255,0.6) inset, 0 8px 28px rgba(45,212,191,0.25)",
      },
      keyframes: {
        float: {
          "0%, 100%": { transform: "translateY(0px)" },
          "50%": { transform: "translateY(-12px)" },
        },
        "float-slow": {
          "0%, 100%": { transform: "translateY(0px) rotate(0deg)" },
          "50%": { transform: "translateY(-8px) rotate(2deg)" },
        },
        wave: {
          "0%": { transform: "translateX(0)" },
          "100%": { transform: "translateX(-50%)" },
        },
        shimmer: {
          "0%": { backgroundPosition: "-200% 0" },
          "100%": { backgroundPosition: "200% 0" },
        },
        "fade-in-up": {
          "0%": { opacity: "0", transform: "translateY(16px)" },
          "100%": { opacity: "1", transform: "translateY(0)" },
        },
        "fade-in": {
          "0%": { opacity: "0" },
          "100%": { opacity: "1" },
        },
        "scale-in": {
          "0%": { opacity: "0", transform: "scale(0.96)" },
          "100%": { opacity: "1", transform: "scale(1)" },
        },
        "pulse-ring": {
          "0%": { transform: "scale(0.95)", opacity: "0.7" },
          "70%": { transform: "scale(1.25)", opacity: "0" },
          "100%": { transform: "scale(1.25)", opacity: "0" },
        },
      },
      animation: {
        float: "float 6s ease-in-out infinite",
        "float-slow": "float-slow 8s ease-in-out infinite",
        wave: "wave 12s linear infinite",
        "wave-slow": "wave 20s linear infinite",
        shimmer: "shimmer 2.5s linear infinite",
        "fade-in-up": "fade-in-up 0.6s ease-out both",
        "fade-in": "fade-in 0.5s ease-out both",
        "scale-in": "scale-in 0.4s ease-out both",
        "pulse-ring": "pulse-ring 2.4s cubic-bezier(0.4,0,0.6,1) infinite",
      },
    },
  },
  plugins: [require("tailwindcss-animate")],
};
