import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { tenantKey } from '@/core/tenant/tenant-key';
import { useTenant } from '@/core/tenant/tenant-context';
import {
  assignUserRoles,
  cloneRole,
  createRole,
  createUser,
  deleteRole,
  deleteUser,
  getPermissionsCatalog,
  getPermissionsMeta,
  getRoleById,
  getRoles,
  getRolesMeta,
  getUserById,
  getUsers,
  updateRole,
  updateUser,
  updateUserStatus,
} from '../api';
import type {
  AssignRolesInput,
  CloneRoleInput,
  CreateRoleInput,
  CreateUserInput,
  ListUsersParams,
  UpdateRoleInput,
  UpdateUserInput,
  UserStatus,
} from '../types';

export function useUsers(params: ListUsersParams = {}) {
  const { tenant } = useTenant();

  return useQuery({
    queryKey: tenantKey(tenant.id, ['users', params]),
    queryFn: () => getUsers(params),
    staleTime: 30 * 1000,
  });
}

export function useUser(id?: string) {
  const { tenant } = useTenant();

  return useQuery({
    queryKey: tenantKey(tenant.id, ['user', id]),
    queryFn: () => (id ? getUserById(id) : Promise.reject(new Error('User ID is required'))),
    enabled: Boolean(id),
    staleTime: 60 * 1000,
  });
}

export function useRolesMeta() {
  return useQuery({
    queryKey: ['roles-meta'],
    queryFn: getRolesMeta,
    staleTime: 10 * 60 * 1000,
  });
}

export function usePermissionsMeta() {
  return useQuery({
    queryKey: ['permissions-meta'],
    queryFn: getPermissionsMeta,
    staleTime: 10 * 60 * 1000,
  });
}

export function useCreateUser() {
  const queryClient = useQueryClient();
  const { tenant } = useTenant();

  return useMutation({
    mutationFn: (input: CreateUserInput) => createUser(input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: tenantKey(tenant.id, ['users']) });
    },
  });
}

export function useUpdateUser() {
  const queryClient = useQueryClient();
  const { tenant } = useTenant();

  return useMutation({
    mutationFn: ({ id, input }: { id: string; input: UpdateUserInput }) =>
      updateUser(id, input),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: tenantKey(tenant.id, ['users']) });
      queryClient.invalidateQueries({ queryKey: tenantKey(tenant.id, ['user', variables.id]) });
    },
  });
}

export function useUpdateUserStatus() {
  const queryClient = useQueryClient();
  const { tenant } = useTenant();

  return useMutation({
    mutationFn: ({ id, status }: { id: string; status: UserStatus }) =>
      updateUserStatus(id, status),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: tenantKey(tenant.id, ['users']) });
      queryClient.invalidateQueries({ queryKey: tenantKey(tenant.id, ['user', variables.id]) });
    },
  });
}

export function useDeleteUser() {
  const queryClient = useQueryClient();
  const { tenant } = useTenant();

  return useMutation({
    mutationFn: (id: string) => deleteUser(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: tenantKey(tenant.id, ['users']) });
    },
  });
}

// Dynamic RBAC Hooks
export function useRoles() {
  const { tenant } = useTenant();

  return useQuery({
    queryKey: tenantKey(tenant.id, ['roles']),
    queryFn: getRoles,
    staleTime: 30 * 1000,
  });
}

export function useRole(id?: string) {
  const { tenant } = useTenant();

  return useQuery({
    queryKey: tenantKey(tenant.id, ['role', id]),
    queryFn: () => (id ? getRoleById(id) : Promise.reject(new Error('Role ID is required'))),
    enabled: Boolean(id),
    staleTime: 60 * 1000,
  });
}

export function usePermissionsCatalog() {
  return useQuery({
    queryKey: ['permissions-catalog'],
    queryFn: getPermissionsCatalog,
    staleTime: 10 * 60 * 1000,
  });
}

export function useCreateRole() {
  const queryClient = useQueryClient();
  const { tenant } = useTenant();

  return useMutation({
    mutationFn: (input: CreateRoleInput) => createRole(input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: tenantKey(tenant.id, ['roles']) });
    },
  });
}

export function useUpdateRole() {
  const queryClient = useQueryClient();
  const { tenant } = useTenant();

  return useMutation({
    mutationFn: ({ id, input }: { id: string; input: UpdateRoleInput }) =>
      updateRole(id, input),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: tenantKey(tenant.id, ['roles']) });
      queryClient.invalidateQueries({ queryKey: tenantKey(tenant.id, ['role', variables.id]) });
    },
  });
}

export function useDeleteRole() {
  const queryClient = useQueryClient();
  const { tenant } = useTenant();

  return useMutation({
    mutationFn: (id: string) => deleteRole(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: tenantKey(tenant.id, ['roles']) });
    },
  });
}

export function useCloneRole() {
  const queryClient = useQueryClient();
  const { tenant } = useTenant();

  return useMutation({
    mutationFn: ({ id, input }: { id: string; input: CloneRoleInput }) =>
      cloneRole(id, input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: tenantKey(tenant.id, ['roles']) });
    },
  });
}

export function useAssignUserRoles() {
  const queryClient = useQueryClient();
  const { tenant } = useTenant();

  return useMutation({
    mutationFn: ({ userId, input }: { userId: string; input: AssignRolesInput }) =>
      assignUserRoles(userId, input),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: tenantKey(tenant.id, ['users']) });
      queryClient.invalidateQueries({ queryKey: tenantKey(tenant.id, ['user', variables.userId]) });
      queryClient.invalidateQueries({ queryKey: tenantKey(tenant.id, ['roles']) });
    },
  });
}
