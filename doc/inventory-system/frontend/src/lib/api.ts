import axios from 'axios';
import type { ApiResponse } from '@/types';

// 创建 Axios 实例
const api = axios.create({
  baseURL: '/api/v1',
  timeout: 15000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// 请求拦截器：自动注入 Token
api.interceptors.request.use(
  (config) => {
    if (typeof window !== 'undefined') {
      const token = localStorage.getItem('token');
      if (token) {
        config.headers.Authorization = `Bearer ${token}`;
      }
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// 响应拦截器：统一处理错误
api.interceptors.response.use(
  (response) => {
    const data = response.data as ApiResponse;
    // 业务错误（code != 0）
    if (data.code !== 0) {
      // Token 过期或未登录
      if (data.code === 40101) {
        if (typeof window !== 'undefined') {
          localStorage.removeItem('token');
          localStorage.removeItem('user');
          window.location.href = '/login';
        }
      }
      return Promise.reject(new Error(data.message || '请求失败'));
    }
    return response;
  },
  (error) => {
    if (error.response) {
      const msg = error.response.data?.message || '服务器错误';
      return Promise.reject(new Error(msg));
    }
    if (error.code === 'ECONNABORTED') {
      return Promise.reject(new Error('请求超时，请稍后重试'));
    }
    return Promise.reject(new Error('网络错误，请检查网络连接'));
  }
);

export default api;
