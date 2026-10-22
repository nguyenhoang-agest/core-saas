import { describe, it, expect } from 'vitest';
import {
  isPublicRoute,
  isProtectedRoute,
  isAuthRoute,
  isInternalRoute,
  normalizePathname,
} from './routes';

describe('Route Governance (routes.ts)', () => {
  describe('normalizePathname', () => {
    it('normalizes root path', () => {
      expect(normalizePathname('')).toBe('/');
      expect(normalizePathname('/')).toBe('/');
      expect(normalizePathname('/?query=1')).toBe('/');
    });

    it('strips trailing slashes and query strings', () => {
      expect(normalizePathname('/users/')).toBe('/users');
      expect(normalizePathname('/users?tab=roles')).toBe('/users');
      expect(normalizePathname('/settings/profile/')).toBe('/settings/profile');
    });
  });

  describe('isPublicRoute & isProtectedRoute', () => {
    it('treats root overview page as public', () => {
      expect(isPublicRoute('/')).toBe(true);
      expect(isProtectedRoute('/')).toBe(false);
    });

    it('treats login, register, and auth error pages as public', () => {
      expect(isPublicRoute('/login')).toBe(true);
      expect(isPublicRoute('/register')).toBe(true);
      expect(isPublicRoute('/auth/error')).toBe(true);
      expect(isProtectedRoute('/login')).toBe(false);
      expect(isProtectedRoute('/register')).toBe(false);
      expect(isProtectedRoute('/auth/error')).toBe(false);
    });

    it('treats management routes as protected', () => {
      const protectedPages = [
        '/users',
        '/users/invite',
        '/tenants',
        '/tenants/new',
        '/settings',
        '/settings/general',
        '/billing',
        '/billing/invoices',
        '/example',
      ];

      for (const path of protectedPages) {
        expect(isPublicRoute(path)).toBe(false);
        expect(isProtectedRoute(path)).toBe(true);
      }
    });

    it('treats internal, API, and static resources as public', () => {
      expect(isInternalRoute('/_next/static/chunk.js')).toBe(true);
      expect(isInternalRoute('/api/auth/session')).toBe(true);
      expect(isInternalRoute('/api/register')).toBe(true);
      expect(isInternalRoute('/favicon.ico')).toBe(true);
      expect(isPublicRoute('/_next/static/chunk.js')).toBe(true);
      expect(isPublicRoute('/api/register')).toBe(true);
    });
  });

  describe('isAuthRoute', () => {
    it('identifies authentication routes correctly', () => {
      expect(isAuthRoute('/login')).toBe(true);
      expect(isAuthRoute('/register')).toBe(true);
      expect(isAuthRoute('/auth/signin')).toBe(true);
      expect(isAuthRoute('/')).toBe(false);
      expect(isAuthRoute('/users')).toBe(false);
    });
  });
});
