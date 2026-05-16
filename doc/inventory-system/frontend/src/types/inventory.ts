// 库存相关类型定义

export interface Inventory {
  id: number;
  material_code: string;
  material_name: string;
  category: string;
  specification: string;
  unit: string;
  quantity: number;
  safety_stock: number;
  location: string;
  status: number;
  remark: string;
  version: number;
  is_deleted: number;
  created_by: number;
  updated_by: number | null;
  created_at: string;
  updated_at: string;
  creator_name?: string;
  warning: boolean;
  attachments?: Attachment[];
}

export interface Attachment {
  id: number;
  inventory_id: number;
  file_name: string;
  display_name: string;
  file_path: string;
  file_size: number;
  file_type: 'image' | 'document';
  mime_type: string;
  extension: string;
  thumbnail_path: string;
  uploaded_by: number;
  uploader_name?: string;
  created_at: string;
  updated_at: string;
}

export interface CreateInventoryRequest {
  material_code: string;
  material_name: string;
  category: string;
  specification?: string;
  unit: string;
  quantity: number;
  safety_stock?: number;
  location?: string;
  status?: number;
  remark?: string;
}

export interface UpdateInventoryRequest {
  material_name?: string;
  category?: string;
  specification?: string;
  unit?: string;
  quantity?: number;
  safety_stock?: number;
  location?: string;
  status?: number;
  remark?: string;
  version: number;
}

export interface InventoryQuery {
  page?: number;
  page_size?: number;
  keyword?: string;
  category?: string;
  status?: number;
  min_quantity?: number;
  max_quantity?: number;
  location?: string;
  show_deleted?: boolean;
}

// 分类选项
export const CATEGORY_OPTIONS = [
  { value: 'raw_material', label: '原料' },
  { value: 'finished_product', label: '成品' },
  { value: 'spare_part', label: '备件' },
  { value: 'other', label: '其他' },
];

// 分类映射
export const CATEGORY_MAP: Record<string, string> = {
  raw_material: '原料',
  finished_product: '成品',
  spare_part: '备件',
  other: '其他',
};
