import { Navigate, useLocation } from "react-router-dom";

import { useAuth } from "@/contexts/AuthContext";
import { isWebCabinet } from "@/lib/runtime";

export function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const { ready, authenticated } = useAuth();
  const location = useLocation();

  if (!ready) {
    return (
      <div className="screen auth-loading">
        <div className="auth-loading-inner">🏄 Surf VPN</div>
      </div>
    );
  }

  if (isWebCabinet() && !authenticated) {
    return <Navigate to="/login" replace state={{ from: location.pathname }} />;
  }

  return <>{children}</>;
}
