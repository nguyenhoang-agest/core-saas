'use client';

import React from 'react';
import { Card, Typography, Alert, Space } from 'antd';
import { TemplateStoreProvider } from '../store/template-store';
import { useTenant } from '@/core/tenant/tenant-context';

const { Title, Paragraph } = Typography;

function Content() {
  const { tenant } = useTenant();

  return (
    <div style={{ maxWidth: 900, margin: '0 auto', padding: '24px 16px' }}>
      <Space orientation="vertical" size="large" style={{ width: '100%' }}>
        <div>
          <Title level={2}>__MODULE_TITLE__</Title>
          <Paragraph type="secondary">
            Scaffolded stub for module <code>__MODULE_NAME__</code>. Tenant: {tenant.name} ({tenant.id})
          </Paragraph>
        </div>

        <Alert
          title="Module Skeleton"
          description="This domain module is registered as a skeleton. Implement domain pages, API integrations, and stores as needed."
          type="info"
          showIcon
        />

        <Card title="Module Overview">
          <Paragraph>
            Refer to <code>src/modules/__MODULE_NAME__/README.md</code> for architecture guidelines and pending TODOs.
          </Paragraph>
        </Card>
      </Space>
    </div>
  );
}

export function TemplatePage() {
  return (
    <TemplateStoreProvider>
      <Content />
    </TemplateStoreProvider>
  );
}
