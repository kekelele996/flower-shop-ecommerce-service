import { get, post, put } from '../utils/request';
import type { ProductVO, PageData, ProductQuery } from './product';
import type { OrderVO, LogisticsVO } from './order';
import type { ReviewVO } from './review';

export interface AuditLogVO {
  id: number;
  user_id: number;
  username: string;
  action: string;
  method: string;
  path: string;
  entity: string;
  entity_id: string;
  detail: string;
  request_id: string;
  ip: string;
  created_at: string;
}

export interface StatsVO {
  total_sales: number;
  order_count: number;
  product_count: number;
  user_count: number;
  pending_shipment: number;
  order_status_count: Record<string, number>;
  top_products: { product_id: number; product_name: string; sales: number; revenue: number }[];
  daily_sales: { date: string; sales: number }[];
}

export const adminListProducts = (params?: ProductQuery) => get<PageData<ProductVO>>('/admin/products', params as Record<string, unknown>);

export const createProduct = (data: Partial<ProductVO>) => post<ProductVO>('/admin/products', data);

export const updateProduct = (id: number, data: Partial<ProductVO>) => put<ProductVO>(`/admin/products/${id}`, data);

export const changeProductStatus = (id: number, status: string) => put<ProductVO>(`/admin/products/${id}/status`, { status });

export const adminListOrders = (params?: { page?: number; page_size?: number; status?: string; order_no?: string }) =>
  get<PageData<OrderVO>>('/admin/orders', params as Record<string, unknown>);

export const shipOrder = (id: number, carrier?: string) => post<{ message: string; order: OrderVO }>(`/admin/orders/${id}/ship`, { carrier });

export const adminCompleteOrder = (id: number) => post<{ message: string; order: OrderVO }>(`/admin/orders/${id}/complete`);

export const adminGetLogistics = (id: number) => get<LogisticsVO>(`/admin/orders/${id}/logistics`);

export const adminListReviews = (params?: { page?: number; page_size?: number }) =>
  get<PageData<ReviewVO>>('/admin/reviews', params as Record<string, unknown>);

export const getStats = () => get<StatsVO>('/admin/stats');

export const listAuditLogs = (params?: { page?: number; page_size?: number; username?: string; action?: string; entity?: string }) =>
  get<PageData<AuditLogVO>>('/admin/audit-logs', params as Record<string, unknown>);
