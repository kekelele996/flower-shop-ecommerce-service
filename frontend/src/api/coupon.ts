import { get, post } from '../utils/request';

export interface CouponVO {
  id: number;
  name: string;
  type: string;
  type_text: string;
  threshold: number;
  amount: number;
  discount_rate: number;
  status: string;
  valid_from: string;
  valid_to: string;
}

export const getMyCoupons = () => get<CouponVO[]>('/coupons');

export const getAvailableCoupons = () => get<CouponVO[]>('/coupons/available');

export const claimCoupon = (id: number) => post<{ message: string; coupon: CouponVO }>(`/coupons/${id}/claim`);

export const createCouponTemplate = (data: {
  name: string;
  type: string;
  threshold: number;
  amount: number;
  discount_rate: number;
  total: number;
  valid_days: number;
}) => post<{ message: string; created: number }>('/admin/coupons', data);
