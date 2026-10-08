'use client';

import React from 'react';
import { Card, Typography, Alert, Space } from 'antd';
import { IamStoreProvider } from '../store/iam-store';
import { useTenant } from '@/core/tenant/tenant-context';

const { Title, Paragraph } = Typography;

function Content() {
  const { tenant } = useTenant();

  return (
    <div style={{ maxWidth: 1000, margin: '0 auto', padding: '24px 16px' }}>
      <Space orientation="vertical" size="large" style={{ width: '100%' }}>
        <div>
          <Title level={2}>Identity & Access Management (IAM)</Title>
          <Paragraph type="secondary">
            Manage organization members, assign RBAC roles, and inspect audit logs. Active Tenant: {tenant.name} ({tenant.id})
          </Paragraph>
        </div>

        <Alert
          title="IAM Module Stub"
          description="This module is a skeleton ready for domain implementation. Business UI has not been implemented yet."
          type="info"
          showIcon
        />

        <Card title="Pending Capabilities">
          <Paragraph>
            Refer to <code>src/modules/iam/README.md</code> for upcoming features and API specs.
          </Paragraph>
        </Card>
      </Space>
    </div>
  );
}

export function IamPage() {
  return (
    <IamStoreProvider>
      <Content />
    </IamStoreProvider>
  );
}
