'use client';

import React from 'react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import {
  Building2,
  LayoutDashboard,
  Settings,
  ShieldAlert,
  Layers,
  ChevronRight,
  ExternalLink,
} from 'lucide-react';
import { cn } from '@/lib/utils';
import { useTenant } from '@/core/tenant/tenant-context';
import { useI18n } from '@/core/i18n/i18n-context';
import { appConfig } from '@/../config/app.config';
import { LocaleSwitcher } from './locale-switcher';

interface NavItem {
  titleKey: string;
  fallback: string;
  href: string;
  icon: React.ComponentType<{ className?: string }>;
  badge?: string;
}

const navItems: NavItem[] = [
  {
    titleKey: 'nav.overview',
    fallback: 'Overview',
    href: '/',
    icon: LayoutDashboard,
  },
  {
    titleKey: 'nav.tenants',
    fallback: 'Tenants',
    href: '/tenants',
    icon: Building2,
    badge: 'Core',
  },
  {
    titleKey: 'nav.settings',
    fallback: 'Settings',
    href: '/settings',
    icon: Settings,
  },
];

export function AdminShell({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();
  const { tenant } = useTenant();
  const { t } = useI18n();

  const getBreadcrumbTitle = () => {
    if (pathname === '/') return t('nav.dashboard', 'Dashboard');
    if (pathname.startsWith('/tenants')) return t('nav.tenants', 'Tenants');
    if (pathname.startsWith('/settings')) return t('nav.settings', 'Settings');
    return pathname.replace('/', '');
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
            {navItems.map((item) => {
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

        {/* Sidebar Footer: Current Tenant Context */}
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
      </aside>

      {/* Main Content Area */}
      <div className="flex-1 flex flex-col min-w-0">
        {/* Top Header */}
        <header className="h-14 border-b border-neutral-200 bg-white/80 backdrop-blur-sm sticky top-0 z-30 flex items-center justify-between px-6 lg:px-8">
          {/* Breadcrumbs */}
          <div className="flex items-center space-x-1.5 text-xs text-neutral-700">
            <span>{t('shell.admin', 'Admin')}</span>
            <ChevronRight className="h-3 w-3 text-neutral-600 stroke-[1.5]" />
            <span className="text-neutral-900 font-medium capitalize">
              {getBreadcrumbTitle()}
            </span>
          </div>

          {/* Right Header Badges */}
          <div className="flex items-center space-x-3 text-xs">
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
          </div>
        </header>

        {/* Page Content */}
        <main className="flex-1 overflow-y-auto p-6 lg:p-8">
          <div className="max-w-6xl mx-auto">{children}</div>
        </main>
      </div>
    </div>
  );
}
