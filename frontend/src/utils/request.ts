import axios, { AxiosError } from 'axios';
import { message } from 'antd';

export interface ApiResponse<T = unknown> {
  code: number;
  message: string;
  data: T;
}

// 统一 axios 实例：携带 JWT、注入 request_id、统一错误拦截（全局错误处理）。
const request = axios.create({
  baseURL: '/api/v1',
  timeout: 15000,
});

request.interceptors.request.use((config) => {
  const token = localStorage.getItem('flowershop_token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  config.headers['X-Request-ID'] = `${Date.now()}-${Math.random().toString(36).slice(2, 10)}`;
  return config;
});

request.interceptors.response.use(
  (response) => {
    const body = response.data as ApiResponse;
    if (body && body.code !== undefined && body.code !== 0) {
      message.error(body.message || '请求失败');
      return Promise.reject(new Error(body.message || '请求失败'));
    }
    return response;
  },
  (error: AxiosError<ApiResponse>) => {
    const status = error.response?.status;
    const msg = error.response?.data?.message || error.message || '网络错误';
    if (status === 401) {
      localStorage.removeItem('flowershop_token');
      localStorage.removeItem('flowershop_user');
      if (!window.location.pathname.startsWith('/login')) {
        message.warning('登录已过期，请重新登录');
        window.location.href = '/login';
      }
    } else {
      message.error(msg);
    }
    return Promise.reject(error);
  },
);

// 通用请求方法
export function get<T>(url: string, params?: Record<string, unknown>): Promise<T> {
  return request.get(url, { params }).then((res) => res.data.data as T);
}

export function post<T>(url: string, data?: unknown): Promise<T> {
  return request.post(url, data).then((res) => res.data.data as T);
}

export function put<T>(url: string, data?: unknown): Promise<T> {
  return request.put(url, data).then((res) => res.data.data as T);
}

export function del<T>(url: string): Promise<T> {
  return request.delete(url).then((res) => res.data.data as T);
}

export default request;
