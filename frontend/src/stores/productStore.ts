import { create } from 'zustand';
import type { ProductVO } from '../api/product';
import { getRecommendations, listProducts } from '../api/product';

interface ProductState {
  products: ProductVO[];
  total: number;
  loading: boolean;
  recommendations: ProductVO[];
  fetchProducts: (params?: Record<string, unknown>) => Promise<void>;
  fetchRecommendations: () => Promise<void>;
}

export const useProductStore = create<ProductState>((set) => ({
  products: [],
  total: 0,
  loading: false,
  recommendations: [],
  fetchProducts: async (params) => {
    set({ loading: true });
    try {
      const data = await listProducts(params);
      set({ products: data.list, total: data.total });
    } finally {
      set({ loading: false });
    }
  },
  fetchRecommendations: async () => {
    try {
      const list = await getRecommendations();
      set({ recommendations: list });
    } catch {
      set({ recommendations: [] });
    }
  },
}));
