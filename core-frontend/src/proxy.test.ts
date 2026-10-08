import { describe, it, expect } from 'vitest';
import { extractTenantFromHost } from './proxy';

describe('extractTenantFromHost (proxy)', () => {
  it('resolves default for null or empty host', () => {
    expect(extractTenantFromHost(null)).toBe('default');
    expect(extractTenantFromHost('')).toBe('default');
  });

  it('resolves default for localhost and IP addresses', () => {
    expect(extractTenantFromHost('localhost')).toBe('default');
    expect(extractTenantFromHost('localhost:3000')).toBe('default');
    expect(extractTenantFromHost('127.0.0.1:8080')).toBe('default');
  });

  it('resolves subdomains correctly', () => {
    expect(extractTenantFromHost('acme.saas.com')).toBe('acme');
    expect(extractTenantFromHost('org42.saas.com:3000')).toBe('org42');
    expect(extractTenantFromHost('beta.localhost:3000')).toBe('beta');
  });

  it('falls back to default for reserved subdomains', () => {
    expect(extractTenantFromHost('www.saas.com')).toBe('default');
    expect(extractTenantFromHost('api.saas.com')).toBe('default');
    expect(extractTenantFromHost('app.saas.com')).toBe('default');
  });
});
