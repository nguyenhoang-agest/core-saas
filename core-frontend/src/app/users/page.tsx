import React from 'react';
import { UserPage } from '@/modules/user';

export const metadata = {
  title: 'User Management - Core SaaS Admin',
  description: 'Manage tenant users, RBAC roles, and granular ACL security access',
};

export default function UsersPageRoute() {
  return <UserPage />;
}
