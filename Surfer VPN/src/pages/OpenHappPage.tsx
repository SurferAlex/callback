import { useEffect } from "react";
import { useSearchParams } from "react-router-dom";

/** Bridge page: Telegram openLink → Safari/WebView → happ:// deep link. */
export function OpenHappPage() {
  const [params] = useSearchParams();

  useEffect(() => {
    const raw = params.get("url")?.trim() ?? "";
    if (!raw.startsWith("happ://")) return;
    window.location.replace(raw);
  }, [params]);

  return (
    <div className="screen">
      <main className="page" style={{ padding: "48px 22px", textAlign: "center" }}>
        <p style={{ margin: 0, color: "var(--ink-2)" }}>Открываем Happ…</p>
      </main>
    </div>
  );
}
