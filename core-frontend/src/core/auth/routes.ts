/**
 * Route governance and access rules for Core-SaaS.
 * Provides a single source of truth for public, auth, and protected management routes.
 */

export const AUTH_ROUTES = ['/login', '/register', '/auth'] as const;

export const INTERNAL_PREFIXES = [
  '/_next',
  '/api',
  '/favicon.ico',
  '/sitemap.xml',
  '/robots.txt',
] as const;

/**
 * Normalizes a pathname by stripping query params and trailing slashes.
 */
export function normalizePathname(pathname: string): string {
  if (!pathname) return '/';
  const clean = pathname.split('?')[0].split('#')[0];
  if (clean !== '/' && clean.endsWith('/')) {
    return clean.slice(0, -1);
  }
  return clean || '/';
}

/**
 * Returns true if the pathname is an internal Next.js or static resource route.
 */
export function isInternalRoute(pathname: string): boolean {
  const normalized = normalizePathname(pathname);
  return INTERNAL_PREFIXES.some(
    (prefix) => normalized === prefix || normalized.startsWith(`${prefix}/`)
  );
}

/**
 * Returns true if the pathname is an authentication page (login, register, auth error).
 */
export function isAuthRoute(pathname: string): boolean {
  const normalized = normalizePathname(pathname);
  return AUTH_ROUTES.some(
    (route) => normalized === route || normalized.startsWith(`${route}/`)
  );
}

/**
 * Returns true if the pathname can be accessed publicly without authentication.
 * By default, only the Overview root ('/') and authentication routes are public.
 */
export function isPublicRoute(pathname: string): boolean {
  const normalized = normalizePathname(pathname);
  if (normalized === '/') return true;
  if (isInternalRoute(normalized)) return true;
  return isAuthRoute(normalized);
}

/**
 * Returns true if the pathname is a protected management page requiring an active session.
 * E.g. /users, /tenants, /settings, /billing, /example.
 */
export function isProtectedRoute(pathname: string): boolean {
  return !isPublicRoute(pathname);
}
