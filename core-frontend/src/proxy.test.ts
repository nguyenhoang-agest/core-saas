import { describe, it, expect } from 'vitest';
import { NextRequest } from 'next/server';
import { extractTenantFromHost, proxy } from './proxy';

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

describe('proxy request guard', () => {
  it('allows unauthenticated requests to public overview root /', () => {
    const req = new NextRequest('http://localhost:3000/');
    const res = proxy(req);
    expect(res.status).toBe(200);
    expect(res.headers.get('location')).toBeNull();
  });

  it('redirects unauthenticated requests to protected route /users to /login', () => {
    const req = new NextRequest('http://localhost:3000/users');
    const res = proxy(req);
    expect(res.status).toBe(307);
    const location = res.headers.get('location');
    expect(location).toContain('/login?callbackUrl=%2Fusers');
  });

  it('allows authenticated requests to protected route /users with session token', () => {
    const req = new NextRequest('http://localhost:3000/users', {
      headers: {
        cookie: 'next-auth.session-token=mock-token',
      },
    });
    const res = proxy(req);
    expect(res.status).toBe(200);
    expect(res.headers.get('location')).toBeNull();
  });
});
