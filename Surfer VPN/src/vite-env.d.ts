/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_API_BASE_URL?: string;
  readonly VITE_SUB_BASE_URL?: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}
