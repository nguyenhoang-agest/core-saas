'use client';

import React, { useState, Suspense } from 'react';
import Link from 'next/link';
import { useRouter, useSearchParams } from 'next/navigation';
import { signIn, getSession } from 'next-auth/react';
import { Building2, Eye, EyeOff, Loader2, AlertCircle, Lock, Mail } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';
import { useI18n } from '@/core/i18n/i18n-context';
import enErrors from '@/locales/en/errors.json';
import viErrors from '@/locales/vi/errors.json';

function LoginForm() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const callbackUrl = searchParams.get('callbackUrl') || '/';
  const emailParam = searchParams.get('email') || '';
  const urlError = searchParams.get('error');
  const { t, locale } = useI18n();

  const errorsDict: Record<string, string> = locale === 'vi' ? viErrors : enErrors;

  const [email, setEmail] = useState(emailParam);
  const [password, setPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [loading, setLoading] = useState(false);
  const [errorMessage, setErrorMessage] = useState<string | null>(() => {
    if (!urlError) return null;
    return errorsDict[urlError] || (urlError === 'CredentialsSignin' ? t('auth.invalid_credentials', 'Incorrect email or password. Please verify your credentials.') : urlError);
  });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setErrorMessage(null);

    if (!email || !password) {
      setErrorMessage(t('auth.fill_all_fields', 'Please enter your email and password.'));
      return;
    }

    setLoading(true);

    try {
      const res = await signIn('credentials', {
        email: email.trim().toLowerCase(),
        password,
        redirect: false,
        callbackUrl,
      });

      if (!res || res.error) {
        const errCode = res?.error || 'E010004';
        const localized = errorsDict[errCode] || t('auth.invalid_credentials', 'Incorrect email or password. Please verify your credentials.');
        setErrorMessage(localized);
        setLoading(false);
        return;
      }

      // Sync active tenant cookie and localStorage from resolved session
      try {
        const session = await getSession();
        const tenantId = (session?.user as any)?.tenantId;
        const tenantName = (session?.user as any)?.tenantName || tenantId;
        if (tenantId && typeof document !== 'undefined') {
          document.cookie = `x-tenant-id=${encodeURIComponent(tenantId)}; path=/; max-age=604800; SameSite=Lax`;
          try {
            localStorage.setItem('core_saas_tenant_id', tenantId);
            localStorage.setItem('core_saas_tenant_name', tenantName);
          } catch {
            // ignore
          }
        }
      } catch {
        // ignore
      }

      router.push(callbackUrl);
      router.refresh();
    } catch (err: any) {
      setErrorMessage(t('auth.login_failed', 'An unexpected error occurred during login. Please try again.'));
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen w-full flex flex-col justify-center items-center bg-neutral-50 px-4 py-12 sm:px-6 lg:px-8">
      <div className="sm:mx-auto sm:w-full sm:max-w-md mb-6 text-center">
        <div className="inline-flex h-12 w-12 items-center justify-center rounded-xl bg-neutral-900 text-white shadow-sm mb-3">
          <Building2 className="h-6 w-6 stroke-[1.5]" />
        </div>
        <h1 className="text-2xl font-bold tracking-tight text-neutral-900">Core SaaS</h1>
        <p className="text-sm text-neutral-500 mt-1">Enterprise Multi-Tenant Management Platform</p>
      </div>

      <Card className="w-full max-w-md shadow-sm border-neutral-200 bg-white">
        <CardHeader className="space-y-1 pb-4">
          <CardTitle className="text-xl font-semibold text-neutral-900">
            {t('auth.sign_in', 'Sign in to your account')}
          </CardTitle>
          <CardDescription className="text-xs text-neutral-500">
            {t('auth.sign_in_desc', 'Enter your work email and password to access your workspace.')}
          </CardDescription>
        </CardHeader>

        <form onSubmit={handleSubmit}>
          <CardContent className="space-y-4 pt-0">
            {errorMessage && (
              <div
                role="alert"
                className="flex items-start gap-2.5 p-3 rounded-md bg-red-50 border border-red-200 text-red-700 text-xs"
              >
                <AlertCircle className="h-4 w-4 shrink-0 mt-0.5 text-red-600" />
                <span className="leading-relaxed">{errorMessage}</span>
              </div>
            )}

            <div className="space-y-1.5">
              <Label htmlFor="email" className="text-xs font-medium text-neutral-700">
                {t('auth.email', 'Email Address')}
              </Label>
              <div className="relative">
                <Input
                  id="email"
                  type="email"
                  autoComplete="email"
                  required
                  placeholder="name@company.com"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  className="pl-9 text-xs h-9"
                  disabled={loading}
                />
                <Mail className="h-4 w-4 text-neutral-400 absolute left-3 top-2.5 pointer-events-none" />
              </div>
            </div>

            <div className="space-y-1.5">
              <div className="flex items-center justify-between">
                <Label htmlFor="password" className="text-xs font-medium text-neutral-700">
                  {t('auth.password', 'Password')}
                </Label>
              </div>
              <div className="relative">
                <Input
                  id="password"
                  type={showPassword ? 'text' : 'password'}
                  autoComplete="current-password"
                  required
                  placeholder="••••••••"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  className="pl-9 pr-9 text-xs h-9"
                  disabled={loading}
                />
                <Lock className="h-4 w-4 text-neutral-400 absolute left-3 top-2.5 pointer-events-none" />
                <button
                  type="button"
                  onClick={() => setShowPassword(!showPassword)}
                  className="absolute right-3 top-2.5 text-neutral-400 hover:text-neutral-600 transition-colors"
                  aria-label={showPassword ? 'Hide password' : 'Show password'}
                  tabIndex={-1}
                >
                  {showPassword ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
                </button>
              </div>
            </div>

            <Button
              type="submit"
              className="w-full h-9 text-xs font-medium bg-neutral-900 text-white hover:bg-neutral-800 transition-colors mt-2"
              disabled={loading}
            >
              {loading ? (
                <>
                  <Loader2 className="h-3.5 w-3.5 animate-spin mr-2" />
                  {t('auth.signing_in', 'Signing in...')}
                </>
              ) : (
                t('auth.sign_in_button', 'Sign In')
              )}
            </Button>
          </CardContent>

          <CardFooter className="flex flex-col space-y-2 border-t border-neutral-100 p-4 bg-neutral-50/50 rounded-b-lg text-center">
            <p className="text-xs text-neutral-600">
              {t('auth.no_account', "Don't have an organization?")}{' '}
              <Link
                href="/register"
                className="font-medium text-neutral-900 hover:underline transition-colors"
              >
                {t('auth.create_workspace', 'Create a workspace')}
              </Link>
            </p>
          </CardFooter>
        </form>
      </Card>
    </div>
  );
}

export default function LoginPage() {
  return (
    <Suspense fallback={<div className="min-h-screen flex items-center justify-center bg-neutral-50 text-neutral-500 text-sm">Loading...</div>}>
      <LoginForm />
    </Suspense>
  );
}
