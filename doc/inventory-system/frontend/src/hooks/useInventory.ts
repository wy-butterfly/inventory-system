'use client';

import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { message } from 'antd';
import * as inventoryService from '@/services/inventoryService';
import type { InventoryQuery, CreateInventoryRequest, UpdateInventoryRequest } from '@/types';

// 库存列表
export function useInventoryList(query: InventoryQuery) {
  return useQuery({
    queryKey: ['inventories', query],
    queryFn: () => inventoryService.getInventoryList(query),
  });
}

// 库存详情
export function useInventoryDetail(id: number) {
  return useQuery({
    queryKey: ['inventory', id],
    queryFn: () => inventoryService.getInventoryById(id),
    enabled: !!id,
  });
}

// 新增库存
export function useCreateInventory() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: CreateInventoryRequest) => inventoryService.createInventory(data),
    onSuccess: () => {
      message.success('新增成功');
      queryClient.invalidateQueries({ queryKey: ['inventories'] });
    },
    onError: (err: Error) => {
      message.error(err.message || '新增失败');
    },
  });
}

// 编辑库存
export function useUpdateInventory() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, data }: { id: number; data: UpdateInventoryRequest }) =>
      inventoryService.updateInventory(id, data),
    onSuccess: () => {
      message.success('编辑成功');
      queryClient.invalidateQueries({ queryKey: ['inventories'] });
      queryClient.invalidateQueries({ queryKey: ['inventory'] });
    },
    onError: (err: Error) => {
      message.error(err.message || '编辑失败');
    },
  });
}

// 删除库存
export function useDeleteInventory() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: number) => inventoryService.deleteInventory(id),
    onSuccess: () => {
      message.success('删除成功');
      queryClient.invalidateQueries({ queryKey: ['inventories'] });
    },
    onError: (err: Error) => {
      message.error(err.message || '删除失败');
    },
  });
}

// 批量删除
export function useBatchDeleteInventory() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (ids: number[]) => inventoryService.batchDeleteInventory(ids),
    onSuccess: () => {
      message.success('批量删除成功');
      queryClient.invalidateQueries({ queryKey: ['inventories'] });
    },
    onError: (err: Error) => {
      message.error(err.message || '批量删除失败');
    },
  });
}

// 恢复
export function useRestoreInventory() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: number) => inventoryService.restoreInventory(id),
    onSuccess: () => {
      message.success('恢复成功');
      queryClient.invalidateQueries({ queryKey: ['inventories'] });
    },
    onError: (err: Error) => {
      message.error(err.message || '恢复失败');
    },
  });
}
