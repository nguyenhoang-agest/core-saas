'use client';

import React from 'react';
import { Card, Typography, Alert, Space } from 'antd';
import { BillingStoreProvider } from '../store/billing-store';
import { useTenant } from '@/core/tenant/tenant-context';

const { Title, Paragraph } = Typography;

function Content() {
  const { tenant } = useTenant();

  return (
    <div style={{ maxWidth: 1000, margin: '0 auto', padding: '24px 16px' }}>
      <Space orientation="vertical" size="large" style={{ width: '100%' }}>
        <div>
          <Title level={2}>Billing & Subscriptions</Title>
          <Paragraph type="secondary">
            Manage subscription tiers, payment methods, seats, and download past invoices. Active Tenant: {tenant.name} ({tenant.id})
          </Paragraph>
        </div>

        <Alert
          title="Billing Module Stub"
          description="This module is a skeleton ready for payment gateway integration (e.g. Stripe/Paddle). Business UI has not been implemented yet."
          type="info"
          showIcon
        />

        <Card title="Pending Capabilities">
          <Paragraph>
            Refer to <code>src/modules/billing/README.md</code> for payment flows and integration TODOs.
          </Paragraph>
        </Card>
      </Space>
    </div>
  );
}

export function BillingPage() {
  return (
    <BillingStoreProvider>
      <Content />
    </BillingStoreProvider>
  );
}
