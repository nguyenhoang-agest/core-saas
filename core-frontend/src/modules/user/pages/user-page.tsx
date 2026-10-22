'use client';

import React, { useState } from 'react';
import { Users as UsersIcon, Plus, CheckCircle2, Shield, Key } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { useI18n } from '@/core/i18n/i18n-context';
import { useTenant } from '@/core/tenant/tenant-context';
import { useSession } from 'next-auth/react';
import type { CustomSession } from '@/core/auth/auth-options';
import { UserStoreProvider } from '../store/user-store';
import {
  UserKpiCards,
  UserListTable,
  CreateUserDialog,
  EditUserDialog,
  RoleListTable,
  RolePermissionMatrix,
  UserInspector,
} from '../components';
import { useUsers } from '../hooks';
import type { User } from '../types';

function UserContent() {
  const { t } = useI18n();
  const { tenant } = useTenant();
  const { data: session } = useSession();
  const customSession = session as CustomSession | null;
  const actorPermissions = customSession?.user?.permissions || ['*'];

  const [activeTab, setActiveTab] = useState<'users' | 'matrix'>('users');
  const [createDialogOpen, setCreateDialogOpen] = useState(false);
  const [editDialogOpen, setEditDialogOpen] = useState(false);
  const [selectedUser, setSelectedUser] = useState<User | null>(null);
  const [editingUser, setEditingUser] = useState<User | null>(null);
  const [newlyCreatedUser, setNewlyCreatedUser] = useState<User | null>(null);

  // Load user data for KPI metrics
  const { data: usersData, isLoading: usersLoading } = useUsers({ page: 1, pageSize: 50 });

  const handleUserCreated = (user: User) => {
    setNewlyCreatedUser(user);
    setSelectedUser(user);
  };

  const handleInspectUser = (user: User) => {
    setSelectedUser(user);
    const el = document.getElementById('user-inspector-section');
    if (el) {
      el.scrollIntoView({ behavior: 'smooth' });
    }
  };

  const handleEditUser = (user: User) => {
    setEditingUser(user);
    setEditDialogOpen(true);
  };

  return (
    <div className="space-y-6">
      {/* Page Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between pb-5 border-b border-neutral-200 gap-4">
        <div>
          <div className="flex items-center gap-2 mb-1">
            <div className="p-1.5 rounded-md border border-neutral-200 bg-white shadow-none">
              <UsersIcon className="h-5 w-5 text-neutral-900 stroke-[1.75]" />
            </div>
            <h1 className="text-xl font-bold tracking-tight text-neutral-900">
              {t('user.page.title', t('iam.page.title', 'User Management'))}
            </h1>
          </div>
          <p className="text-xs text-neutral-500 max-w-2xl">
            {t(
              'user.page.description',
              t(
                'iam.page.description',
                'Manage organization members, assign RBAC roles, configure granular ACL permissions, and inspect security access.'
              )
            )}{' '}
            <span className="text-neutral-700 font-medium">
              Active Tenant: {tenant.name} ({tenant.id})
            </span>
          </p>
        </div>

        {/* Primary Action Button */}
        <Button
          onClick={() => setCreateDialogOpen(true)}
          className="gap-1.5 self-start sm:self-auto shrink-0 text-xs"
        >
          <Plus className="h-4 w-4 stroke-[2]" />
          <span>{t('user.page.create_btn', t('iam.page.invite_btn', 'Create User'))}</span>
        </Button>
      </div>

      {/* Success Notification Banner after User Creation */}
      {newlyCreatedUser && (
        <div className="rounded-lg border border-emerald-200 bg-emerald-50/60 p-4 text-xs text-emerald-900 flex items-start justify-between gap-3 duration-200 animate-in fade-in-0">
          <div className="flex items-start gap-2.5">
            <CheckCircle2 className="h-4 w-4 text-emerald-600 mt-0.5 shrink-0" />
            <div>
              <div className="font-semibold text-emerald-800">
                {t(
                  'iam.page.created_success',
                  { name: newlyCreatedUser.name },
                  `User "${newlyCreatedUser.name}" successfully invited!`
                )}
              </div>
              <div className="text-emerald-700 text-[11px] mt-0.5">
                {t(
                  'iam.page.created_details',
                  { role: newlyCreatedUser.role, email: newlyCreatedUser.email },
                  `Role: ${newlyCreatedUser.role} • Email: ${newlyCreatedUser.email}`
                )}
              </div>
            </div>
          </div>
          <Button
            size="sm"
            variant="outline"
            className="text-xs bg-white text-emerald-900 border-emerald-200 hover:bg-emerald-50"
            onClick={() => setNewlyCreatedUser(null)}
          >
            {t('common.dismiss', 'Dismiss')}
          </Button>
        </div>
      )}

      {/* Tab Navigation */}
      <div className="flex items-center gap-2 border-b border-neutral-200 pb-px">
        <button
          type="button"
          onClick={() => setActiveTab('users')}
          className={`flex items-center gap-2 pb-2.5 px-3 text-xs font-medium border-b-2 transition-colors ${
            activeTab === 'users'
              ? 'border-neutral-900 text-neutral-900'
              : 'border-transparent text-neutral-500 hover:text-neutral-700'
          }`}
        >
          <UsersIcon className="h-3.5 w-3.5" />
          <span>{t('iam.page.tab_users', 'Members & Directory')}</span>
        </button>

        <button
          type="button"
          onClick={() => setActiveTab('matrix')}
          className={`flex items-center gap-2 pb-2.5 px-3 text-xs font-medium border-b-2 transition-colors ${
            activeTab === 'matrix'
              ? 'border-neutral-900 text-neutral-900'
              : 'border-transparent text-neutral-500 hover:text-neutral-700'
          }`}
        >
          <Key className="h-3.5 w-3.5" />
          <span>{t('iam.page.tab_matrix', 'Roles & Permissions Matrix')}</span>
        </button>
      </div>

      {/* Tab 1: Users Directory */}
      {activeTab === 'users' && (
        <div className="space-y-6">
          <UserKpiCards
            users={usersData?.items || []}
            totalCount={usersData?.pagination?.totalItems}
            isLoading={usersLoading}
          />

          <UserListTable
            onInspectUser={handleInspectUser}
            onEditUser={handleEditUser}
            onInviteClick={() => setCreateDialogOpen(true)}
          />

          <div id="user-inspector-section">
            <UserInspector
              user={selectedUser}
              onClose={() => setSelectedUser(null)}
            />
          </div>
        </div>
      )}

      {/* Tab 2: Roles & Permissions */}
      {activeTab === 'matrix' && (
        <div className="space-y-6">
          <RoleListTable actorPermissions={actorPermissions} />
          <RolePermissionMatrix />
        </div>
      )}

      {/* Create User Modal */}
      <CreateUserDialog
        open={createDialogOpen}
        onOpenChange={setCreateDialogOpen}
        onSuccess={handleUserCreated}
      />

      {/* Edit User Modal */}
      <EditUserDialog
        user={editingUser}
        open={editDialogOpen}
        onOpenChange={setEditDialogOpen}
        onSuccess={(updated) => {
          if (selectedUser?.id === updated.id) {
            setSelectedUser(updated);
          }
        }}
      />
    </div>
  );
}

export function UserPage() {
  return (
    <UserStoreProvider>
      <UserContent />
    </UserStoreProvider>
  );
}
