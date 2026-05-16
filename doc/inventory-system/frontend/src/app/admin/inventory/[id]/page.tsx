'use client';

import React from 'react';
import { Typography, Descriptions, Tag, Spin, Button, Space, Card, Tabs, Table, Image, Empty } from 'antd';
import { EditOutlined, ArrowLeftOutlined, WarningOutlined } from '@ant-design/icons';
import { useRouter, useParams } from 'next/navigation';
import dayjs from 'dayjs';
import { useAuth } from '@/contexts/AuthContext';
import { useInventoryDetail } from '@/hooks/useInventory';
import { useAttachmentList, useDeleteAttachment } from '@/hooks/useAttachment';
import { showConfirm } from '@/components/common/ConfirmModal';
import { CATEGORY_MAP } from '@/types';
import { getPreviewUrl, getDownloadUrl } from '@/services/attachmentService';
import FileUploader from '@/components/attachment/FileUploader';
import type { Attachment } from '@/types';

const { Title } = Typography;

export default function InventoryDetailPage() {
  const router = useRouter();
  const params = useParams();
  const id = Number(params.id);
  const { hasWritePermission } = useAuth();

  const { data: inventory, isLoading } = useInventoryDetail(id);
  const { data: attachments, refetch: refetchAttachments } = useAttachmentList(id);
  const deleteAttachmentMutation = useDeleteAttachment();

  const handleDeleteAttachment = (record: Attachment) => {
    showConfirm({
      title: '删除附件',
      content: `确定删除附件「${record.display_name}」？`,
      onConfirm: () => deleteAttachmentMutation.mutate(record.id),
    });
  };

  if (isLoading) {
    return <div className="flex justify-center py-20"><Spin size="large" /></div>;
  }

  if (!inventory) {
    return <Empty description="物料不存在" />;
  }

  const attachmentColumns = [
    {
      title: '预览',
      dataIndex: 'file_type',
      width: 80,
      render: (_: string, record: Attachment) => {
        if (record.file_type === 'image') {
          return (
            <Image 
              src={getPreviewUrl(record.id, true)} 
              width={48} 
              height={48} 
              alt={record.display_name} 
              style={{ objectFit: 'cover', borderRadius: 4, cursor: 'pointer' }}
              onClick={() => window.open(getPreviewUrl(record.id, false), '_blank')}
            />
          );
        }
        
        // 非图片文件显示图标，点击尝试预览
        return (
          <div 
            className="w-12 h-12 bg-gray-100 flex items-center justify-center rounded text-xs text-gray-500 cursor-pointer hover:bg-gray-200"
            onClick={() => {
              const previewUrl = getPreviewUrl(record.id, false);
              // 尝试在新窗口打开预览，如果失败则下载
              window.open(previewUrl, '_blank') || window.open(getDownloadUrl(record.id), '_blank');
            }}
          >
            {record.extension.toUpperCase()}
          </div>
        );
      },
    },
    { title: '文件名', dataIndex: 'display_name', ellipsis: true },
    {
      title: '大小',
      dataIndex: 'file_size',
      width: 100,
      render: (size: number) => {
        if (size < 1024) return `${size} B`;
        if (size < 1024 * 1024) return `${(size / 1024).toFixed(1)} KB`;
        return `${(size / 1024 / 1024).toFixed(1)} MB`;
      },
    },
    {
      title: '上传时间',
      dataIndex: 'created_at',
      width: 160,
      render: (t: string) => dayjs(t).format('YYYY-MM-DD HH:mm'),
    },
    {
      title: '操作',
      key: 'action',
      width: 150,
      render: (_: any, record: Attachment) => (
        <Space size="small">
          <a href={getDownloadUrl(record.id)} target="_blank" rel="noopener noreferrer">下载</a>
          {hasWritePermission && (
            <a style={{ color: '#ff4d4f' }} onClick={() => handleDeleteAttachment(record)}>删除</a>
          )}
        </Space>
      ),
    },
  ];

  return (
    <div>
      <div className="flex justify-between items-center mb-4">
        <Space>
          <Button icon={<ArrowLeftOutlined />} onClick={() => router.back()}>返回</Button>
          <Title level={4} className="!mb-0">物料详情</Title>
        </Space>
        {hasWritePermission && (
          <Button type="primary" icon={<EditOutlined />} onClick={() => router.push(`/admin/inventory/${id}/edit`)}>
            编辑
          </Button>
        )}
      </div>

      <Card className="mb-4">
        <Descriptions title="基本信息" bordered column={2}>
          <Descriptions.Item label="物料编码">
            {inventory.material_code}
            {inventory.warning && <WarningOutlined style={{ color: '#faad14', marginLeft: 8 }} />}
          </Descriptions.Item>
          <Descriptions.Item label="物料名称">{inventory.material_name}</Descriptions.Item>
          <Descriptions.Item label="分类">{CATEGORY_MAP[inventory.category] || inventory.category}</Descriptions.Item>
          <Descriptions.Item label="规格型号">{inventory.specification || '-'}</Descriptions.Item>
          <Descriptions.Item label="库存数量">
            <span style={{ color: inventory.warning ? '#ff4d4f' : undefined, fontWeight: inventory.warning ? 'bold' : undefined }}>
              {inventory.quantity} {inventory.unit}
            </span>
          </Descriptions.Item>
          <Descriptions.Item label="安全库存">{inventory.safety_stock} {inventory.unit}</Descriptions.Item>
          <Descriptions.Item label="仓库位置">{inventory.location || '-'}</Descriptions.Item>
          <Descriptions.Item label="状态">
            {inventory.status === 1 ? <Tag color="green">启用</Tag> : <Tag color="red">停用</Tag>}
          </Descriptions.Item>
          <Descriptions.Item label="创建人">{inventory.creator_name || '-'}</Descriptions.Item>
          <Descriptions.Item label="创建时间">{dayjs(inventory.created_at).format('YYYY-MM-DD HH:mm:ss')}</Descriptions.Item>
          <Descriptions.Item label="更新时间" span={2}>{dayjs(inventory.updated_at).format('YYYY-MM-DD HH:mm:ss')}</Descriptions.Item>
          <Descriptions.Item label="备注" span={2}>{inventory.remark || '-'}</Descriptions.Item>
        </Descriptions>
      </Card>

      <Card>
        <Tabs
          items={[
            {
              key: 'attachments',
              label: `附件 (${attachments?.length || 0})`,
              children: (
                <div>
                  {hasWritePermission && (
                    <div className="mb-4">
                      <FileUploader inventoryId={id} onSuccess={() => refetchAttachments()} />
                    </div>
                  )}
                  <Table
                    columns={attachmentColumns}
                    dataSource={attachments || []}
                    rowKey="id"
                    pagination={false}
                    locale={{ emptyText: '暂无附件' }}
                  />
                </div>
              ),
            },
          ]}
        />
      </Card>
    </div>
  );
}
