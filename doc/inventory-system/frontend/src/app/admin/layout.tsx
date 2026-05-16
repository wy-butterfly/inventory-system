'use client';

import React, { useEffect } from 'react';
import { Layout, Menu, Button, Dropdown, Avatar, Typography } from 'antd';
import {
  DatabaseOutlined,
  UserOutlined,
  FileTextOutlined,
  DashboardOutlined,
  LogoutOutlined,
  MenuFoldOutlined,
  MenuUnfoldOutlined,
} from '@ant-design/icons';
import { useRouter, usePathname } from 'next/navigation';
import { useAuth } from '@/contexts/AuthContext';

const { Header, Sider, Content } = Layout;
const { Text } = Typography;

export default function AdminLayout({ children }: { children: React.ReactNode }) {
  const { user, loading, logout, isAdmin } = useAuth();
  const router = useRouter();
  const pathname = usePathname();
  const [collapsed, setCollapsed] = React.useState(false);

  useEffect(() => {
    if (!loading && !user) {
      router.push('/login');
    }
  }, [user, loading, router]);

  if (loading || !user) return null;

  const menuItems = [
    {
      key: '/admin/inventory',
      icon: <DatabaseOutlined />,
      label: '库存管理',
    },
    ...(isAdmin
      ? [
          {
            key: '/admin/users',
            icon: <UserOutlined />,
            label: '用户管理',
          },
          {
            key: '/admin/logs',
            icon: <FileTextOutlined />,
            label: '操作日志',
          },
        ]
      : []),
  ];

  const userMenuItems = [
    {
      key: 'logout',
      icon: <LogoutOutlined />,
      label: '退出登录',
      onClick: logout,
    },
  ];

  const roleLabel: Record<string, string> = {
    admin: '管理员',
    operator: '操作员',
    viewer: '查看者',
  };

  return (
    <Layout style={{ minHeight: '100vh', height: '100vh' }}>
      <Sider trigger={null} collapsible collapsed={collapsed} theme="dark">
        <div className="h-16 flex items-center justify-center">
          <span style={{ color: '#fff', fontSize: 18, fontWeight: 'bold' }}>
            {collapsed ? '库存' : '库存管理系统'}
          </span>
        </div>
        <Menu
          theme="dark"
          mode="inline"
          selectedKeys={[pathname]}
          items={menuItems}
          onClick={({ key }) => router.push(key)}
        />
      </Sider>
      <Layout>
        <Header className="bg-white px-4 flex items-center justify-between shadow-sm" style={{ padding: '0 24px' }}>
          <Button
            type="text"
            icon={collapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
            onClick={() => setCollapsed(!collapsed)}
          />
          <Dropdown menu={{ items: userMenuItems }} placement="bottomRight">
            <div className="flex items-center gap-2 cursor-pointer">
              <Avatar icon={<UserOutlined />} size="small" />
              <span style={{ color: '#333' }}>{user.display_name}</span>
              <span style={{ color: '#999', fontSize: 12 }}>({roleLabel[user.role]})</span>
            </div>
          </Dropdown>
        </Header>
        <Content className="m-6" style={{ overflow: 'auto' }}>
          {children}
        </Content>
      </Layout>
    </Layout>
  );
}
