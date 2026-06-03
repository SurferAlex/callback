import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { BrowserRouter } from "react-router-dom";
import { initTelegram, applyTelegramTheme } from "@/lib/telegram";
import App from "./App";
import "./index.css";
import "@/styles/surfer.css";

// Boot the Telegram Mini App runtime (no-op in a plain browser).
initTelegram();
applyTelegramTheme();

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <BrowserRouter>
      <App />
    </BrowserRouter>
  </StrictMode>
);
