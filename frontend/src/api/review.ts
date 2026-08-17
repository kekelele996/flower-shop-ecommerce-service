import { get, post } from '../utils/request';
import type { PageData } from './product';

export interface ReviewVO {
  id: number;
  user_id: number;
  username: string;
  order_id: number;
  order_item_id: number;
  product_id: number;
  product_name: string;
  rating: number;
  content: string;
  images: string[];
  reply: string;
  replied_at: string;
  created_at: string;
}

export const createReview = (data: { order_item_id: number; rating: number; content: string; images?: string[] }) =>
  post<{ message: string; review: ReviewVO }>('/reviews', data);

export const listReviews = (params?: { page?: number; page_size?: number; product_id?: number; order_id?: number; user_id?: number }) =>
  get<PageData<ReviewVO>>('/reviews', params as Record<string, unknown>);

export const replyReview = (id: number, reply: string) => post<{ message: string; review: ReviewVO }>(`/admin/reviews/${id}/reply`, { reply });
