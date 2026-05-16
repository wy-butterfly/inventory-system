'use client';

import React, { useState } from 'react';
import { Table, Button, Space, Typography, Tag, Modal, Form, Input, Select, message } from 'antd';
import { PlusOutlined, EditOutlined, KeyOutlined } from '@ant-design/icons';
import dayjs from 'dayjs';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import * as userService from '@/services/userService';
import type { User, CreateUserRequest, UpdateUserRequest } from '@/types';

const { Title } = Typography;

const ROLE_OPTIONS = [
  { value: 'admin', label: '管理员' },
  { value: 'operator', label: '操作员' },
  { value: 'viewer', label: '查看者' },
];

const ROLE_COLORS: Record<string, string> = {
  admin: 'red',
  operator: 'blue',
  viewer: 'green',
};

const ROLE_LABELS: Record<string, string> = {
  admin: '管理员',
  operator: '操作员',
  viewer: '查看者',
};

export default function UsersPage() {
  const queryClient = useQueryClient();
  const [query, setQuery] = useState({ page: 1, page_size: 20 });
  const [modalOpen, setModalOpen] = useState(false);
  const [editingUser, setEditingUser] = useState<User | null>(null);
  const [resetPwdOpen, setResetPwdOpen] = useState(false);
  const [resetUserId, setResetUserId] = useState<number>(0);
  const [form] = Form.useForm();
  const [pwdForm] = Form.useForm();

  const { data, isLoading } = useQuery({
    queryKey: ['users', query],
    queryFn: () => userService.getUserList(query),
  });

  const createMutation = useMutation({
    mutationFn: (data: CreateUserRequest) => userService.createUser(data),
    onSuccess: () => {
      message.success('创建成功');
      queryClient.invalidateQueries({ queryKey: ['users'] });
      setModalOpen(false);
      form.resetFields();
    },
    onError: (err: Error) => message.error(err.message),
  });

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: number; data: UpdateUserRequest }) => userService.updateUser(id, data),
    onSuccess: () => {
      message.success('更新成功');
      queryClient.invalidateQueries({ queryKey: ['users'] });
      setModalOpen(false);
      setEditingUser(null);
      form.resetFields();
    },
    onError: (err: Error) => message.error(err.message),
  });

  const resetPwdMutation = useMutation({
    mutationFn: ({ id, password }: { id: number; password: string }) => userService.resetPassword(id, password),
    onSuccess: () => {
      message.success('密码重置成功');
      setResetPwdOpen(false);
      pwdForm.resetFields();
    },
    onError: (err: Error) => message.error(err.message),
  });

  const handleAdd = () => {
    setEditingUser(null);
    form.resetFields();
    setModalOpen(true);
  };

  const handleEdit = (record: User) => {
    setEditingUser(record);
    form.setFieldsValue(record);
    setModalOpen(true);
  };

  const handleModalOk = () => {
    form.validateFields().then((values) => {
      if (editingUser) {
        updateMutation.mutate({ id: editingUser.id, data: values });
      } else {
        createMutation.mutate(values);
      }
    });
  };

  const columns = [
    { title: 'ID', dataIndex: 'id', width: 60 },
    { title: '用户名', dataIndex: 'username', width: 120 },
    { title: '显示名', dataIndex: 'display_name', width: 140 },
    {
      title: '角色',
      dataIndex: 'role',
      width: 100,
      render: (role: string) => <Tag color={ROLE_COLORS[role]}>{ROLE_LABELS[role]}</Tag>,
    },
    {
      title: '状态',
      dataIndex: 'status',
      width: 80,
      render: (status: number) => status === 1 ? <Tag color="green">启用</Tag> : <Tag color="red">停用</Tag>,
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      width: 160,
      render: (t: string) => dayjs(t).format('YYYY-MM-DD HH:mm'),
    },
    {
      title: '操作',
      key: 'action',
      width: 150,
      render: (_: any, record: User) => (
        <Space size="small">
          <Button type="link" size="small" icon={<EditOutlined />} onClick={() => handleEdit(record)}>
            编辑
          </Button>
          <Button
            type="link"
            size="small"
            icon={<KeyOutlined />}
            onClick={() => { setResetUserId(record.id); setResetPwdOpen(true); }}
          >
            重置密码
          </Button>
        </Space>
      ),
    },
  ];

  return (
    <div>
      <div className="flex justify-between items-center mb-4">
        <Title level={4} className="!mb-0">用户管理</Title>
        <Button type="primary" icon={<PlusOutlined />} onClick={handleAdd}>
          新增用户
        </Button>
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
            onChange: (page, pageSize) => setQuery({ page, page_size: pageSize }),
          }}
        />
      </div>

      {/* 新增/编辑用户弹窗 */}
      <Modal
        title={editingUser ? '编辑用户' : '新增用户'}
        open={modalOpen}
        onOk={handleModalOk}
        onCancel={() => { setModalOpen(false); setEditingUser(null); form.resetFields(); }}
        confirmLoading={createMutation.isPending || updateMutation.isPending}
      >
        <Form form={form} layout="vertical">
          <Form.Item
            label="用户名"
            name="username"
            rules={[{ required: true, message: '请输入用户名' }]}
          >
            <Input disabled={!!editingUser} placeholder="3-50个字符" />
          </Form.Item>
          {!editingUser && (
            <Form.Item
              label="密码"
              name="password"
              rules={[{ required: true, message: '请输入密码' }, { min: 6, message: '至少6个字符' }]}
            >
              <Input.Password placeholder="至少6个字符" />
            </Form.Item>
          )}
          <Form.Item label="显示名" name="display_name">
            <Input placeholder="用于界面展示" />
          </Form.Item>
          <Form.Item
            label="角色"
            name="role"
            rules={[{ required: true, message: '请选择角色' }]}
          >
            <Select options={ROLE_OPTIONS} placeholder="请选择角色" />
          </Form.Item>
          {editingUser && (
            <Form.Item label="状态" name="status">
              <Select
                options={[
                  { value: 1, label: '启用' },
                  { value: 0, label: '停用' },
                ]}
              />
            </Form.Item>
          )}
        </Form>
      </Modal>

      {/* 重置密码弹窗 */}
      <Modal
        title="重置密码"
        open={resetPwdOpen}
        onOk={() => {
          pwdForm.validateFields().then((values) => {
            resetPwdMutation.mutate({ id: resetUserId, password: values.new_password });
          });
        }}
        onCancel={() => { setResetPwdOpen(false); pwdForm.resetFields(); }}
        confirmLoading={resetPwdMutation.isPending}
      >
        <Form form={pwdForm} layout="vertical">
          <Form.Item
            label="新密码"
            name="new_password"
            rules={[{ required: true, message: '请输入新密码' }, { min: 6, message: '至少6个字符' }]}
          >
            <Input.Password placeholder="请输入新密码" />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
}
