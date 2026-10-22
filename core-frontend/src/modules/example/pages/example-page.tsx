'use client';

import React from 'react';
import { Card, Table, Typography, Input, Button, Space, Tag } from 'antd';
import { ExampleStoreProvider, useExampleStore } from '../store/example-store';
import { useExampleItems } from '../hooks/example-hooks';
import type { ExampleItem } from '../api/example-api';
import { Can } from '@/core/auth/can';
import { useTenant } from '@/core/tenant/tenant-context';

const { Title, Paragraph, Text } = Typography;

function ExampleContent() {
  const { tenant } = useTenant();
  const filter = useExampleStore((s) => s.filter);
  const setFilter = useExampleStore((s) => s.setFilter);
  const counter = useExampleStore((s) => s.counter);
  const increment = useExampleStore((s) => s.increment);

  const { data: items = [], isLoading } = useExampleItems();

  const filteredItems = items.filter((item) =>
    item.title.toLowerCase().includes(filter.toLowerCase())
  );

  const columns = [
    {
      title: 'Title',
      dataIndex: 'title',
      key: 'title',
      render: (text: string) => <Text strong>{text}</Text>,
    },
    {
      title: 'Status',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => {
        const color = status === 'completed' ? 'green' : status === 'pending' ? 'gold' : 'blue';
        return <Tag color={color}>{status.toUpperCase()}</Tag>;
      },
    },
    {
      title: 'Created At',
      dataIndex: 'createdAt',
      key: 'createdAt',
    },
  ];

  return (
    <div style={{ maxWidth: 1000, margin: '0 auto', padding: '24px 16px' }}>
      <Space orientation="vertical" size="large" style={{ width: '100%' }}>
        <div>
          <Title level={2}>Example Module Slice</Title>
          <Paragraph type="secondary">
            A thin vertical slice demonstrating integration across Ant Design, React Query
            (tenant-scoped), Zustand (context-safe store factory), and UI permission guards. Active
            Tenant:{' '}
            <Tag color="blue">
              {tenant.name} ({tenant.id})
            </Tag>
          </Paragraph>
        </div>

        <Card title="Module Zustand Store (Context-Bound)">
          <Space orientation="vertical" size="middle" style={{ width: '100%' }}>
            <Space wrap>
              <Input
                placeholder="Filter items by title..."
                value={filter}
                onChange={(e) => setFilter(e.target.value)}
                style={{ width: 280 }}
                allowClear
              />
              <Button onClick={increment} type="default">
                Store Action Clicks: {counter}
              </Button>
            </Space>

            <Can
              permission="example:write"
              fallback={
                <Text type="secondary">
                  (Protected mutation button hidden - requires <code>example:write</code>{' '}
                  permission)
                </Text>
              }
            >
              <Button type="primary">Protected Mutation Action</Button>
            </Can>
          </Space>
        </Card>

        <Card title="Tenant-Scoped API Data (React Query)">
          <Table<ExampleItem>
            rowKey="id"
            columns={columns}
            dataSource={filteredItems}
            loading={isLoading}
            pagination={false}
          />
        </Card>
      </Space>
    </div>
  );
}

export function ExamplePage() {
  return (
    <ExampleStoreProvider>
      <ExampleContent />
    </ExampleStoreProvider>
  );
}
