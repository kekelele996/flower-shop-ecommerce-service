import { get, post, put } from '../utils/request';

export interface UserVO {
  id: number;
  username: string;
  nickname: string;
  email: string;
  phone: string;
  avatar: string;
  role: string;
  status: string;
  created_at: string;
}

export interface TokenVO {
  access_token: string;
  token_type: string;
  expires_in: number;
  user: UserVO;
}

export const register = (data: { username: string; password: string; nickname?: string; email?: string }) =>
  post<{ message: string; user: UserVO }>('/auth/register', data);

export const login = (data: { username: string; password: string }) => post<TokenVO>('/auth/login', data);

export const getProfile = () => get<UserVO>('/auth/profile');

export const updateProfile = (data: { nickname?: string; email?: string; phone?: string; avatar?: string }) =>
  put<UserVO>('/auth/profile', data);

export const changePassword = (data: { old_password: string; new_password: string }) =>
  put<{ message: string }>('/auth/password', data);
