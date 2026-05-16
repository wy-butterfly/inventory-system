'use client';

import React, { useState } from 'react';
import { Button, Space, Typography } from 'antd';
import { PlusOutlined, DeleteOutlined, DownloadOutlined } from '@ant-design/icons';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/contexts/AuthContext';
import { useInventoryList, useDeleteInventory, useBatchDeleteInventory } from '@/hooks/useInventory';
import { exportInventory } from '@/services/inventoryService';
import InventoryFilter from '@/components/inventory/InventoryFilter';
import InventoryTable from '@/components/inventory/InventoryTable';
import { showConfirm } from '@/components/common/ConfirmModal';
import type { Inventory, InventoryQuery } from '@/types';
import { message } from 'antd';

const { Title } = Typography;

export default function InventoryListPage() {
  const router = useRouter();
  const { hasWritePermission } = useAuth();
  const [query, setQuery] = useState<InventoryQuery>({ page: 1, page_size: 20 });
  const [selectedRowKeys, setSelectedRowKeys] = useState<React.Key[]>([]);

  const { data, isLoading } = useInventoryList(query);
  const deleteMutation = useDeleteInventory();
  const batchDeleteMutation = useBatchDeleteInventory();

  const handleFilter = (values: Partial<InventoryQuery>) => {
    setQuery({ ...query, ...values, page: 1 });
  };

  const handleReset = () => {
    setQuery({ page: 1, page_size: 20 });
  };

  const handleView = (record: Inventory) => {
    router.push(`/admin/inventory/${record.id}`);
  };

  const handleEdit = (record: Inventory) => {
    router.push(`/admin/inventory/${record.id}/edit`);
  };

  const handleDelete = (record: Inventory) => {
    showConfirm({
      content: `确定删除物料「${record.material_name}」？删除后可在回收站恢复。`,
      onConfirm: () => deleteMutation.mutate(record.id),
    });
  };

  const handleBatchDelete = () => {
    if (selectedRowKeys.length === 0) {
      message.warning('请先选择要删除的物料');
      return;
    }
    showConfirm({
      content: `确定删除选中的 ${selectedRowKeys.length} 条物料？`,
      onConfirm: () => {
        batchDeleteMutation.mutate(selectedRowKeys as number[], {
          onSuccess: () => setSelectedRowKeys([]),
        });
      },
    });
  };

  const handleExport = async () => {
    try {
      const blob = await exportInventory(query);
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `库存数据_${new Date().toISOString().slice(0, 10)}.xlsx`;
      a.click();
      window.URL.revokeObjectURL(url);
      message.success('导出成功');
    } catch (err: any) {
      message.error(err.message || '导出失败');
    }
  };

  return (
    <div>
      <div className="flex justify-between items-center mb-4">
        <Title level={4} className="!mb-0">库存管理</Title>
        <Space>
          <Button icon={<DownloadOutlined />} onClick={handleExport}>
            导出Excel
          </Button>
          {hasWritePermission && (
            <>
              <Button danger icon={<DeleteOutlined />} onClick={handleBatchDelete} disabled={selectedRowKeys.length === 0}>
                批量删除
              </Button>
              <Button type="primary" icon={<PlusOutlined />} onClick={() => router.push('/admin/inventory/create')}>
                新增物料
              </Button>
            </>
          )}
        </Space>
      </div>

      <InventoryFilter values={query} onChange={handleFilter} onReset={handleReset} />

      <div className="bg-white rounded-lg shadow-sm p-4">
        <InventoryTable
          data={data?.list || []}
          total={data?.total || 0}
          page={query.page || 1}
          pageSize={query.page_size || 20}
          loading={isLoading}
          onPageChange={(page, pageSize) => setQuery({ ...query, page, page_size: pageSize })}
          onView={handleView}
          onEdit={hasWritePermission ? handleEdit : undefined}
          onDelete={hasWritePermission ? handleDelete : undefined}
          selectedRowKeys={hasWritePermission ? selectedRowKeys : undefined}
          onSelectChange={hasWritePermission ? setSelectedRowKeys : undefined}
        />
      </div>
    </div>
  );
}
