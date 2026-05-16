'use client';

import React from 'react';
import { Upload, message } from 'antd';
import { InboxOutlined } from '@ant-design/icons';
import type { UploadProps } from 'antd';

const { Dragger } = Upload;

interface FileUploaderProps {
  inventoryId: number;
  onSuccess: () => void;
}

export default function FileUploader({ inventoryId, onSuccess }: FileUploaderProps) {
  const token = typeof window !== 'undefined' ? localStorage.getItem('token') : '';

  const props: UploadProps = {
    name: 'file',
    multiple: true,
    action: `/api/v1/inventories/${inventoryId}/attachments`,
    headers: {
      Authorization: `Bearer ${token}`,
    },
    accept: '.jpg,.jpeg,.png,.webp,.pdf,.doc,.docx,.xls,.xlsx,.txt',
    onChange(info) {
      const { status } = info.file;
      if (status === 'done') {
        if (info.file.response?.code === 0) {
          message.success(`${info.file.name} 上传成功`);
          onSuccess();
        } else {
          message.error(info.file.response?.message || '上传失败');
        }
      } else if (status === 'error') {
        message.error(`${info.file.name} 上传失败`);
      }
    },
    showUploadList: false,
  };

  return (
    <Dragger {...props}>
      <p className="text-4xl text-blue-400 mb-2">
        <InboxOutlined />
      </p>
      <p className="text-base">点击或拖拽文件到此区域上传</p>
      <p className="text-gray-400 text-sm mt-1">
        支持图片（JPG/PNG/WebP）和文档（PDF/Word/Excel/TXT），单文件最大20MB
      </p>
    </Dragger>
  );
}
