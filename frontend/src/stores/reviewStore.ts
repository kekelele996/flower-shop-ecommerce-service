import { create } from 'zustand';
import type { ReviewVO } from '../api/review';
import { listReviews } from '../api/review';

interface ReviewState {
  reviews: ReviewVO[];
  total: number;
  loading: boolean;
  fetchReviews: (params?: Record<string, unknown>) => Promise<void>;
}

export const useReviewStore = create<ReviewState>((set) => ({
  reviews: [],
  total: 0,
  loading: false,
  fetchReviews: async (params) => {
    set({ loading: true });
    try {
      const data = await listReviews(params);
      set({ reviews: data.list, total: data.total });
    } finally {
      set({ loading: false });
    }
  },
}));
