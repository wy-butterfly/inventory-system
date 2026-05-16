'use client';

import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { message } from 'antd';
import * as attachmentService from '@/services/attachmentService';

// 附件列表
export function useAttachmentList(inventoryId: number) {
  return useQuery({
    queryKey: ['attachments', inventoryId],
    queryFn: () => attachmentService.getAttachmentList(inventoryId),
    enabled: !!inventoryId,
  });
}

// 上传附件
export function useUploadAttachment() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ inventoryId, file }: { inventoryId: number; file: File }) =>
      attachmentService.uploadAttachment(inventoryId, file),
    onSuccess: () => {
      message.success('上传成功');
      queryClient.invalidateQueries({ queryKey: ['attachments'] });
      queryClient.invalidateQueries({ queryKey: ['inventory'] });
    },
    onError: (err: Error) => {
      message.error(err.message || '上传失败');
    },
  });
}

// 删除附件
export function useDeleteAttachment() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: number) => attachmentService.deleteAttachment(id),
    onSuccess: () => {
      message.success('删除成功');
      queryClient.invalidateQueries({ queryKey: ['attachments'] });
      queryClient.invalidateQueries({ queryKey: ['inventory'] });
    },
    onError: (err: Error) => {
      message.error(err.message || '删除失败');
    },
  });
}

// 重命名附件
export function useRenameAttachment() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, displayName }: { id: number; displayName: string }) =>
      attachmentService.renameAttachment(id, displayName),
    onSuccess: () => {
      message.success('重命名成功');
      queryClient.invalidateQueries({ queryKey: ['attachments'] });
    },
    onError: (err: Error) => {
      message.error(err.message || '重命名失败');
    },
  });
}
