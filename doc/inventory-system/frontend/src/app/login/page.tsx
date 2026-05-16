'use client';

import React, { useState } from 'react';
import { Form, Input, Button, Card, Typography, message } from 'antd';
import { UserOutlined, LockOutlined } from '@ant-design/icons';
import { useAuth } from '@/contexts/AuthContext';
import type { LoginRequest } from '@/types';

const { Title, Text } = Typography;

export default function LoginPage() {
  const { login } = useAuth();
  const [loading, setLoading] = useState(false);

  const onFinish = async (values: LoginRequest) => {
    setLoading(true);
    try {
      await login({
        username: values.username.trim(),
        password: values.password.trim(),
      });
    } catch (err: any) {
      message.error(err.message || '登录失败');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 to-indigo-100">
      <Card className="w-[400px] shadow-xl rounded-2xl" variant="borderless">
        <div className="text-center mb-8">
          <Title level={2} className="!mb-2">库存管理系统</Title>
          <Text type="secondary">请登录您的账户</Text>
        </div>

        <Form
          name="login"
          onFinish={onFinish}
          autoComplete="off"
          size="large"
        >
          <Form.Item
            name="username"
            rules={[{ required: true, message: '请输入用户名' }]}
          >
            <Input prefix={<UserOutlined />} placeholder="用户名" />
          </Form.Item>

          <Form.Item
            name="password"
            rules={[{ required: true, message: '请输入密码' }]}
          >
            <Input.Password prefix={<LockOutlined />} placeholder="密码" />
          </Form.Item>

          <Form.Item>
            <Button type="primary" htmlType="submit" loading={loading} block className="h-10">
              登 录
            </Button>
          </Form.Item>
        </Form>

        <div className="text-center text-sm text-gray-400 mt-4">
          <div>测试账号：admin / operator / viewer</div>
          <div>密码：admin123</div>
        </div>
      </Card>
    </div>
  );
}
