import { create } from 'zustand';
import type { OrderVO } from '../api/order';
import { listOrders } from '../api/order';

interface OrderState {
  orders: OrderVO[];
  total: number;
  loading: boolean;
  fetchOrders: (params?: { page?: number; page_size?: number; status?: string }) => Promise<void>;
}

export const useOrderStore = create<OrderState>((set) => ({
  orders: [],
  total: 0,
  loading: false,
  fetchOrders: async (params) => {
    set({ loading: true });
    try {
      const data = await listOrders(params);
      set({ orders: data.list, total: data.total });
    } finally {
      set({ loading: false });
    }
  },
}));
