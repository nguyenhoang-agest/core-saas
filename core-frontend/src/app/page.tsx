'use client';

import React from 'react';
import Link from 'next/link';
import {
  Building2,
  Users,
  Database,
  ShieldCheck,
  Layers,
  ArrowRight,
  Zap,
  Lock,
  Globe,
  Radio,
  Sliders,
  CreditCard,
  Settings,
  ExternalLink,
  CheckCircle2,
  LogIn,
} from 'lucide-react';
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { useSession } from 'next-auth/react';
import { useTenant } from '@/core/tenant/tenant-context';
import { useI18n } from '@/core/i18n/i18n-context';

export default function HomePage() {
  const { tenant } = useTenant();
  const { t } = useI18n();
  const { data: session, status } = useSession();
  const isAuthenticated = status === 'authenticated';
  const isHost = !!(session?.user as any)?.isHost;

  const coreModules = [
    {
      title: 'Tenant Workspaces',
      description: 'Provision and inspect organizations with shared or dedicated database isolation.',
      href: '/tenants',
      icon: Building2,
      badge: 'Core Engine',
      badgeVariant: 'secondary' as const,
      features: ['Dynamic DB Routing', 'Dedicated TLS Encryption', 'Idle Pool Reaping'],
    },
    {
      title: 'Users & Access Control',
      description: 'Manage workspace members, assign RBAC roles, and inspect fine-grained permissions.',
      href: '/users',
      icon: Users,
      badge: 'RBAC Engine',
      badgeVariant: 'secondary' as const,
      features: ['System vs Custom Roles', 'Anti-Lockout Protection', 'Real-Time Invalidation'],
    },
    {
      title: 'Billing & Subscriptions',
      description: 'Manage organization subscription tiers, usage quotas, and download invoices.',
      href: '/billing',
      icon: CreditCard,
      badge: 'Commerce',
      badgeVariant: 'outline' as const,
      features: ['Tiered Pricing Plans', 'Add-on Quotas', 'Invoice History'],
    },
    {
      title: 'Settings & Branding',
      description: 'Customize workspace themes, color tokens, and configure localization preferences.',
      href: '/settings',
      icon: Settings,
      badge: 'System',
      badgeVariant: 'outline' as const,
      features: ['Dynamic Custom Properties', 'En/Vi Localization', 'Domain Branding'],
    },
  ];

  const architecturePillars = [
    {
      icon: Database,
      title: 'Dynamic DB Routing',
      subtitle: 'pkg/db/router',
      description:
        'Context-aware connection pool router supporting shared RLS multi-tenancy and dedicated tenant databases with mandatory TLS.',
    },
    {
      icon: Lock,
      title: 'Granular RBAC & Anti-Lockout',
      subtitle: 'modules/user & pkg/rbac',
      description:
        'Immutable system roles (Owner, Admin, Member, Viewer) with Anti-Lockout owner protection and strict privilege escalation checks.',
    },
    {
      icon: Radio,
      title: 'Real-Time Cache Invalidation',
      subtitle: 'PostgreSQL LISTEN/NOTIFY',
      description:
        'Instant distributed authorization cache invalidation via PostgreSQL NOTIFY authz_invalidation and versioned epochs.',
    },
    {
      icon: ShieldCheck,
      title: '3-Layer Error Architecture',
      subtitle: 'E031001+ / E032001+',
      description:
        'Typed domain exceptions with module-scoped error code offsets, field validation details, and two-tier i18n dictionary lookup.',
    },
  ];

  return (
    <div className="space-y-6">
      {/* Hero Welcome Banner */}
      <div className="rounded-xl border border-neutral-200 bg-white p-6 lg:p-8 space-y-5">
        <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
          <div className="space-y-1.5">
            <div className="inline-flex items-center gap-1.5 px-2 py-0.5 rounded border border-neutral-200 bg-neutral-50 text-[11px] font-semibold uppercase tracking-wider text-neutral-700">
              <Zap className="h-3 w-3 text-neutral-900" />
              <span>Core SaaS Platform Admin • v1.0</span>
            </div>
            <h1 className="text-2xl font-bold tracking-tight text-neutral-900">
              Modular Multi-Tenant SaaS Console
            </h1>
            <p className="text-xs text-neutral-500 max-w-2xl leading-relaxed">
              Enterprise-grade multi-tenant foundation engineered with Go Fiber v3, Next.js 16,
              Tailwind CSS, and TanStack Query v5. Implements Flat 2.0 Monochrome design with strict
              modular boundaries and zero cross-tenant leakage.
            </p>
          </div>

          <div className="flex flex-wrap items-center gap-2 shrink-0">
            {isAuthenticated ? (
              <>
                {isHost && (
                  <Link href="/tenants">
                    <Button className="gap-2">
                      <Building2 className="h-4 w-4" />
                      <span>Manage Tenants</span>
                      <ArrowRight className="h-3.5 w-3.5" />
                    </Button>
                  </Link>
                )}
                <Link href="/users">
                  <Button variant={isHost ? 'outline' : 'default'} className="gap-2">
                    <Users className="h-4 w-4" />
                    <span>Manage Users & Roles</span>
                    {!isHost && <ArrowRight className="h-3.5 w-3.5" />}
                  </Button>
                </Link>
                {!isHost && (
                  <Link href="/settings">
                    <Button variant="outline" className="gap-2">
                      <Settings className="h-4 w-4" />
                      <span>Workspace Settings</span>
                    </Button>
                  </Link>
                )}
              </>
            ) : (
              <>
                <Link href="/login">
                  <Button className="gap-2">
                    <LogIn className="h-4 w-4" />
                    <span>Sign In to Console</span>
                    <ArrowRight className="h-3.5 w-3.5" />
                  </Button>
                </Link>
                <Link href="/register">
                  <Button variant="outline" className="gap-2">
                    <span>Create Organization</span>
                  </Button>
                </Link>
              </>
            )}
          </div>
        </div>

        {/* Current Active Workspace Indicator Bar */}
        <div className="pt-4 border-t border-neutral-100 flex flex-wrap items-center justify-between gap-4 text-xs">
          {isAuthenticated ? (
            <div className="flex items-center gap-2">
              <span className="text-neutral-500">Active Workspace:</span>
              <div className="inline-flex items-center gap-1.5 font-semibold text-neutral-900 bg-neutral-100 px-2 py-0.5 rounded border border-neutral-200">
                <span className="h-2 w-2 rounded-full bg-emerald-500" />
                <span>{tenant.name}</span>
              </div>
              <span className="font-mono text-[11px] text-neutral-400">
                ({tenant.id === 'default' ? 'Default Host' : tenant.id})
              </span>
            </div>
          ) : (
            <div className="flex items-center gap-2">
              <span className="text-neutral-500">System Mode:</span>
              <div className="inline-flex items-center gap-1.5 font-semibold text-neutral-900 bg-neutral-100 px-2 py-0.5 rounded border border-neutral-200">
                <span className="h-2 w-2 rounded-full bg-amber-500" />
                <span>Public Preview • Guest</span>
              </div>
              <span className="text-[11px] text-neutral-400">
                (Sign in required to access management consoles)
              </span>
            </div>
          )}

          <div className="flex items-center gap-3">
            <div className="flex items-center gap-1.5 text-neutral-600">
              <Database className="h-3.5 w-3.5 text-neutral-500" />
              <span>Isolation:</span>
              <Badge variant="outline" className="capitalize text-[10px] py-0 px-1.5 font-medium">
                {tenant.isolateMode || 'shared'} DB
              </Badge>
            </div>
            <div className="flex items-center gap-1.5 text-neutral-600">
              <Globe className="h-3.5 w-3.5 text-neutral-500" />
              <span>Sync:</span>
              <span className="text-[11px] font-mono text-neutral-700">Cookie + Header</span>
            </div>
          </div>
        </div>
      </div>

      {/* Core Modules Grid */}
      <div>
        <div className="flex items-center justify-between mb-3">
          <h2 className="text-xs font-semibold uppercase tracking-wider text-neutral-500">
            Core Modules & Workspaces
          </h2>
          <span className="text-[11px] text-neutral-400">
            {coreModules.length} Modules Available
          </span>
        </div>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
          {coreModules.map((m) => {
            const Icon = m.icon;
            return (
              <Card
                key={m.title}
                className="hover:border-neutral-300 transition-colors flex flex-col justify-between"
              >
                <CardHeader className="pb-3">
                  <div className="flex items-center justify-between mb-2">
                    <div className="p-2 rounded-md border border-neutral-200 bg-neutral-50 text-neutral-900">
                      <Icon className="h-4 w-4 stroke-[1.75]" />
                    </div>
                    <Badge variant={m.badgeVariant} className="text-[10px] font-semibold uppercase">
                      {m.badge}
                    </Badge>
                  </div>
                  <CardTitle className="text-sm font-semibold">{m.title}</CardTitle>
                  <CardDescription className="text-xs text-neutral-500 mt-1 line-clamp-2">
                    {m.description}
                  </CardDescription>
                </CardHeader>
                <CardContent className="pt-0 space-y-3">
                  <div className="space-y-1">
                    {m.features.map((feat) => (
                      <div
                        key={feat}
                        className="flex items-center gap-1.5 text-[11px] text-neutral-600"
                      >
                        <CheckCircle2 className="h-3 w-3 text-neutral-400 shrink-0" />
                        <span>{feat}</span>
                      </div>
                    ))}
                  </div>
                  <div className="pt-2 border-t border-neutral-100">
                    {isAuthenticated ? (
                      <Link
                        href={m.href}
                        className="inline-flex items-center text-xs font-medium text-neutral-900 hover:text-neutral-700 gap-1 group"
                      >
                        <span>Open Module</span>
                        <ArrowRight className="h-3 w-3 stroke-[1.5] group-hover:translate-x-0.5 transition-transform" />
                      </Link>
                    ) : (
                      <Link
                        href={`/login?callbackUrl=${encodeURIComponent(m.href)}`}
                        className="inline-flex items-center text-xs font-medium text-neutral-700 hover:text-neutral-900 gap-1 group"
                      >
                        <Lock className="h-3 w-3 text-neutral-400" />
                        <span>Sign In to Access</span>
                        <ArrowRight className="h-3 w-3 stroke-[1.5] group-hover:translate-x-0.5 transition-transform" />
                      </Link>
                    )}
                  </div>
                </CardContent>
              </Card>
            );
          })}
        </div>
      </div>

      {/* Architecture Highlights */}
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <div className="space-y-1">
              <CardTitle className="text-sm flex items-center gap-2">
                <Layers className="h-4 w-4 text-neutral-900 stroke-[1.75]" />
                <span>Architecture & Security Engineering Pillars</span>
              </CardTitle>
              <CardDescription className="text-xs">
                Hardened infrastructure patterns powering multi-tenant reliability and compliance.
              </CardDescription>
            </div>
          </div>
        </CardHeader>
        <CardContent className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 text-xs">
          {architecturePillars.map((p) => {
            const Icon = p.icon;
            return (
              <div
                key={p.title}
                className="rounded border border-neutral-200 p-3.5 bg-neutral-50/40 space-y-1.5 flex flex-col justify-between"
              >
                <div className="space-y-1">
                  <div className="flex items-center gap-2">
                    <Icon className="h-4 w-4 text-neutral-800 shrink-0 stroke-[1.75]" />
                    <span className="font-semibold text-neutral-900">{p.title}</span>
                  </div>
                  <div className="font-mono text-[10px] text-neutral-400">{p.subtitle}</div>
                  <p className="text-neutral-500 text-[11px] leading-relaxed pt-1">
                    {p.description}
                  </p>
                </div>
              </div>
            );
          })}
        </CardContent>
      </Card>
    </div>
  );
}
