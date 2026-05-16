'use client';

import React, { createContext, useContext, useState, useEffect, useCallback } from 'react';
import { useRouter } from 'next/navigation';
import { message } from 'antd';
import * as authService from '@/services/authService';
import type { User, LoginRequest } from '@/types';

interface AuthContextType {
  user: Pick<User, 'id' | 'username' | 'display_name' | 'role'> | null;
  loading: boolean;
  login: (data: LoginRequest) => Promise<void>;
  logout: () => void;
  isAdmin: boolean;
  isOperator: boolean;
  hasWritePermission: boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<AuthContextType['user']>(null);
  const [loading, setLoading] = useState(true);
  const router = useRouter();

  // 初始化：从 localStorage 恢复登录状态
  useEffect(() => {
    const storedUser = localStorage.getItem('user');
    const token = localStorage.getItem('token');
    if (storedUser && token) {
      try {
        setUser(JSON.parse(storedUser));
      } catch {
        localStorage.removeItem('user');
        localStorage.removeItem('token');
      }
    }
    setLoading(false);
  }, []);

  // 登录
  const login = useCallback(async (data: LoginRequest) => {
    const res = await authService.login(data);
    localStorage.setItem('token', res.token);
    localStorage.setItem('user', JSON.stringify(res.user));
    setUser(res.user);
    message.success('登录成功');

    router.push('/admin/inventory');
  }, [router]);

  // 登出
  const logout = useCallback(() => {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    setUser(null);
    router.push('/login');
  }, [router]);

  const isAdmin = user?.role === 'admin';
  const isOperator = user?.role === 'operator';
  const hasWritePermission = isAdmin || isOperator;

  return (
    <AuthContext.Provider value={{ user, loading, login, logout, isAdmin, isOperator, hasWritePermission }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}
