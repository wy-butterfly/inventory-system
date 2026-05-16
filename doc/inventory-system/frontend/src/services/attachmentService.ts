import api from '@/lib/api';
import type { ApiResponse, Attachment } from '@/types';

// 获取库存的附件列表
export async function getAttachmentList(inventoryId: number): Promise<Attachment[]> {
  const res = await api.get<ApiResponse<Attachment[]>>(`/inventories/${inventoryId}/attachments`);
  return res.data.data;
}

// 上传附件
export async function uploadAttachment(inventoryId: number, file: File): Promise<Attachment> {
  const formData = new FormData();
  formData.append('file', file);
  const res = await api.post<ApiResponse<Attachment>>(
    `/inventories/${inventoryId}/attachments`,
    formData,
    { headers: { 'Content-Type': 'multipart/form-data' } }
  );
  return res.data.data;
}

// 重命名附件
export async function renameAttachment(id: number, displayName: string): Promise<void> {
  await api.put(`/attachments/${id}`, { display_name: displayName });
}

// 删除附件
export async function deleteAttachment(id: number): Promise<void> {
  await api.delete(`/attachments/${id}`);
}

// 获取附件预览URL
export function getPreviewUrl(id: number, thumbnail = false): string {
  const base = `/api/v1/attachments/${id}/preview`;
  return thumbnail ? `${base}?thumbnail=true` : base;
}

// 获取附件下载URL
export function getDownloadUrl(id: number): string {
  return `/api/v1/attachments/${id}/download`;
}
