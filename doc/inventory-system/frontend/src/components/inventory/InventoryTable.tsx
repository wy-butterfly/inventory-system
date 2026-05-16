'use client';

import React from 'react';
import { Table, Tag, Space, Button, Tooltip } from 'antd';
import { EyeOutlined, EditOutlined, DeleteOutlined, WarningOutlined } from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import type { Inventory } from '@/types';
import { CATEGORY_MAP } from '@/types';
import dayjs from 'dayjs';

interface InventoryTableProps {
  data: Inventory[];
  total: number;
  page: number;
  pageSize: number;
  loading: boolean;
  onPageChange: (page: number, pageSize: number) => void;
  onView: (record: Inventory) => void;
  onEdit?: (record: Inventory) => void;
  onDelete?: (record: Inventory) => void;
  showActions?: boolean;
  selectedRowKeys?: React.Key[];
  onSelectChange?: (keys: React.Key[]) => void;
}

export default function InventoryTable({
  data,
  total,
  page,
  pageSize,
  loading,
  onPageChange,
  onView,
  onEdit,
  onDelete,
  showActions = true,
  selectedRowKeys,
  onSelectChange,
}: InventoryTableProps) {
  const columns: ColumnsType<Inventory> = [
    {
      title: '物料编码',
      dataIndex: 'material_code',
      width: 160,
      render: (code: string, record) => (
        <Space>
          <a onClick={() => onView(record)}>{code}</a>
          {record.warning && (
            <Tooltip title="库存低于安全库存">
              <WarningOutlined style={{ color: '#faad14' }} />
            </Tooltip>
          )}
        </Space>
      ),
    },
    {
      title: '物料名称',
      dataIndex: 'material_name',
      width: 180,
      ellipsis: true,
    },
    {
      title: '分类',
      dataIndex: 'category',
      width: 80,
      render: (cat: string) => CATEGORY_MAP[cat] || cat,
    },
    {
      title: '规格型号',
      dataIndex: 'specification',
      width: 150,
      ellipsis: true,
    },
    {
      title: '库存数量',
      dataIndex: 'quantity',
      width: 100,
      render: (qty: number, record) => (
        <span style={{ color: record.warning ? '#ff4d4f' : undefined, fontWeight: record.warning ? 'bold' : undefined }}>
          {qty} {record.unit}
        </span>
      ),
    },
    {
      title: '安全库存',
      dataIndex: 'safety_stock',
      width: 100,
      render: (val: number, record) => `${val} ${record.unit}`,
    },
    {
      title: '仓库位置',
      dataIndex: 'location',
      width: 120,
      ellipsis: true,
    },
    {
      title: '状态',
      dataIndex: 'status',
      width: 80,
      render: (status: number) =>
        status === 1 ? <Tag color="green">启用</Tag> : <Tag color="red">停用</Tag>,
    },
    {
      title: '更新时间',
      dataIndex: 'updated_at',
      width: 160,
      render: (t: string) => dayjs(t).format('YYYY-MM-DD HH:mm'),
    },
  ];

  if (showActions) {
    columns.push({
      title: '操作',
      key: 'action',
      width: 150,
      fixed: 'right',
      render: (_, record) => (
        <Space size="small">
          <Button type="link" size="small" icon={<EyeOutlined />} onClick={() => onView(record)}>
            查看
          </Button>
          {onEdit && (
            <Button type="link" size="small" icon={<EditOutlined />} onClick={() => onEdit(record)}>
              编辑
            </Button>
          )}
          {onDelete && (
            <Button type="link" size="small" danger icon={<DeleteOutlined />} onClick={() => onDelete(record)}>
              删除
            </Button>
          )}
        </Space>
      ),
    });
  }

  return (
    <Table
      columns={columns}
      dataSource={data}
      rowKey="id"
      loading={loading}
      scroll={{ x: 1200 }}
      pagination={{
        current: page,
        pageSize,
        total,
        showSizeChanger: true,
        showQuickJumper: true,
        showTotal: (t) => `共 ${t} 条记录`,
        onChange: onPageChange,
      }}
      rowSelection={
        onSelectChange
          ? {
              selectedRowKeys,
              onChange: onSelectChange,
            }
          : undefined
      }
    />
  );
}
