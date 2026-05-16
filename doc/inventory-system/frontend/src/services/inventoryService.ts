import api from '@/lib/api';
import type {
  ApiResponse,
  PageResult,
  Inventory,
  CreateInventoryRequest,
  UpdateInventoryRequest,
  InventoryQuery,
} from '@/types';

// 获取库存列表
export async function getInventoryList(query: InventoryQuery): Promise<PageResult<Inventory>> {
  const res = await api.get<ApiResponse<PageResult<Inventory>>>('/inventories', { params: query });
  return res.data.data;
}

// 获取库存详情
export async function getInventoryById(id: number): Promise<Inventory> {
  const res = await api.get<ApiResponse<Inventory>>(`/inventories/${id}`);
  return res.data.data;
}

// 新增库存
export async function createInventory(data: CreateInventoryRequest): Promise<Inventory> {
  const res = await api.post<ApiResponse<Inventory>>('/inventories', data);
  return res.data.data;
}

// 编辑库存
export async function updateInventory(id: number, data: UpdateInventoryRequest): Promise<void> {
  await api.put(`/inventories/${id}`, data);
}

// 删除库存
export async function deleteInventory(id: number): Promise<void> {
  await api.delete(`/inventories/${id}`);
}

// 批量删除
export async function batchDeleteInventory(ids: number[]): Promise<void> {
  await api.post('/inventories/batch-delete', { ids });
}

// 恢复已删除
export async function restoreInventory(id: number): Promise<void> {
  await api.put(`/inventories/${id}/restore`);
}

// 导出Excel
export async function exportInventory(query: InventoryQuery): Promise<Blob> {
  const res = await api.get('/inventories/export', {
    params: query,
    responseType: 'blob',
  });
  return res.data;
}
