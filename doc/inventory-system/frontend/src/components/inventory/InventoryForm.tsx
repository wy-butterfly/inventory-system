'use client';

import React from 'react';
import { Form, Input, InputNumber, Select, Button, Space } from 'antd';
import { CATEGORY_OPTIONS } from '@/types';
import type { CreateInventoryRequest } from '@/types';

interface InventoryFormProps {
  initialValues?: Partial<CreateInventoryRequest & { version?: number }>;
  onSubmit: (values: any) => void;
  onCancel: () => void;
  loading?: boolean;
  isEdit?: boolean;
}

export default function InventoryForm({
  initialValues,
  onSubmit,
  onCancel,
  loading = false,
  isEdit = false,
}: InventoryFormProps) {
  const [form] = Form.useForm();

  const handleFinish = (values: any) => {
    if (isEdit && initialValues?.version) {
      values.version = initialValues.version;
    }
    onSubmit(values);
  };

  return (
    <Form
      form={form}
      layout="vertical"
      initialValues={{ status: 1, ...initialValues }}
      onFinish={handleFinish}
      className="max-w-2xl"
    >
      <div className="grid grid-cols-2 gap-x-6">
        <Form.Item
          label="物料编码"
          name="material_code"
          rules={[{ required: true, message: '请输入物料编码' }]}
        >
          <Input placeholder="如: YL-202603-0001" disabled={isEdit} maxLength={50} />
        </Form.Item>

        <Form.Item
          label="物料名称"
          name="material_name"
          rules={[{ required: true, message: '请输入物料名称' }]}
        >
          <Input placeholder="请输入物料名称" maxLength={200} />
        </Form.Item>

        <Form.Item
          label="分类"
          name="category"
          rules={[{ required: true, message: '请选择分类' }]}
        >
          <Select placeholder="请选择分类" options={CATEGORY_OPTIONS} />
        </Form.Item>

        <Form.Item label="规格型号" name="specification">
          <Input placeholder="请输入规格型号" maxLength={200} />
        </Form.Item>

        <Form.Item
          label="单位"
          name="unit"
          rules={[{ required: true, message: '请输入单位' }]}
        >
          <Input placeholder="如: 个、台、米、张" maxLength={20} />
        </Form.Item>

        <Form.Item
          label="库存数量"
          name="quantity"
          rules={[{ required: true, message: '请输入库存数量' }]}
        >
          <InputNumber min={0} precision={2} style={{ width: '100%' }} placeholder="0.00" />
        </Form.Item>

        <Form.Item label="安全库存" name="safety_stock">
          <InputNumber min={0} precision={2} style={{ width: '100%' }} placeholder="低于此值会预警" />
        </Form.Item>

        <Form.Item label="仓库位置" name="location">
          <Input placeholder="如: A区-01-01" maxLength={200} />
        </Form.Item>

        <Form.Item label="状态" name="status">
          <Select
            options={[
              { value: 1, label: '启用' },
              { value: 0, label: '停用' },
            ]}
          />
        </Form.Item>
      </div>

      <Form.Item label="备注" name="remark">
        <Input.TextArea rows={3} placeholder="备注信息（可选）" maxLength={1000} showCount />
      </Form.Item>

      <Form.Item>
        <Space>
          <Button type="primary" htmlType="submit" loading={loading}>
            {isEdit ? '保存修改' : '确认新增'}
          </Button>
          <Button onClick={onCancel}>取消</Button>
        </Space>
      </Form.Item>
    </Form>
  );
}
