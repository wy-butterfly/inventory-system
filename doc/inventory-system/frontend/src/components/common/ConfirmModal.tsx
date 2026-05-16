'use client';

import React from 'react';
import { Modal } from 'antd';
import { ExclamationCircleFilled } from '@ant-design/icons';

interface ConfirmModalProps {
  title?: string;
  content: string;
  onConfirm: () => void;
  onCancel?: () => void;
}

export function showConfirm({ title = '确认操作', content, onConfirm, onCancel }: ConfirmModalProps) {
  Modal.confirm({
    title,
    icon: <ExclamationCircleFilled />,
    content,
    okText: '确认',
    cancelText: '取消',
    onOk: onConfirm,
    onCancel,
  });
}
