import { describe, it, expect } from 'vitest';
import { tenantKey } from './tenant-key';

describe('tenantKey helper', () => {
  it('prefixes query keys with the provided tenant ID', () => {
    const key = tenantKey('org-42', ['billing', 'invoices']);
    expect(key).toEqual(['tenant', 'org-42', 'billing', 'invoices']);
  });

  it('falls back to "default" when tenantId is null or undefined or empty', () => {
    expect(tenantKey(undefined, ['users'])).toEqual(['tenant', 'default', 'users']);
    expect(tenantKey(null, ['settings'])).toEqual(['tenant', 'default', 'settings']);
    expect(tenantKey('', ['settings'])).toEqual(['tenant', 'default', 'settings']);
  });
});
