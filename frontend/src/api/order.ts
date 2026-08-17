import { get, post } from '../utils/request';
import type { PageData } from './product';

export interface OrderItemVO {
  id: number;
  product_id: number;
  product_name: string;
  product_image: string;
  price: number;
  quantity: number;
  total_price: number;
  reviewed: boolean;
}

export interface OrderVO {
  id: number;
  order_no: string;
  user_id: number;
  status: string;
  status_text: string;
  total_amount: number;
  discount_amount: number;
  shipping_fee: number;
  pay_amount: number;
  coupon_id: number;
  receiver_name: string;
  receiver_phone: string;
  receiver_addr: string;
  remark: string;
  paid_at: string;
  shipped_at: string;
  completed_at: string;
  cancelled_at: string;
  created_at: string;
  items: OrderItemVO[];
}

export interface LogisticsEventVO {
  time: string;
  status: string;
  desc: string;
}

export interface LogisticsVO {
  tracking_no: string;
  carrier: string;
  status: string;
  status_text: string;
  events: LogisticsEventVO[];
}

export interface CheckoutRequest {
  cart_item_ids: number[];
  coupon_id?: number;
  receiver_name: string;
  receiver_phone: string;
  receiver_addr: string;
  remark?: string;
}

export const checkout = (data: CheckoutRequest) => post<{ message: string; order: OrderVO }>('/orders', data);

export const listOrders = (params?: { page?: number; page_size?: number; status?: string; order_no?: string }) =>
  get<PageData<OrderVO>>('/orders', params as Record<string, unknown>);

export const getOrder = (id: number) => get<OrderVO>(`/orders/${id}`);

export const payOrder = (id: number, channel = 'ALIPAY_SANDBOX') => post<{ message: string; order: OrderVO }>(`/orders/${id}/pay`, { channel });

export const cancelOrder = (id: number, reason?: string) => post<{ message: string; order: OrderVO }>(`/orders/${id}/cancel`, { reason });

export const completeOrder = (id: number) => post<{ message: string; order: OrderVO }>(`/orders/${id}/complete`);

export const getLogistics = (id: number) => get<LogisticsVO>(`/orders/${id}/logistics`);
