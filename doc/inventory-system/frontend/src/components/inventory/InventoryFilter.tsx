'use client';

import React from 'react';
import { Input, Select, Button, Space } from 'antd';
import { SearchOutlined, ReloadOutlined } from '@ant-design/icons';
import { CATEGORY_OPTIONS } from '@/types';

interface FilterValues {
  keyword?: string;
  category?: string;
  status?: number;
}

interface InventoryFilterProps {
  values: FilterValues;
  onChange: (values: FilterValues) => void;
  onReset: () => void;
}

export default function InventoryFilter({ values, onChange, onReset }: InventoryFilterProps) {
  return (
    <div className="bg-white p-4 rounded-lg shadow-sm mb-4">
      <Space wrap size="middle">
        <Input
          placeholder="搜索物料编码/名称"
          prefix={<SearchOutlined />}
          value={values.keyword}
          onChange={(e) => onChange({ ...values, keyword: e.target.value })}
          style={{ width: 220 }}
          allowClear
        />
        <Select
          placeholder="选择分类"
          value={values.category}
          onChange={(val) => onChange({ ...values, category: val })}
          style={{ width: 140 }}
          allowClear
          options={CATEGORY_OPTIONS}
        />
        <Select
          placeholder="状态"
          value={values.status}
          onChange={(val) => onChange({ ...values, status: val })}
          style={{ width: 120 }}
          allowClear
          options={[
            { value: 1, label: '启用' },
            { value: 0, label: '停用' },
          ]}
        />
        <Button icon={<ReloadOutlined />} onClick={onReset}>
          重置
        </Button>
      </Space>
    </div>
  );
}
