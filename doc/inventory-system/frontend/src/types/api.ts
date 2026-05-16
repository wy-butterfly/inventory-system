// API 通用类型定义

// 统一API响应格式
export interface ApiResponse<T = any> {
  code: number;
  message: string;
  data: T;
}

// 分页结果
export interface PageResult<T = any> {
  list: T[];
  total: number;
  page: number;
  page_size: number;
}

// 分页查询参数
export interface PageQuery {
  page?: number;
  page_size?: number;
}
