/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_API_BASE_URL?: string;
  readonly VITE_SUB_BASE_URL?: string;
  readonly VITE_USE_MOCK?: string;
  readonly VITE_TELEGRAM_BOT_USERNAME?: string;
  readonly VITE_DEV_API_PROXY?: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}
