'use client';

import React, { useState } from 'react';
import { Table, Typography, DatePicker, Space, Tag } from 'antd';
import dayjs from 'dayjs';
import { useQuery } from '@tanstack/react-query';
import api from '@/lib/api';
import type { ApiResponse, PageResult, OperationLog } from '@/types';
import { ACTION_MAP } from '@/types';

const { Title } = Typography;
const { RangePicker } = DatePicker;

export default function LogsPage() {
  const [query, setQuery] = useState<{ page: number; page_size: number; start_time?: string; end_time?: string }>({
    page: 1,
    page_size: 20,
  });

  const { data, isLoading } = useQuery({
    queryKey: ['operation-logs', query],
    queryFn: async () => {
      const res = await api.get<ApiResponse<PageResult<OperationLog>>>('/logs/operations', { params: query });
      return res.data.data;
    },
  });

  const ACTION_COLORS: Record<string, string> = {
    create: 'green',
    update: 'blue',
    delete: 'red',
    login: 'purple',
    export: 'orange',
    upload: 'cyan',
  };

  const columns = [
    { title: 'ID', dataIndex: 'id', width: 60 },
    { title: '操作人', dataIndex: 'user_name', width: 100 },
    {
      title: '操作类型',
      dataIndex: 'action',
      width: 100,
      render: (action: string) => (
        <Tag color={ACTION_COLORS[action] || 'default'}>{ACTION_MAP[action] || action}</Tag>
      ),
    },
    { title: '目标类型', dataIndex: 'target_type', width: 100 },
    { title: '描述', dataIndex: 'description', ellipsis: true },
    { title: 'IP地址', dataIndex: 'ip_address', width: 140 },
    {
      title: '时间',
      dataIndex: 'created_at',
      width: 170,
      render: (t: string) => dayjs(t).format('YYYY-MM-DD HH:mm:ss'),
    },
  ];

  return (
    <div>
      <div className="flex justify-between items-center mb-4">
        <Title level={4} className="!mb-0">操作日志</Title>
        <Space>
          <RangePicker
            onChange={(dates) => {
              if (dates && dates[0] && dates[1]) {
                setQuery({
                  ...query,
                  page: 1,
                  start_time: dates[0].format('YYYY-MM-DD'),
                  end_time: dates[1].format('YYYY-MM-DD'),
                });
              } else {
                const { start_time, end_time, ...rest } = query;
                setQuery({ ...rest, page: 1 });
              }
            }}
          />
        </Space>
      </div>

      <div className="bg-white rounded-lg shadow-sm p-4">
        <Table
          columns={columns}
          dataSource={data?.list || []}
          rowKey="id"
          loading={isLoading}
          pagination={{
            current: query.page,
            pageSize: query.page_size,
            total: data?.total || 0,
            showTotal: (t) => `共 ${t} 条`,
            onChange: (page, pageSize) => setQuery({ ...query, page, page_size: pageSize }),
          }}
        />
      </div>
    </div>
  );
}
