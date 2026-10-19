'use client';

import React, { createContext, useContext, useState, useEffect, useMemo } from 'react';
import axios from 'axios';
import enErrors from '@/locales/en/errors.json';
import viErrors from '@/locales/vi/errors.json';
import enCommon from '@/locales/en/common.json';
import viCommon from '@/locales/vi/common.json';

export type Locale = 'en' | 'vi';

export const SUPPORTED_LOCALES: { code: Locale; name: string; flag: string }[] = [
  { code: 'en', name: 'English', flag: 'EN' },
  { code: 'vi', name: 'Tiếng Việt', flag: 'VI' },
];

export interface FormattedError {
  message: string;
  code?: string;
  fieldErrors?: Record<string, string>;
}

export type TranslationParams = Record<string, string | number>;

interface I18nContextValue {
  locale: Locale;
  setLocale: (locale: Locale) => void;
  t: (key: string, paramsOrFallback?: TranslationParams | string, fallback?: string) => string;
  formatError: (error: unknown) => FormattedError;
}

const errorDictionaries: Record<Locale, Record<string, string>> = {
  en: enErrors,
  vi: viErrors,
};

const commonDictionaries: Record<Locale, Record<string, unknown>> = {
  en: enCommon,
  vi: viCommon,
};

export function getClientLocale(): Locale {
  if (typeof document !== 'undefined') {
    // 1. Check cookie
    const match = document.cookie.match(/(?:^|; )x-locale=([^;]*)/);
    if (match && match[1]) {
      const decoded = decodeURIComponent(match[1]).toLowerCase();
      if (decoded === 'vi' || decoded === 'en') {
        return decoded;
      }
    }

    // 2. Check localStorage
    try {
      const stored = localStorage.getItem('locale');
      if (stored === 'vi' || stored === 'en') {
        return stored;
      }
    } catch {
      // Ignore localStorage errors in SSR/sandboxes
    }

    // 3. Check navigator language
    if (navigator.language && navigator.language.toLowerCase().startsWith('vi')) {
      return 'vi';
    }
  }

  return 'en';
}

const I18nContext = createContext<I18nContextValue>({
  locale: 'en',
  setLocale: () => {},
  t: (key, paramsOrFallback, fallback) => {
    if (typeof paramsOrFallback === 'string') return paramsOrFallback;
    return fallback || key;
  },
  formatError: () => ({ message: 'An unexpected error occurred.' }),
});

export function I18nProvider({
  children,
  initialLocale,
}: {
  children: React.ReactNode;
  initialLocale?: Locale;
}) {
  const [locale, setLocaleState] = useState<Locale>(() => initialLocale || getClientLocale());

  const setLocale = (newLocale: Locale) => {
    setLocaleState(newLocale);
    if (typeof document !== 'undefined') {
      document.cookie = `x-locale=${newLocale}; path=/; max-age=31536000; SameSite=Lax`;
      try {
        localStorage.setItem('locale', newLocale);
      } catch {
        // Ignore
      }
    }
  };

  const t = useMemo(() => {
    return (
      key: string,
      paramsOrFallback?: TranslationParams | string,
      fallback?: string
    ): string => {
      let params: TranslationParams | undefined;
      let effectiveFallback: string | undefined = fallback;

      if (typeof paramsOrFallback === 'string') {
        effectiveFallback = paramsOrFallback;
      } else if (typeof paramsOrFallback === 'object' && paramsOrFallback !== null) {
        params = paramsOrFallback;
      }

      const keys = key.split('.');
      let current: unknown = commonDictionaries[locale];
      for (const k of keys) {
        if (current && typeof current === 'object' && k in current) {
          current = (current as Record<string, unknown>)[k];
        } else {
          current = undefined;
          break;
        }
      }

      let res: string = typeof current === 'string' ? current : effectiveFallback || key;

      if (params) {
        for (const [pKey, pVal] of Object.entries(params)) {
          res = res.replaceAll(`{${pKey}}`, String(pVal));
        }
      }

      return res;
    };
  }, [locale]);

  const formatError = useMemo(() => {
    return (error: unknown): FormattedError => {
      const dict = (errorDictionaries[locale] || errorDictionaries.en) as Record<string, string>;

      if (axios.isAxiosError(error) && error.response?.data) {
        const data = error.response.data as {
          message_code?: string;
          message?: string;
          error_detail?: Array<{ field?: string; message_code?: string; message?: string }>;
        };

        const code = data.message_code;

        // Two-tier resolution: Check FE dictionary first -> then BE localized message -> then fallback
        const message =
          (code && dict[code]) ||
          data.message ||
          (code
            ? `Error: ${code}`
            : dict.ERR_UNEXPECTED || 'An unexpected request error occurred.');

        const fieldErrors: Record<string, string> = {};
        if (data.error_detail && Array.isArray(data.error_detail)) {
          for (const item of data.error_detail) {
            if (item.field) {
              const fieldMsg =
                (item.message_code && dict[item.message_code]) ||
                item.message ||
                dict.ERR_INVALID_FIELD ||
                'Invalid field value';
              fieldErrors[item.field] = fieldMsg;
            }
          }
        }

        return { message, code, fieldErrors };
      }

      if (error instanceof Error) {
        return { message: dict[error.message] || error.message };
      }

      return {
        message: dict.ERR_NETWORK || 'An unexpected network error occurred.',
      };
    };
  }, [locale]);

  const value = useMemo(() => ({ locale, setLocale, t, formatError }), [locale, t, formatError]);

  return <I18nContext.Provider value={value}>{children}</I18nContext.Provider>;
}

export function useI18n(): I18nContextValue {
  return useContext(I18nContext);
}
