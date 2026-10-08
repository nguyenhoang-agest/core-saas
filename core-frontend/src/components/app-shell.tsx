'use client';

import React, { useState } from 'react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { Layout, Menu, Typography, Space, Button } from 'antd';
import { Slot } from '@/core/slots/slot';
import { moduleRegistry } from '@/core/modules/registry';
import { appConfig } from '@/../config/app.config';
import { useTenant } from '@/core/tenant/tenant-context';

const { Header, Content, Footer } = Layout;
const { Title, Text } = Typography;

export function AppShell({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();
  const { tenant } = useTenant();
  const menuItems = moduleRegistry.getMenuItems();

  const navItems = [
    {
      key: '/',
      label: <Link href="/">Dashboard</Link>,
    },
    ...menuItems.map((item) => ({
      key: item.path,
      label: <Link href={item.path}>{item.label}</Link>,
    })),
  ];

  return (
    <Layout style={{ minHeight: '100vh', display: 'flex', flexDirection: 'column' }}>
      <Header
        style={{
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'space-between',
          background: '#ffffff',
          padding: '0 24px',
          borderBottom: '1px solid #f0f0f0',
          position: 'sticky',
          top: 0,
          zIndex: 100,
        }}
      >
        <Space size="large" align="center">
          <Link href="/" style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
            <div
              style={{
                width: 28,
                height: 28,
                borderRadius: 6,
                background: 'var(--primary-color, #1677ff)',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                color: '#fff',
                fontWeight: 'bold',
                fontSize: 16,
              }}
            >
              C
            </div>
            <Title level={4} style={{ margin: 0, whiteSpace: 'nowrap' }}>
              {appConfig.branding.appName}
            </Title>
          </Link>

          <Menu
            mode="horizontal"
            selectedKeys={[pathname || '/']}
            items={navItems}
            style={{ borderBottom: 'none', minWidth: 280 }}
          />
        </Space>

        <Space size="middle" align="center">
          <Text type="secondary" style={{ fontSize: 12 }}>
            Tenant: <strong>{tenant.id}</strong>
          </Text>
          <Slot name="header-actions" />
        </Space>
      </Header>

      <Content style={{ flex: 1, background: '#f5f7fa' }}>{children}</Content>

      <Footer style={{ textAlign: 'center', background: '#ffffff', borderTop: '1px solid #f0f0f0', padding: '16px 24px' }}>
        <Text type="secondary" style={{ fontSize: 13 }}>
          {appConfig.branding.appName} &copy; {new Date().getFullYear()} &mdash; SaaS Starter Kit
        </Text>
      </Footer>
    </Layout>
  );
}
