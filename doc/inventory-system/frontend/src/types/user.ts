// 用户相关类型定义

export interface User {
  id: number;
  username: string;
  display_name: string;
  role: 'admin' | 'operator' | 'viewer';
  status: number;
  created_at: string;
  updated_at: string;
}

export interface LoginRequest {
  username: string;
  password: string;
}

export interface LoginResponse {
  token: string;
  user: Pick<User, 'id' | 'username' | 'display_name' | 'role'>;
}

export interface CreateUserRequest {
  username: string;
  password: string;
  display_name?: string;
  role: 'admin' | 'operator' | 'viewer';
}

export interface UpdateUserRequest {
  display_name?: string;
  role?: 'admin' | 'operator' | 'viewer';
  status?: number;
}

export interface UserQuery {
  page?: number;
  page_size?: number;
  keyword?: string;
  role?: string;
  status?: number;
}
