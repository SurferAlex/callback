import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { BrowserRouter } from "react-router-dom";
import { AuthProvider } from "@/contexts/AuthContext";
import { initTelegram, applyTelegramTheme, isTelegram } from "@/lib/telegram";
import App from "./App";
import "./index.css";
import "@/styles/surfer.css";

if (isTelegram()) {
  document.documentElement.classList.add("tma");
  document.documentElement.classList.remove("cabinet");
} else {
  document.documentElement.classList.add("cabinet");
  document.documentElement.classList.remove("tma");
}
initTelegram();
applyTelegramTheme();

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <BrowserRouter>
      <AuthProvider>
        <App />
      </AuthProvider>
    </BrowserRouter>
  </StrictMode>
);
