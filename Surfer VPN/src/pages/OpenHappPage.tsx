import { useEffect } from "react";
import { useSearchParams } from "react-router-dom";

/** Bridge for Telegram Mini App: external browser opens, then redirects to happ:// */
export function OpenHappPage() {
  const [params] = useSearchParams();
  const url = params.get("url") ?? "";

  useEffect(() => {
    if (url.startsWith("happ://")) {
      window.location.replace(url);
    }
  }, [url]);

  return (
    <main className="page" style={{ padding: "24px", textAlign: "center" }}>
      <p>Открываем Happ…</p>
    </main>
  );
}
