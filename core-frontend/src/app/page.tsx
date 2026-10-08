'use client';

import React from 'react';
import Link from 'next/link';
import { Card, Typography, Space, Tag, Button, Row, Col } from 'antd';

const { Title, Paragraph, Text } = Typography;

export default function HomePage() {
  const stackItems = [
    { label: 'Next.js', version: '16 (App Router)' },
    { label: 'React', version: '19' },
    { label: 'TypeScript', version: 'Strict Mode' },
    { label: 'Ant Design', version: 'v6 (@ant-design/cssinjs)' },
    { label: 'State & Query', version: 'TanStack Query v5 & Zustand v5' },
    { label: 'Auth & Forms', version: 'NextAuth v4, RHF, Zod' },
    { label: 'Testing & Tools', version: 'Vitest, ESLint, Prettier' },
  ];

  return (
    <main style={{ minHeight: 'calc(100vh - 120px)', padding: '48px 24px' }}>
      <Row justify="center">
        <Col xs={24} sm={20} md={16} lg={14}>
          <Card
            variant="borderless"
            style={{
              boxShadow: '0 4px 12px rgba(0, 0, 0, 0.05)',
              borderRadius: 8,
            }}
          >
            <Space orientation="vertical" size="large" style={{ width: '100%' }}>
              <div>
                <Title level={2} style={{ marginBottom: 8 }}>
                  Core SaaS Frontend
                </Title>
                <Paragraph type="secondary">
                  Modular SaaS Starter Kit scaffolded with Next.js 16, Ant Design v6, and isolated domain modules.
                </Paragraph>
              </div>

              <div>
                <Text strong style={{ display: 'block', marginBottom: 12 }}>
                  Configured Tech Stack
                </Text>
                <Space wrap size={[8, 8]}>
                  {stackItems.map((item) => (
                    <Tag key={item.label} color="blue">
                      {item.label}: {item.version}
                    </Tag>
                  ))}
                </Space>
              </div>

              <div>
                <Text strong style={{ display: 'block', marginBottom: 12 }}>
                  Explore Starter Kit Slices
                </Text>
                <Space wrap>
                  <Link href="/example">
                    <Button type="primary">
                      Open Example Module Slice &rarr;
                    </Button>
                  </Link>
                  <Button href="https://nextjs.org/docs" target="_blank">
                    Next.js Docs
                  </Button>
                  <Button href="https://ant.design/components/overview" target="_blank">
                    Ant Design Docs
                  </Button>
                </Space>
              </div>
            </Space>
          </Card>
        </Col>
      </Row>
    </main>
  );
}
