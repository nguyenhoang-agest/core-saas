import React from 'react';
import { TenantAdminPage } from '@/modules/tenant';

export const metadata = {
  title: 'Tenants Management - Core SaaS Admin',
  description: 'Manage tenant workspaces, database isolation modes, and runtime contexts',
};

export default function TenantsPage() {
  return <TenantAdminPage />;
}
