import { create } from 'zustand';
import type { UserVO } from '../api/auth';
import { getProfile } from '../api/auth';

interface UserState {
  token: string | null;
  user: UserVO | null;
  setToken: (token: string) => void;
  setUser: (user: UserVO | null) => void;
  logout: () => void;
  refreshProfile: () => Promise<void>;
  isLogin: () => boolean;
  isAdmin: () => boolean;
}

export const useUserStore = create<UserState>((set, get) => ({
  token: localStorage.getItem('flowershop_token'),
  user: (() => {
    try {
      return JSON.parse(localStorage.getItem('flowershop_user') || 'null');
    } catch {
      return null;
    }
  })(),
  setToken: (token) => {
    localStorage.setItem('flowershop_token', token);
    set({ token });
  },
  setUser: (user) => {
    if (user) {
      localStorage.setItem('flowershop_user', JSON.stringify(user));
    } else {
      localStorage.removeItem('flowershop_user');
    }
    set({ user });
  },
  logout: () => {
    localStorage.removeItem('flowershop_token');
    localStorage.removeItem('flowershop_user');
    set({ token: null, user: null });
  },
  refreshProfile: async () => {
    if (!get().token) return;
    try {
      const profile = await getProfile();
      get().setUser(profile);
    } catch {
      // 401 由拦截器统一处理
    }
  },
  isLogin: () => !!get().token,
  isAdmin: () => get().user?.role === 'ADMIN',
}));
