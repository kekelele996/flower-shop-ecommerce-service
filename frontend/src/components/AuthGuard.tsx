import { Navigate, useLocation } from 'react-router-dom';
import type { ReactNode } from 'react';
import { useUserStore } from '../stores/userStore';

interface AuthGuardProps {
  children: ReactNode;
  adminOnly?: boolean;
}

// AuthGuard：路由守卫（未登录跳登录页，非管理员跳首页）。
export default function AuthGuard({ children, adminOnly = false }: AuthGuardProps) {
  const location = useLocation();
  const { isLogin, isAdmin } = useUserStore();
  if (!isLogin()) {
    return <Navigate to="/login" state={{ from: location.pathname }} replace />;
  }
  if (adminOnly && !isAdmin()) {
    return <Navigate to="/" replace />;
  }
  return <>{children}</>;
}
