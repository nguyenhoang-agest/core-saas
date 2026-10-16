import React from 'react';
import Link from 'next/link';
import {
  Building2,
  Database,
  ShieldCheck,
  Layers,
  ArrowRight,
  ExternalLink,
  Server,
  Zap,
} from 'lucide-react';
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';

export default function HomePage() {
  const features = [
    {
      title: 'Tenant Workspaces',
      description: 'Provision and inspect organizations with shared or dedicated database isolation.',
      href: '/tenants',
      icon: Building2,
      badge: 'Implemented',
    },
    {
      title: 'Dynamic DB Routing',
      description: 'Tenant-aware PostgreSQL connection pool management with idle pool reaping.',
      href: '/tenants',
      icon: Database,
      badge: 'Core Backend',
    },
    {
      title: '3-Layer Error Handling',
      description: 'Typed exceptions, field-level validation, and module-scoped error code offsets.',
      href: '/tenants',
      icon: ShieldCheck,
      badge: 'E031001+',
    },
  ];

  return (
    <div className="space-y-6">
      {/* Hero Welcome Banner */}
      <div className="rounded-xl border border-neutral-200 bg-white p-6 lg:p-8 space-y-4">
        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
          <div className="space-y-1">
            <div className="inline-flex items-center gap-1.5 px-2 py-0.5 rounded border border-neutral-200 bg-neutral-50 text-[11px] font-semibold uppercase tracking-wider text-neutral-700">
              <Zap className="h-3 w-3 text-neutral-900" />
              <span>Core SaaS Platform Admin</span>
            </div>
            <h1 className="text-2xl font-bold tracking-tight text-neutral-900">
              Modular Multi-Tenant SaaS Console
            </h1>
            <p className="text-xs text-neutral-500 max-w-2xl">
              Engineered with Next.js 16, Tailwind CSS, TanStack Query v5, and Go Fiber. Designed
              with an Outlined Minimalist / Flat 2.0 Monochrome UI for enterprise observability.
            </p>
          </div>

          <Link href="/tenants">
            <Button className="gap-2 shrink-0">
              <Building2 className="h-4 w-4" />
              <span>Manage Tenants</span>
              <ArrowRight className="h-3.5 w-3.5" />
            </Button>
          </Link>
        </div>
      </div>

      {/* Feature Cards Grid */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        {features.map((f) => {
          const Icon = f.icon;
          return (
            <Card key={f.title} className="hover:border-neutral-300 transition-colors flex flex-col justify-between">
              <CardHeader>
                <div className="flex items-center justify-between mb-2">
                  <div className="p-2 rounded-md border border-neutral-200 bg-neutral-50 text-neutral-900">
                    <Icon className="h-4 w-4 stroke-[1.75]" />
                  </div>
                  <Badge variant="outline" className="text-[10px] font-semibold uppercase">
                    {f.badge}
                  </Badge>
                </div>
                <CardTitle className="text-sm font-semibold">{f.title}</CardTitle>
                <CardDescription className="text-xs text-neutral-500 mt-1">
                  {f.description}
                </CardDescription>
              </CardHeader>
              <CardContent className="pt-0">
                <Link
                  href={f.href}
                  className="inline-flex items-center text-xs font-medium text-neutral-900 hover:text-neutral-700 gap-1"
                >
                  <span>Open Module</span>
                  <ArrowRight className="h-3 w-3 stroke-[1.5]" />
                </Link>
              </CardContent>
            </Card>
          );
        })}
      </div>

      {/* Architecture Highlights */}
      <Card>
        <CardHeader>
          <CardTitle className="text-sm flex items-center gap-2">
            <Layers className="h-4 w-4 text-neutral-900 stroke-[1.75]" />
            <span>Architecture & Design Pillars</span>
          </CardTitle>
          <CardDescription className="text-xs">
            Core principles defining the administrative experience of Core SaaS.
          </CardDescription>
        </CardHeader>
        <CardContent className="grid grid-cols-1 sm:grid-cols-3 gap-4 text-xs">
          <div className="rounded border border-neutral-200 p-3 bg-neutral-50/30 space-y-1">
            <div className="font-semibold text-neutral-900">Border-focused Minimalist</div>
            <p className="text-neutral-500 text-[11px] leading-relaxed">
              Thin 1px strokes replace heavy drop shadows for clean, uncluttered enterprise information hierarchy.
            </p>
          </div>
          <div className="rounded border border-neutral-200 p-3 bg-neutral-50/30 space-y-1">
            <div className="font-semibold text-neutral-900">Monochrome Grayscale</div>
            <p className="text-neutral-500 text-[11px] leading-relaxed">
              100% grayscale palette keeping focus on data, reserving status colors solely for live indicators.
            </p>
          </div>
          <div className="rounded border border-neutral-200 p-3 bg-neutral-50/30 space-y-1">
            <div className="font-semibold text-neutral-900">Isolated Multi-Tenancy</div>
            <p className="text-neutral-500 text-[11px] leading-relaxed">
              Every request carries tenant context through cookies and headers to dynamically route queries.
            </p>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
