import { NextResponse, type NextRequest } from 'next/server';
import { isProtectedRoute } from './core/auth/routes';

const RESERVED_SUBDOMAINS = new Set(['www', 'api', 'app', 'admin', 'auth', 'mail', 'static']);

/**
 * Extracts subdomain from host.
 * Examples:
 * - 'acme.saas.local:3000' -> 'acme'
 * - 'tenant1.example.com' -> 'tenant1'
 * - 'localhost:3000' -> 'default'
 * - '127.0.0.1' -> 'default'
 */
export function extractTenantFromHost(host: string | null): string {
  if (!host) return 'default';

  // Strip port
  const hostname = host.split(':')[0].toLowerCase();

  // If localhost or IP, default
  if (hostname === 'localhost' || hostname === '127.0.0.1' || /^(\d{1,3}\.){3}\d{1,3}$/.test(hostname)) {
    return 'default';
  }

  const parts = hostname.split('.');

  // Need at least 2 parts for a domain, 3 parts for a subdomain (e.g., tenant.domain.com or tenant.domain.local)
  if (parts.length >= 2) {
    if (parts.length === 2 && parts[1] === 'localhost') {
      const sub = parts[0];
      return RESERVED_SUBDOMAINS.has(sub) ? 'default' : sub;
    }

    if (parts.length >= 3) {
      const sub = parts[0];
      return RESERVED_SUBDOMAINS.has(sub) ? 'default' : sub;
    }
  }

  return 'default';
}

/**
 * Request interception using Next.js 16 proxy convention.
 */
export function proxy(request: NextRequest) {
  const pathname = request.nextUrl.pathname;

  // Protected management routes require an active NextAuth session
  if (isProtectedRoute(pathname)) {
    const sessionToken =
      request.cookies.get('next-auth.session-token')?.value ||
      request.cookies.get('__Secure-next-auth.session-token')?.value;

    if (!sessionToken) {
      const loginUrl = new URL('/login', request.url);
      loginUrl.searchParams.set('callbackUrl', pathname);
      return NextResponse.redirect(loginUrl);
    }
  }

  const host = request.headers.get('host');

  // Allow explicit header, query parameter, or cookie override
  const explicitHeader = request.headers.get('x-tenant-id');
  const queryTenant = request.nextUrl.searchParams.get('tenant');
  const cookieTenant = request.cookies.get('x-tenant-id')?.value;
  const hostTenant = extractTenantFromHost(host);

  let tenantId = 'default';
  if (explicitHeader && explicitHeader.trim() !== '') {
    tenantId = explicitHeader.trim();
  } else if (queryTenant && queryTenant.trim() !== '') {
    tenantId = queryTenant.trim();
  } else if (hostTenant && hostTenant !== 'default') {
    tenantId = hostTenant;
  } else if (cookieTenant && cookieTenant.trim() !== '' && cookieTenant !== 'default') {
    tenantId = cookieTenant.trim();
  } else {
    tenantId = hostTenant || 'default';
  }

  const requestHeaders = new Headers(request.headers);
  requestHeaders.set('x-tenant-id', tenantId);

  const response = NextResponse.next({
    request: {
      headers: requestHeaders,
    },
  });

  // Set cookie so client-side Axios & components can access current tenant synchronously
  response.cookies.set('x-tenant-id', tenantId, {
    path: '/',
    sameSite: 'lax',
  });
  response.headers.set('x-tenant-id', tenantId);

  return response;
}

export const config = {
  matcher: [
    /*
     * Match all request paths except for static assets:
     * - api/auth (NextAuth endpoints)
     * - _next/static (static files)
     * - _next/image (image optimization files)
     * - favicon.ico, sitemap.xml, robots.txt
     */
    '/((?!api/auth|_next/static|_next/image|favicon.ico|sitemap.xml|robots.txt).*)',
  ],
};
