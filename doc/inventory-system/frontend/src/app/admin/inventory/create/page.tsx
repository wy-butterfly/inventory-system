'use client';

import React from 'react';
import { Typography } from 'antd';
import { useRouter } from 'next/navigation';
import { useCreateInventory } from '@/hooks/useInventory';
import InventoryForm from '@/components/inventory/InventoryForm';

const { Title } = Typography;

export default function CreateInventoryPage() {
  const router = useRouter();
  const createMutation = useCreateInventory();

  const handleSubmit = (values: any) => {
    createMutation.mutate(values, {
      onSuccess: () => router.push('/admin/inventory'),
    });
  };

  return (
    <div className="bg-white rounded-lg shadow-sm p-6">
      <Title level={4}>新增物料</Title>
      <InventoryForm
        onSubmit={handleSubmit}
        onCancel={() => router.back()}
        loading={createMutation.isPending}
      />
    </div>
  );
}
