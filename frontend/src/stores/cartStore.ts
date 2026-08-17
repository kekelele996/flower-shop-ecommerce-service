import { create } from 'zustand';
import type { CartSummaryVO } from '../api/cart';
import { getCart } from '../api/cart';

interface CartState {
  summary: CartSummaryVO | null;
  loading: boolean;
  fetchCart: () => Promise<void>;
  clear: () => void;
}

export const useCartStore = create<CartState>((set) => ({
  summary: null,
  loading: false,
  fetchCart: async () => {
    set({ loading: true });
    try {
      const summary = await getCart();
      set({ summary });
    } finally {
      set({ loading: false });
    }
  },
  clear: () => set({ summary: null }),
}));
