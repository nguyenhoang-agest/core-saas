'use client';

import React from 'react';
import Link from 'next/link';
import { usePathname, useRouter } from 'next/navigation';
import { useSession, signOut } from 'next-auth/react';
import {
  Building2,
  LayoutDashboard,
  Settings,
  ShieldAlert,
  Layers,
  ChevronRight,
  ExternalLink,
  Users as UsersIcon,
  LogOut,
  LogIn,
} from 'lucide-react';
import { cn } from '@/lib/utils';
import { useTenant } from '@/core/tenant/tenant-context';
import { useI18n } from '@/core/i18n/i18n-context';
import { isProtectedRoute, isAuthRoute } from '@/core/auth/routes';
import { appConfig } from '@/../config/app.config';
import { moduleRegistry } from '@/core/modules/registry';
import { LocaleSwitcher } from './locale-switcher';
import { TenantSwitcher } from './tenant-switcher';

interface NavItem {
  titleKey: string;
  fallback: string;
  href: string;
  icon: React.ComponentType<{ className?: string }>;
  badge?: string;
  hostOnly?: boolean;
  authRequired?: boolean;
}

const navItems: NavItem[] = [
  {
    titleKey: 'nav.overview',
    fallback: 'Overview',
    href: '/',
    icon: LayoutDashboard,
    authRequired: false,
  },
  {
    titleKey: 'nav.tenants',
    fallback: 'Tenants',
    href: '/tenants',
    icon: Building2,
    badge: 'Platform',
    hostOnly: true,
    authRequired: true,
  },
  {
    titleKey: 'nav.users',
    fallback: 'Users',
    href: '/users',
    icon: UsersIcon,
    badge: 'RBAC',
    authRequired: true,
  },
  {
    titleKey: 'nav.settings',
    fallback: 'Settings',
    href: '/settings',
    icon: Settings,
    authRequired: true,
  },
];

export function AdminShell({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();
  const router = useRouter();
  const { data: session, status } = useSession();
  const { tenant } = useTenant();
  const { t } = useI18n();

  const isAuthenticated = status === 'authenticated';
  const isHost = !!(session?.user as any)?.isHost;
  const isAuth = isAuthRoute(pathname || '');
  const isProtected = isProtectedRoute(pathname || '');

  React.useEffect(() => {
    if (isProtected && status === 'unauthenticated') {
      const redirectUrl = `/login?callbackUrl=${encodeURIComponent(pathname || '/')}`;
      router.replace(redirectUrl);
    }
  }, [isProtected, status, pathname, router]);

  const visibleNavItems = navItems.filter((item) => {
    if (item.authRequired !== false && !isAuthenticated) return false;
    if (item.hostOnly && !isHost) return false;
    return true;
  });

  if (isAuth) {
    return <main className="min-h-screen bg-neutral-50">{children}</main>;
  }

  const getBreadcrumbs = () => {
    const rootItem = { label: t('shell.admin', 'Admin'), href: '/' };

    if (!pathname || pathname === '/') {
      return [rootItem, { label: t('nav.dashboard', 'Dashboard'), href: '/' }];
    }

    const segments = pathname.split('/').filter(Boolean);
    const items = [rootItem];

    let accumulatedPath = '';
    for (const segment of segments) {
      accumulatedPath += `/${segment}`;

      // 1. Try dictionary lookup for nav.<segment>
      let label = t(`nav.${segment}`, '');

      // 2. Try registered module lookup
      if (!label || label === `nav.${segment}`) {
        const mod = moduleRegistry.get(segment);
        if (mod) {
          label = mod.name;
        } else {
          // 3. Fallback: convert kebab-case or slug to title case
          label = segment.replace(/[-_]+/g, ' ').replace(/\b\w/g, (char) => char.toUpperCase());
        }
      }

      items.push({
        label,
        href: accumulatedPath,
      });
    }

    return items;
  };

  return (
    <div className="flex min-h-screen bg-neutral-50/50 text-neutral-900 font-sans antialiased">
      {/* Sidebar - Outlined Minimalist */}
      <aside className="w-64 border-r border-neutral-200 bg-white flex flex-col justify-between shrink-0">
        <div>
          {/* Logo & Brand Header */}
          <div className="h-14 flex items-center px-5 border-b border-neutral-200 gap-3">
            <div className="h-7 w-7 rounded-md border border-neutral-300 bg-neutral-900 flex items-center justify-center text-white font-bold text-xs">
              CS
            </div>
            <div className="flex flex-col">
              <span className="text-xs font-semibold tracking-tight text-neutral-900 leading-tight">
                Core SaaS
              </span>
              <span className="text-[10px] text-neutral-500 tracking-wider uppercase font-medium">
                {t('shell.admin_console', 'Admin Console')}
              </span>
            </div>
          </div>

          {/* Navigation Items */}
          <div className="px-3 py-4 space-y-1">
            <div className="px-2 pb-2 text-[10px] font-semibold tracking-wider text-neutral-600 uppercase">
              {t('shell.platform_modules', 'Platform Modules')}
            </div>
            {visibleNavItems.map((item) => {
              const Icon = item.icon;
              const isActive =
                item.href === '/' ? pathname === '/' : pathname.startsWith(item.href);

              return (
                <Link
                  key={item.href}
                  href={item.href}
                  className={cn(
                    'flex items-center justify-between px-3 py-2 rounded-md text-xs font-medium transition-colors border',
                    isActive
                      ? 'bg-neutral-100 text-neutral-900 border-neutral-300'
                      : 'text-neutral-600 hover:text-neutral-900 hover:bg-neutral-50 border-transparent'
                  )}
                >
                  <div className="flex items-center gap-2.5">
                    <Icon
                      className={cn(
                        'h-4 w-4 stroke-[1.75]',
                        isActive ? 'text-neutral-900' : 'text-neutral-600'
                      )}
                    />
                    <span>{t(item.titleKey, item.fallback)}</span>
                  </div>
                  {item.badge && (
                    <span
                      className={cn(
                        'text-[10px] px-1.5 py-0.2 rounded border uppercase font-medium',
                        isActive
                          ? 'bg-white text-neutral-900 border-neutral-300'
                          : 'bg-neutral-100 text-neutral-700 border-neutral-200'
                      )}
                    >
                      {item.badge}
                    </span>
                  )}
                </Link>
              );
            })}
          </div>
        </div>

        {/* Sidebar Footer: Current Tenant Context or Guest Mode */}
        {isAuthenticated ? (
          <div className="p-3 border-t border-neutral-200 bg-neutral-50/50">
            <div className="rounded-md border border-neutral-200 bg-white p-2.5 text-xs">
              <div className="flex items-center justify-between text-neutral-600 mb-1">
                <span className="text-[10px] uppercase font-semibold tracking-wider">
                  {t('shell.tenant_context', 'Tenant Context')}
                </span>
                <span className="flex h-1.5 w-1.5 rounded-full bg-emerald-500 ring-2 ring-emerald-100" />
              </div>
              <div className="font-medium text-neutral-900 truncate">{tenant.name}</div>
              <div className="text-[10px] font-mono text-neutral-600 truncate mt-0.5">
                id: {tenant.id}
              </div>
            </div>
          </div>
        ) : (
          <div className="p-3 border-t border-neutral-200 bg-neutral-50/50">
            <div className="rounded-md border border-neutral-200 bg-white p-2.5 text-xs text-center space-y-2">
              <span className="text-[11px] text-neutral-500 block">
                {t('shell.guest_mode', 'Browsing as Guest')}
              </span>
              <Link
                href="/login"
                className="flex items-center justify-center gap-1.5 w-full py-1.5 rounded bg-neutral-900 text-white text-xs font-medium hover:bg-neutral-800 transition-colors"
              >
                <LogIn className="h-3.5 w-3.5" />
                <span>{t('shell.sign_in', 'Sign In')}</span>
              </Link>
            </div>
          </div>
        )}
      </aside>

      {/* Main Content Area */}
      <div className="flex-1 flex flex-col min-w-0">
        {/* Top Header */}
        <header className="h-14 border-b border-neutral-200 bg-white/80 backdrop-blur-sm sticky top-0 z-30 flex items-center justify-between px-6 lg:px-8">
          {/* Breadcrumbs */}
          <nav aria-label="Breadcrumb" className="flex items-center">
            <ol className="flex items-center space-x-1.5 text-xs text-neutral-500">
              {getBreadcrumbs().map((item, index, arr) => {
                const isLast = index === arr.length - 1;
                return (
                  <li key={item.href + index} className="flex items-center space-x-1.5">
                    {index > 0 && (
                      <ChevronRight
                        className="h-3 w-3 text-neutral-400 stroke-[1.5] shrink-0"
                        aria-hidden="true"
                      />
                    )}
                    {isLast ? (
                      <span className="text-neutral-900 font-medium capitalize" aria-current="page">
                        {item.label}
                      </span>
                    ) : (
                      <Link
                        href={item.href}
                        className="text-neutral-500 hover:text-neutral-900 transition-colors capitalize"
                      >
                        {item.label}
                      </Link>
                    )}
                  </li>
                );
              })}
            </ol>
          </nav>

          {/* Right Header Badges */}
          <div className="flex items-center space-x-3 text-xs">
            {isHost ? (
              <TenantSwitcher />
            ) : isAuthenticated ? (
              <div
                className="flex items-center gap-2 px-3 py-1.5 rounded-md border border-neutral-200 bg-white text-neutral-800 text-xs font-medium max-w-[200px]"
                title={tenant.name || tenant.id}
              >
                <Building2 className="h-3.5 w-3.5 text-neutral-600 shrink-0" />
                <span className="truncate">{tenant.name || tenant.id}</span>
              </div>
            ) : null}
            <LocaleSwitcher />
            <div className="hidden sm:flex items-center gap-1.5 px-2.5 py-1 rounded-md border border-neutral-200 bg-neutral-50 text-neutral-600 text-[11px]">
              <Layers className="h-3.5 w-3.5 text-neutral-700 stroke-[1.5]" />
              <span>{t('shell.multi_tenant_arch', 'Multi-Tenant Architecture')}</span>
            </div>
            <a
              href="https://github.com/nguyenhoang-agest/core-saas"
              target="_blank"
              rel="noreferrer"
              className="flex items-center gap-1 text-neutral-700 hover:text-neutral-900 text-xs transition-colors"
            >
              <span>{t('shell.docs', 'Docs')}</span>
              <ExternalLink className="h-3 w-3 stroke-[1.5]" />
            </a>

            {/* User Session Profile & Sign Out */}
            {status === 'authenticated' && session?.user ? (
              <div className="flex items-center gap-2 pl-2 border-l border-neutral-200">
                <div className="flex items-center gap-2">
                  <div className="h-7 w-7 rounded-full bg-neutral-900 text-white flex items-center justify-center font-semibold text-[11px] shadow-sm uppercase">
                    {(session.user.name || session.user.email || 'U')[0]}
                  </div>
                  <div className="hidden lg:flex flex-col text-left">
                    <span className="text-xs font-medium text-neutral-900 leading-tight">
                      {session.user.name || session.user.email}
                    </span>
                    <span className="text-[10px] text-neutral-500 font-mono leading-tight">
                      {session.user.email}
                    </span>
                  </div>
                </div>
                <button
                  type="button"
                  onClick={() => {
                    if (typeof document !== 'undefined') {
                      document.cookie = 'x-tenant-id=; path=/; max-age=0; SameSite=Lax';
                      try {
                        localStorage.removeItem('core_saas_tenant_id');
                        localStorage.removeItem('core_saas_tenant_name');
                      } catch {
                        // ignore
                      }
                    }
                    signOut({ callbackUrl: '/login' });
                  }}
                  title={t('shell.sign_out', 'Sign Out')}
                  className="p-1.5 rounded-md text-neutral-500 hover:text-red-600 hover:bg-neutral-100 transition-colors"
                >
                  <LogOut className="h-4 w-4" />
                </button>
              </div>
            ) : status === 'unauthenticated' ? (
              <Link
                href="/login"
                className="flex items-center gap-1.5 px-2.5 py-1 rounded-md bg-neutral-900 text-white text-xs font-medium hover:bg-neutral-800 transition-colors"
              >
                <LogIn className="h-3.5 w-3.5" />
                <span>{t('shell.sign_in', 'Sign In')}</span>
              </Link>
            ) : null}
          </div>
        </header>

        {/* Page Content */}
        <main className="flex-1 overflow-y-auto p-6 lg:p-8">
          <div className="max-w-6xl mx-auto">
            {isProtected && status === 'unauthenticated' ? (
              <div className="min-h-[360px] flex flex-col items-center justify-center text-center space-y-3">
                <div className="h-6 w-6 animate-spin rounded-full border-2 border-neutral-900 border-t-transparent" />
                <p className="text-xs text-neutral-500 font-medium">
                  {t('shell.redirecting_to_login', 'Sign in required. Redirecting to login...')}
                </p>
              </div>
            ) : isProtected && status === 'loading' ? (
              <div className="min-h-[360px] flex flex-col items-center justify-center text-center space-y-3">
                <div className="h-6 w-6 animate-spin rounded-full border-2 border-neutral-300 border-t-neutral-800" />
                <p className="text-xs text-neutral-400 font-medium">
                  {t('common.loading', 'Loading...')}
                </p>
              </div>
            ) : (
              children
            )}
          </div>
        </main>
      </div>
    </div>
  );
}
