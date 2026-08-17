import { get, post, put } from '../utils/request';

export interface ProductVO {
  id: number;
  category_id: number;
  category_name: string;
  name: string;
  sub_title: string;
  description: string;
  detail: string;
  cover_image: string;
  images: string[];
  price: number;
  original_price: number;
  stock: number;
  sales: number;
  rating: number;
  rating_count: number;
  shipping_from: string;
  free_shipping: boolean;
  status: string;
  created_at: string;
}

export interface PageData<T> {
  list: T[];
  total: number;
  page: number;
  page_size: number;
}

export interface ProductQuery {
  page?: number;
  page_size?: number;
  keyword?: string;
  category_id?: number;
  min_price?: number;
  max_price?: number;
  free_shipping?: boolean;
  shipping_from?: string;
  sort?: string;
  status?: string;
}

export const listProducts = (params?: ProductQuery) => get<PageData<ProductVO>>('/products', params as Record<string, unknown>);

export const getProduct = (id: number) => get<ProductVO>(`/products/${id}`);

export const recordView = (id: number) => post<{ message: string }>(`/products/${id}/view`);

export const getRecommendations = () => get<ProductVO[]>('/products/recommendations');
