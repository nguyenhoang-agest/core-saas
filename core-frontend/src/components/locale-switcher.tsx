'use client';

import React, { useState, useRef, useEffect } from 'react';
import { Globe, Check } from 'lucide-react';
import { useI18n, SUPPORTED_LOCALES, type Locale } from '@/core/i18n/i18n-context';
import { cn } from '@/lib/utils';

export function LocaleSwitcher() {
  const { locale, setLocale } = useI18n();
  const [open, setOpen] = useState(false);
  const containerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      if (containerRef.current && !containerRef.current.contains(event.target as Node)) {
        setOpen(false);
      }
    }
    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  const current = SUPPORTED_LOCALES.find((l) => l.code === locale) || SUPPORTED_LOCALES[0];

  return (
    <div className="relative inline-block text-left" ref={containerRef}>
      <button
        type="button"
        onClick={() => setOpen(!open)}
        className="flex items-center gap-1.5 px-2.5 py-1 rounded-md border border-neutral-200 bg-white text-neutral-800 hover:border-neutral-300 hover:bg-neutral-50 text-xs font-medium transition-colors shadow-none"
        title="Change application language"
      >
        <Globe className="h-3.5 w-3.5 text-neutral-600 stroke-[1.5]" />
        <span>{current.flag}</span>
      </button>

      {open && (
        <div className="absolute right-0 mt-1.5 w-36 rounded-lg border border-neutral-200 bg-white p-1 text-xs shadow-md z-50 duration-150 animate-in fade-in-0 zoom-in-95">
          <div className="px-2 py-1 text-[10px] font-semibold uppercase tracking-wider text-neutral-600">
            Select Language
          </div>
          {SUPPORTED_LOCALES.map((item) => {
            const isSelected = item.code === locale;
            return (
              <button
                key={item.code}
                type="button"
                onClick={() => {
                  setLocale(item.code as Locale);
                  setOpen(false);
                }}
                className={cn(
                  'flex items-center justify-between w-full px-2 py-1.5 rounded-md text-xs font-medium transition-colors',
                  isSelected
                    ? 'bg-neutral-100 text-neutral-900 font-semibold'
                    : 'text-neutral-600 hover:bg-neutral-50 hover:text-neutral-900'
                )}
              >
                <span>{item.name}</span>
                {isSelected && <Check className="h-3.5 w-3.5 text-neutral-900" />}
              </button>
            );
          })}
        </div>
      )}
    </div>
  );
}
