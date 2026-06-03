import { useCallback, useEffect, useRef, useState } from "react";

export interface UseClipboardResult {
  copied: boolean;
  copy: (text: string) => Promise<boolean>;
}

/**
 * Copy text to the clipboard with a graceful fallback for environments where
 * the async Clipboard API is unavailable (older WebViews / insecure contexts).
 *
 * `copied` flips to `true` on success and resets to `false` after `timeout` ms.
 */
export function useClipboard(timeout = 2000): UseClipboardResult {
  const [copied, setCopied] = useState<boolean>(false);
  const timerRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  // Clear any pending reset on unmount.
  useEffect(() => {
    return () => {
      if (timerRef.current !== null) {
        clearTimeout(timerRef.current);
      }
    };
  }, []);

  const copy = useCallback(
    async (text: string): Promise<boolean> => {
      let ok = false;

      try {
        if (
          typeof navigator !== "undefined" &&
          navigator.clipboard?.writeText
        ) {
          await navigator.clipboard.writeText(text);
          ok = true;
        } else {
          ok = fallbackCopy(text);
        }
      } catch {
        // Async API can reject (permissions / focus) — try the legacy path.
        ok = fallbackCopy(text);
      }

      if (ok) {
        setCopied(true);
        if (timerRef.current !== null) {
          clearTimeout(timerRef.current);
        }
        timerRef.current = setTimeout(() => {
          setCopied(false);
          timerRef.current = null;
        }, timeout);
      }

      return ok;
    },
    [timeout]
  );

  return { copied, copy };
}

/** Legacy clipboard write using a hidden textarea + `execCommand`. */
function fallbackCopy(text: string): boolean {
  if (typeof document === "undefined") return false;

  const textarea = document.createElement("textarea");
  textarea.value = text;
  textarea.setAttribute("readonly", "");
  textarea.style.position = "fixed";
  textarea.style.top = "-9999px";
  textarea.style.left = "-9999px";
  textarea.style.opacity = "0";
  document.body.appendChild(textarea);

  let ok = false;
  try {
    textarea.focus();
    textarea.select();
    ok = document.execCommand("copy");
  } catch {
    ok = false;
  } finally {
    document.body.removeChild(textarea);
  }

  return ok;
}
