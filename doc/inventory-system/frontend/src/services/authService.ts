import api from '@/lib/api';
import type { ApiResponse, LoginRequest, LoginResponse, User } from '@/types';

// 登录
export async function login(data: LoginRequest): Promise<LoginResponse> {
  const res = await api.post<ApiResponse<LoginResponse>>('/auth/login', data);
  return res.data.data;
}

// 获取当前用户信息
export async function getMe(): Promise<User> {
  const res = await api.get<ApiResponse<User>>('/auth/me');
  return res.data.data;
}

// 登出
export async function logout(): Promise<void> {
  await api.post('/auth/logout');
}
