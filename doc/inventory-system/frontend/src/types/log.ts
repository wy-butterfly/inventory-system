// 日志相关类型定义

export interface OperationLog {
  id: number;
  user_id: number;
  user_name: string;
  action: string;
  target_type: string;
  target_id: number | null;
  description: string;
  ip_address: string;
  created_at: string;
}

export interface ChangeLog {
  id: number;
  inventory_id: number;
  user_id: number;
  user_name: string;
  field_name: string;
  field_label: string;
  old_value: string;
  new_value: string;
  created_at: string;
}

// 操作类型映射
export const ACTION_MAP: Record<string, string> = {
  create: '新增',
  update: '编辑',
  delete: '删除',
  restore: '恢复',
  upload: '上传',
  download: '下载',
  export: '导出',
  login: '登录',
  logout: '登出',
};
