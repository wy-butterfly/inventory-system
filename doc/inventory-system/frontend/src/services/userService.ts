import api from '@/lib/api';
import type { ApiResponse, PageResult, User, CreateUserRequest, UpdateUserRequest, UserQuery } from '@/types';

// 获取用户列表
export async function getUserList(query: UserQuery): Promise<PageResult<User>> {
  const res = await api.get<ApiResponse<PageResult<User>>>('/users', { params: query });
  return res.data.data;
}

// 创建用户
export async function createUser(data: CreateUserRequest): Promise<User> {
  const res = await api.post<ApiResponse<User>>('/users', data);
  return res.data.data;
}

// 编辑用户
export async function updateUser(id: number, data: UpdateUserRequest): Promise<void> {
  await api.put(`/users/${id}`, data);
}

// 重置密码
export async function resetPassword(id: number, newPassword: string): Promise<void> {
  await api.put(`/users/${id}/reset-password`, { new_password: newPassword });
}
