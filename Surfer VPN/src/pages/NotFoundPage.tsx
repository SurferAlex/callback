import { Link } from "react-router-dom";

import { FullLogo } from "@/components/surf/logos";

/** Minimal on-brand 404. */
export function NotFoundPage() {
  return (
    <div className="screen">
      <div
        style={{
          minHeight: "100dvh",
          display: "flex",
          flexDirection: "column",
          alignItems: "center",
          justifyContent: "center",
          gap: 24,
          padding: "24px",
          textAlign: "center",
        }}
      >
        <FullLogo className="splash-logo" />
        <h1
          style={{
            margin: 0,
            fontSize: 22,
            fontWeight: 800,
            color: "#0f3a5f",
            letterSpacing: "-0.01em",
          }}
        >
          Страница не найдена
        </h1>
        <Link to="/" className="btn-ghost" style={{ textDecoration: "none" }}>
          <span>На главную</span>
        </Link>
      </div>
    </div>
  );
}
