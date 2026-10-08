/**
 * Helper to prefix React Query keys with the active tenant ID.
 * Prevents cache collision across tenants in multi-tenant environments.
 *
 * @example
 * const queryKey = tenantKey('tenant-123', ['users', userId]);
 * // Returns: ['tenant', 'tenant-123', 'users', userId]
 */
export function tenantKey(
  tenantId: string | null | undefined,
  queryKey: readonly unknown[]
): readonly unknown[] {
  const safeTenant = tenantId || 'default';
  return ['tenant', safeTenant, ...queryKey];
}
