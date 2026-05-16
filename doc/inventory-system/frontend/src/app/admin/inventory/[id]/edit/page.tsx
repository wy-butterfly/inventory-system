'use client';

import React from 'react';
import { Typography, Spin, Empty } from 'antd';
import { useRouter, useParams } from 'next/navigation';
import { useInventoryDetail, useUpdateInventory } from '@/hooks/useInventory';
import InventoryForm from '@/components/inventory/InventoryForm';

const { Title } = Typography;

export default function EditInventoryPage() {
  const router = useRouter();
  const params = useParams();
  const id = Number(params.id);

  const { data: inventory, isLoading } = useInventoryDetail(id);
  const updateMutation = useUpdateInventory();

  if (isLoading) {
    return <div className="flex justify-center py-20"><Spin size="large" /></div>;
  }

  if (!inventory) {
    return <Empty description="物料不存在" />;
  }

  const handleSubmit = (values: any) => {
    updateMutation.mutate(
      { id, data: values },
      { onSuccess: () => router.push('/admin/inventory') }
    );
  };

  return (
    <div className="bg-white rounded-lg shadow-sm p-6">
      <Title level={4}>编辑物料</Title>
      <InventoryForm
        initialValues={inventory}
        onSubmit={handleSubmit}
        onCancel={() => router.push('/admin/inventory')}
        loading={updateMutation.isPending}
        isEdit
      />
    </div>
  );
}
