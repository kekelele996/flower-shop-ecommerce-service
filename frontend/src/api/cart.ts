import { del, get, post, put } from '../utils/request';
import type { ProductVO } from './product';

export interface CartItemVO {
  id: number;
  user_id: number;
  product_id: number;
  quantity: number;
  selected: boolean;
  product: ProductVO;
  subtotal: number;
}

export interface CartSummaryVO {
  items: CartItemVO[];
  total_quantity: number;
  total_price: number;
}

export const getCart = () => get<CartSummaryVO>('/cart');

export const addCart = (data: { product_id: number; quantity: number }) => post<CartItemVO>('/cart', data);

export const updateCart = (id: number, data: { quantity?: number; selected?: boolean }) => put<{ message: string }>(`/cart/${id}`, data);

export const removeCart = (id: number) => del<{ message: string }>(`/cart/${id}`);
