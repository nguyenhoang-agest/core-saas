'use client';

import React, { useState, useRef, useEffect, useMemo } from 'react';
import { Building2, Check, ChevronsUpDown, Search, Database } from 'lucide-react';
import { useQueryClient } from '@tanstack/react-query';
import { useTenant } from '@/core/tenant/tenant-context';
import { useTenants, type Tenant } from '@/modules/tenant';
import { cn } from '@/lib/utils';
import { Badge } from '@/components/ui/badge';

export function TenantSwitcher() {
  const { tenant, setTenant } = useTenant();
  const queryClient = useQueryClient();
  const [open, setOpen] = useState(false);
  const [search, setSearch] = useState('');
  const containerRef = useRef<HTMLDivElement>(null);

  const { data: tenantsData, isLoading } = useTenants({ page: 1, pageSize: 50 });
  const tenants = tenantsData?.items || [];

  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      if (containerRef.current && !containerRef.current.contains(event.target as Node)) {
        setOpen(false);
      }
    }
    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  const filteredTenants = useMemo(() => {
    if (!search.trim()) return tenants;
    const query = search.toLowerCase();
    return tenants.filter(
      (t) =>
        t.name.toLowerCase().includes(query) ||
        t.id.toLowerCase().includes(query)
    );
  }, [tenants, search]);

  const handleSelectTenant = (selected: Tenant) => {
    setTenant({
      id: selected.id,
      name: selected.name,
      isolateMode: selected.isolate_mode,
      theme: tenant.theme,
    });
    // Invalidate queries so modules re-fetch data for the new tenant
    queryClient.clear();
    setOpen(false);
  };

  return (
    <div className="relative inline-block text-left" ref={containerRef}>
      <button
        type="button"
        onClick={() => setOpen(!open)}
        className="flex items-center gap-2 px-3 py-1.5 rounded-md border border-neutral-200 bg-white text-neutral-800 hover:border-neutral-300 hover:bg-neutral-50 text-xs font-medium transition-colors shadow-none max-w-[200px]"
        title="Switch active organization / tenant"
      >
        <Building2 className="h-3.5 w-3.5 text-neutral-600 shrink-0" />
        <span className="truncate">{tenant.name || tenant.id}</span>
        <ChevronsUpDown className="h-3 w-3 text-neutral-400 shrink-0 ml-auto" />
      </button>

      {open && (
        <div className="absolute right-0 mt-1.5 w-72 rounded-lg border border-neutral-200 bg-white p-2 text-xs shadow-lg z-50 animate-in fade-in-0 zoom-in-95">
          <div className="flex items-center gap-1.5 px-2 pb-2 mb-1 border-b border-neutral-100">
            <Search className="h-3.5 w-3.5 text-neutral-400 shrink-0" />
            <input
              type="text"
              value ={search}
              onChange={(e) => setSearch(e.target.value)}
              placeholder="Search workspaces..."
              className="w-full text-xs bg-transparent border-none outline-none placeholder:text-neutral-400"
              autoFocus
            />
          </div>

          <div className="max-h-56 overflow-y-auto space-y-1 py-1">
            {isLoading ? (
              <div className="py-3 text-center text-neutral-400 text-xs">Loading workspaces...</div>
            ) : filteredTenants.length === 0 ? (
              <div className="py-3 text-center text-neutral-400 text-xs">No workspaces found</div>
            ) : (
              filteredTenants.map((item) => {
                const isSelected = item.id === tenant.id;
                return (
                  <button
                    key={item.id}
                    type="button"
                    onClick={() => handleSelectTenant(item)}
                    className={cn(
                      'flex items-center justify-between w-full px-2.5 py-2 rounded-md text-left transition-colors',
                      isSelected
                        ? 'bg-neutral-100 text-neutral-900 font-medium'
                        : 'text-neutral-700 hover:bg-neutral-50 hover:text-neutral-900'
                    )}
                  >
                    <div className="min-w-0 pr-2">
                      <div className="flex items-center gap-1.5">
                        <span className="truncate font-medium">{item.name}</span>
                        {item.isolate_mode === 'dedicated' && (
                          <Badge variant="outline" className="text-[9px] px-1 py-0 border-amber-300 text-amber-700 bg-amber-50">
                            Dedicated
                          </Badge>
                        )}
                        {item.status === 'suspended' && (
                          <Badge variant="outline" className="text-[9px] px-1 py-0 border-red-300 text-red-700 bg-red-50">
                            Suspended
                          </Badge>
                        )}
                      </div>
                      <div className="text-[10px] text-neutral-400 font-mono truncate">
                        {item.id}
                      </div>
                    </div>
                    {isSelected && <Check className="h-4 w-4 text-neutral-900 shrink-0 ml-2" />}
                  </button>
                );
              })
            )}
          </div>
        </div>
      )}
    </div>
  );
}
