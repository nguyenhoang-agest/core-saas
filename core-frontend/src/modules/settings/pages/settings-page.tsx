'use client';

import React from 'react';
import { Card, Typography, Alert, Space } from 'antd';
import { SettingsStoreProvider } from '../store/settings-store';
import { useTenant } from '@/core/tenant/tenant-context';

const { Title, Paragraph } = Typography;

function Content() {
  const { tenant } = useTenant();

  return (
    <div style={{ maxWidth: 1000, margin: '0 auto', padding: '24px 16px' }}>
      <Space orientation="vertical" size="large" style={{ width: '100%' }}>
        <div>
          <Title level={2}>Tenant Settings</Title>
          <Paragraph type="secondary">
            Configure tenant branding, primary colors, subdomains, and notification preferences. Active Tenant: {tenant.name} ({tenant.id})
          </Paragraph>
        </div>

        <Alert
          title="Settings Module Stub"
          description="This module is a skeleton ready for domain implementation. Business UI has not been implemented yet."
          type="info"
          showIcon
        />

        <Card title="Pending Capabilities">
          <Paragraph>
            Refer to <code>src/modules/settings/README.md</code> for full scope and configuration options.
          </Paragraph>
        </Card>
      </Space>
    </div>
  );
}

export function SettingsPage() {
  return (
    <SettingsStoreProvider>
      <Content />
    </SettingsStoreProvider>
  );
}
