'use client';

import React, { useState, useMemo } from 'react';
import {
  Search,
  Filter,
  RefreshCw,
  Plus,
  User as UserIcon,
  Shield,
  Key,
  ChevronLeft,
  ChevronRight,
  AlertCircle,
  MoreVertical,
  Edit2,
  Trash2,
  UserX,
  UserCheck,
} from 'lucide-react';
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { useI18n } from '@/core/i18n/i18n-context';
import { useUsers, useUpdateUserStatus, useDeleteUser } from '../hooks';
import type { User, UserRole, UserStatus } from '../types';

interface UserListTableProps {
  onInspectUser?: (user: User) => void;
  onEditUser?: (user: User) => void;
  onInviteClick?: () => void;
}

export function UserListTable({ onInspectUser, onEditUser, onInviteClick }: UserListTableProps) {
  const { t } = useI18n();

  // Search & Filter state
  const [searchInput, setSearchInput] = useState('');
  const [selectedRole, setSelectedRole] = useState<string>('all');
  const [selectedStatus, setSelectedStatus] = useState<string>('all');
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [actionMenuOpenId, setActionMenuOpenId] = useState<string | null>(null);

  // Queries & Mutations
  const {
    data: usersData,
    isLoading,
    isFetching,
    error,
    refetch,
  } = useUsers({
    page,
    pageSize,
    keyword: searchInput,
    role: selectedRole !== 'all' ? selectedRole : undefined,
    status: selectedStatus !== 'all' ? selectedStatus : undefined,
  });

  const updateStatusMutation = useUpdateUserStatus();
  const deleteMutation = useDeleteUser();

  const users = usersData?.items || [];
  const pagination = usersData?.pagination;
  const totalItems = pagination?.totalItems ?? pagination?.total_items ?? 0;
  const totalPages = pagination?.totalPages ?? pagination?.total_pages ?? 1;

  const handleStatusToggle = async (user: User) => {
    const nextStatus: UserStatus = user.status === 'active' ? 'suspended' : 'active';
    try {
      await updateStatusMutation.mutateAsync({ id: user.id, status: nextStatus });
      setActionMenuOpenId(null);
    } catch {
      // Handled by toast/mutation error state
    }
  };

  const handleDelete = async (user: User) => {
    if (
      window.confirm(
        t(
          'iam.list.delete_confirm',
          `Are you sure you want to remove user "${user.name}" (${user.email})?`
        )
      )
    ) {
      try {
        await deleteMutation.mutateAsync(user.id);
        setActionMenuOpenId(null);
      } catch {
        // Handled by query client
      }
    }
  };

  return (
    <Card className="shadow-none border border-neutral-200">
      <CardHeader className="pb-4">
        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
          <div>
            <CardTitle className="text-sm font-semibold text-neutral-900 flex items-center gap-2">
              <UserIcon className="h-4 w-4 text-neutral-700" />
              <span>{t('iam.list.title', 'Organization Members')}</span>
            </CardTitle>
            <CardDescription className="text-xs text-neutral-500 mt-1">
              {t(
                'iam.list.subtitle',
                'Directory of all registered accounts, assigned roles, and permission levels.'
              )}
            </CardDescription>
          </div>

          <div className="flex items-center gap-2">
            <Button
              variant="outline"
              size="sm"
              onClick={() => refetch()}
              disabled={isFetching}
              className="text-xs gap-1.5 h-8 text-neutral-600"
            >
              <RefreshCw className={`h-3.5 w-3.5 ${isFetching ? 'animate-spin' : ''}`} />
              <span>Refresh</span>
            </Button>

            {onInviteClick && (
              <Button size="sm" onClick={onInviteClick} className="text-xs gap-1.5 h-8">
                <Plus className="h-3.5 w-3.5" />
                <span>{t('user.page.create_btn', t('iam.page.invite_btn', 'Create User'))}</span>
              </Button>
            )}
          </div>
        </div>

        {/* Search & Filter Bar */}
        <div className="grid grid-cols-1 sm:grid-cols-12 gap-2.5 pt-3">
          {/* Keyword Search */}
          <div className="sm:col-span-6 relative">
            <Search className="absolute left-2.5 top-2.5 h-3.5 w-3.5 text-neutral-400" />
            <Input
              type="text"
              placeholder={t('iam.list.search_placeholder', 'Search users by name or email...')}
              value={searchInput}
              onChange={(e) => {
                setSearchInput(e.target.value);
                setPage(1);
              }}
              className="pl-8 text-xs h-9 bg-neutral-50/50"
            />
          </div>

          {/* Role Filter */}
          <div className="sm:col-span-3">
            <select
              value={selectedRole}
              onChange={(e) => {
                setSelectedRole(e.target.value);
                setPage(1);
              }}
              className="w-full h-9 rounded-md border border-neutral-200 bg-white px-3 py-1 text-xs text-neutral-800 focus:outline-none focus:ring-1 focus:ring-neutral-900"
              aria-label={t('iam.list.filter_role', 'Filter by Role')}
            >
              <option value="all">{t('iam.list.all_roles', 'All Roles')}</option>
              <option value="owner">{t('iam.roles.owner', 'Owner')}</option>
              <option value="admin">{t('iam.roles.admin', 'Administrator')}</option>
              <option value="member">{t('iam.roles.member', 'Member')}</option>
              <option value="viewer">{t('iam.roles.viewer', 'Viewer')}</option>
            </select>
          </div>

          {/* Status Filter */}
          <div className="sm:col-span-3">
            <select
              value={selectedStatus}
              onChange={(e) => {
                setSelectedStatus(e.target.value);
                setPage(1);
              }}
              className="w-full h-9 rounded-md border border-neutral-200 bg-white px-3 py-1 text-xs text-neutral-800 focus:outline-none focus:ring-1 focus:ring-neutral-900"
              aria-label={t('iam.list.filter_status', 'Filter by Status')}
            >
              <option value="all">{t('iam.list.all_statuses', 'All Statuses')}</option>
              <option value="active">{t('iam.statuses.active', 'Active')}</option>
              <option value="suspended">{t('iam.statuses.suspended', 'Suspended')}</option>
              <option value="invited">{t('iam.statuses.invited', 'Invited')}</option>
            </select>
          </div>
        </div>
      </CardHeader>

      <CardContent className="p-0">
        {/* Table Content */}
        <div className="overflow-x-auto">
          <table className="w-full text-left border-collapse text-xs">
            <thead>
              <tr className="border-y border-neutral-200 bg-neutral-50/70 text-neutral-600 font-semibold uppercase tracking-wider text-[11px]">
                <th className="py-2.5 px-4">{t('iam.list.col_user', 'User')}</th>
                <th className="py-2.5 px-4">{t('iam.list.col_role', 'Role')}</th>
                <th className="py-2.5 px-4">{t('iam.list.col_acl', 'Access Control (ACL)')}</th>
                <th className="py-2.5 px-4">{t('iam.list.col_status', 'Status')}</th>
                <th className="py-2.5 px-4">{t('iam.list.col_created', 'Joined Date')}</th>
                <th className="py-2.5 px-4 text-right">{t('iam.list.col_actions', 'Actions')}</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-neutral-100">
              {isLoading ? (
                <tr>
                  <td colSpan={6} className="py-8 text-center text-neutral-500">
                    <div className="flex items-center justify-center gap-2">
                      <RefreshCw className="h-4 w-4 animate-spin text-neutral-400" />
                      <span>{t('common.loading', 'Loading members...')}</span>
                    </div>
                  </td>
                </tr>
              ) : error ? (
                <tr>
                  <td colSpan={6} className="py-8 text-center text-red-600">
                    <div className="flex items-center justify-center gap-1.5 text-xs">
                      <AlertCircle className="h-4 w-4" />
                      <span>Failed to load users: {(error as Error).message}</span>
                    </div>
                  </td>
                </tr>
              ) : users.length === 0 ? (
                <tr>
                  <td colSpan={6} className="py-12 text-center">
                    <div className="max-w-sm mx-auto space-y-3">
                      <UserIcon className="h-8 w-8 text-neutral-300 mx-auto" />
                      <div className="font-semibold text-neutral-800 text-sm">
                        {searchInput || selectedRole !== 'all' || selectedStatus !== 'all'
                          ? t('iam.list.empty_search', 'No matching members found')
                          : t('iam.list.empty_title', 'No members registered yet')}
                      </div>
                      <p className="text-xs text-neutral-500">
                        {searchInput || selectedRole !== 'all' || selectedStatus !== 'all'
                          ? 'Try adjusting your search keyword or active filters.'
                          : t('iam.list.empty_initial', 'Click "Create User" to onboard your first team member.')}
                      </p>
                      {onInviteClick && !searchInput && selectedRole === 'all' && selectedStatus === 'all' && (
                        <div className="pt-2">
                          <Button size="sm" onClick={onInviteClick} className="text-xs gap-1.5 h-8">
                            <Plus className="h-3.5 w-3.5" />
                            <span>{t('user.page.create_btn', 'Create User')}</span>
                          </Button>
                        </div>
                      )}
                    </div>
                  </td>
                </tr>
              ) : (
                users.map((u) => {
                  const isOwner = u.role === 'owner';
                  const customPermCount = u.permissions?.length || 0;

                  return (
                    <tr
                      key={u.id}
                      className="hover:bg-neutral-50/60 transition-colors group cursor-pointer"
                      onClick={() => onInspectUser?.(u)}
                    >
                      {/* User Info */}
                      <td className="py-3 px-4">
                        <div className="flex items-center gap-3">
                          <div className="h-8 w-8 rounded-full border border-neutral-300 bg-neutral-100 flex items-center justify-center text-neutral-700 font-semibold text-xs shrink-0">
                            {u.avatar_url ? (
                              // eslint-disable-next-line @next/next/no-img-element
                              <img
                                src={u.avatar_url}
                                alt={u.name}
                                className="h-full w-full rounded-full object-cover"
                              />
                            ) : (
                              u.name.slice(0, 2).toUpperCase()
                            )}
                          </div>
                          <div className="min-w-0">
                            <div className="font-medium text-neutral-900 truncate">{u.name}</div>
                            <div className="text-[11px] font-mono text-neutral-500 truncate">
                              {u.email}
                            </div>
                          </div>
                        </div>
                      </td>

                      {/* Role */}
                      <td className="py-3 px-4 whitespace-nowrap">
                        <div className="inline-flex items-center gap-1.5 px-2 py-0.5 rounded border border-neutral-200 bg-neutral-50 font-medium text-neutral-800 capitalize">
                          <Shield className="h-3 w-3 text-neutral-500" />
                          <span>{t(`iam.roles.${u.role}`, u.role)}</span>
                        </div>
                      </td>

                      {/* ACL Overrides */}
                      <td className="py-3 px-4">
                        {isOwner ? (
                          <span className="text-[11px] font-mono text-indigo-700 bg-indigo-50 px-2 py-0.5 rounded border border-indigo-200">
                            * (Universal)
                          </span>
                        ) : customPermCount > 0 ? (
                          <div className="flex items-center gap-1.5">
                            <span className="text-[11px] font-mono font-medium text-amber-800 bg-amber-50 px-2 py-0.5 rounded border border-amber-200 flex items-center gap-1">
                              <Key className="h-3 w-3 text-amber-600" />
                              <span>+{customPermCount} custom ACL</span>
                            </span>
                          </div>
                        ) : (
                          <span className="text-[11px] text-neutral-400">
                            {t('iam.list.default_role_badge', 'Standard role permissions')}
                          </span>
                        )}
                      </td>

                      {/* Status */}
                      <td className="py-3 px-4 whitespace-nowrap">
                        <Badge
                          variant={
                            u.status === 'active'
                              ? 'success'
                              : u.status === 'invited'
                              ? 'outline'
                              : 'secondary'
                          }
                          className="text-[10px] uppercase font-mono tracking-wider"
                        >
                          {t(`iam.statuses.${u.status}`, u.status)}
                        </Badge>
                      </td>

                      {/* Joined Date */}
                      <td className="py-3 px-4 whitespace-nowrap text-neutral-500 font-mono text-[11px]">
                        {u.create_at ? new Date(u.create_at).toLocaleDateString() : '—'}
                      </td>

                      {/* Actions */}
                      <td
                        className="py-3 px-4 text-right whitespace-nowrap"
                        onClick={(e) => e.stopPropagation()}
                      >
                        <div className="flex items-center justify-end gap-1">
                          {/* Quick Inspect Button */}
                          <Button
                            variant="ghost"
                            size="sm"
                            onClick={() => onInspectUser?.(u)}
                            className="h-7 px-2 text-[11px] text-neutral-600 hover:text-neutral-900"
                          >
                            {t('iam.list.inspect', 'Inspect')}
                          </Button>

                          {/* Quick Edit */}
                          {onEditUser && (
                            <Button
                              variant="ghost"
                              size="sm"
                              onClick={() => onEditUser(u)}
                              className="h-7 w-7 p-0 text-neutral-500 hover:text-neutral-900"
                              title={t('iam.list.edit', 'Edit Member')}
                            >
                              <Edit2 className="h-3.5 w-3.5" />
                            </Button>
                          )}

                          {/* Actions Dropdown Toggle */}
                          <div className="relative inline-block text-left">
                            <Button
                              variant="ghost"
                              size="sm"
                              onClick={() =>
                                setActionMenuOpenId(actionMenuOpenId === u.id ? null : u.id)
                              }
                              className="h-7 w-7 p-0 text-neutral-500 hover:text-neutral-900"
                            >
                              <MoreVertical className="h-3.5 w-3.5" />
                            </Button>

                            {actionMenuOpenId === u.id && (
                              <div className="absolute right-0 mt-1 w-44 rounded-md bg-white shadow-lg border border-neutral-200 py-1 z-20 text-xs">
                                <button
                                  type="button"
                                  onClick={() => handleStatusToggle(u)}
                                  className="w-full text-left px-3 py-1.5 flex items-center gap-2 hover:bg-neutral-50 text-neutral-700"
                                >
                                  {u.status === 'active' ? (
                                    <>
                                      <UserX className="h-3.5 w-3.5 text-amber-600" />
                                      <span>Suspend Member</span>
                                    </>
                                  ) : (
                                    <>
                                      <UserCheck className="h-3.5 w-3.5 text-emerald-600" />
                                      <span>Activate Member</span>
                                    </>
                                  )}
                                </button>

                                {!isOwner && (
                                  <button
                                    type="button"
                                    onClick={() => handleDelete(u)}
                                    className="w-full text-left px-3 py-1.5 flex items-center gap-2 hover:bg-red-50 text-red-600 border-t border-neutral-100"
                                  >
                                    <Trash2 className="h-3.5 w-3.5" />
                                    <span>{t('iam.list.delete', 'Remove Member')}</span>
                                  </button>
                                )}
                              </div>
                            )}
                          </div>
                        </div>
                      </td>
                    </tr>
                  );
                })
              )}
            </tbody>
          </table>
        </div>

        {/* Pagination Bar */}
        {totalItems > 0 && (
          <div className="flex flex-col sm:flex-row items-center justify-between gap-3 px-4 py-3 border-t border-neutral-200 bg-neutral-50/50 text-xs text-neutral-500">
            <div>
              {t(
                'iam.list.showing_results',
                {
                  from: (page - 1) * pageSize + 1,
                  to: Math.min(page * pageSize, totalItems),
                  total: totalItems,
                },
                `Showing ${(page - 1) * pageSize + 1} - ${Math.min(
                  page * pageSize,
                  totalItems
                )} of ${totalItems} members`
              )}
            </div>

            <div className="flex items-center gap-2">
              <div className="flex items-center gap-1 mr-2">
                <span className="text-[11px] text-neutral-500">
                  {t('iam.list.page_size', 'Per page')}:
                </span>
                <select
                  value={pageSize}
                  onChange={(e) => {
                    setPageSize(Number(e.target.value));
                    setPage(1);
                  }}
                  className="rounded border border-neutral-200 bg-white px-2 py-0.5 text-xs text-neutral-700"
                >
                  <option value={10}>10</option>
                  <option value={25}>25</option>
                  <option value={50}>50</option>
                </select>
              </div>

              <div className="flex items-center gap-1">
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => setPage((p) => Math.max(1, p - 1))}
                  disabled={page <= 1}
                  className="h-7 px-2 text-xs"
                >
                  <ChevronLeft className="h-3.5 w-3.5" />
                  <span className="sr-only">Previous Page</span>
                </Button>

                <span className="text-xs text-neutral-600 px-2 font-mono">
                  {page} / {totalPages}
                </span>

                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
                  disabled={page >= totalPages}
                  className="h-7 px-2 text-xs"
                >
                  <ChevronRight className="h-3.5 w-3.5" />
                  <span className="sr-only">Next Page</span>
                </Button>
              </div>
            </div>
          </div>
        )}
      </CardContent>
    </Card>
  );
}
