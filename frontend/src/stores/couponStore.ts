import { create } from 'zustand';
import type { CouponVO } from '../api/coupon';
import { getAvailableCoupons, getMyCoupons } from '../api/coupon';

interface CouponState {
  coupons: CouponVO[];
  available: CouponVO[];
  fetchCoupons: () => Promise<void>;
  fetchAvailable: () => Promise<void>;
}

export const useCouponStore = create<CouponState>((set) => ({
  coupons: [],
  available: [],
  fetchCoupons: async () => {
    const list = await getMyCoupons();
    set({ coupons: list });
  },
  fetchAvailable: async () => {
    const list = await getAvailableCoupons();
    set({ available: list });
  },
}));
