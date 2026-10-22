'use client';

import React, { useState } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { signIn } from 'next-auth/react';
import { Building2, Eye, EyeOff, Loader2, AlertCircle, Lock, Mail, User } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';
import { useI18n } from '@/core/i18n/i18n-context';
import enErrors from '@/locales/en/errors.json';
import viErrors from '@/locales/vi/errors.json';

export default function RegisterPage() {
  const router = useRouter();
  const { t, locale } = useI18n();

  const [orgName, setOrgName] = useState('');
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [loading, setLoading] = useState(false);
  const [errorMessage, setErrorMessage] = useState<string | null>(null);

  const errorsDict: Record<string, string> = locale === 'vi' ? viErrors : enErrors;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setErrorMessage(null);

    // Client-side validation
    if (orgName.trim().length < 3) {
      setErrorMessage(t('auth.org_name_min', 'Organization name must be at least 3 characters.'));
      return;
    }
    if (!name.trim()) {
      setErrorMessage(t('auth.name_required', 'Please enter your full name.'));
      return;
    }
    if (!email.trim()) {
      setErrorMessage(t('auth.email_required', 'Please enter your work email.'));
      return;
    }
    if (password.length < 8) {
      setErrorMessage(t('auth.password_min', 'Password must be at least 8 characters.'));
      return;
    }
    if (password !== confirmPassword) {
      setErrorMessage(t('auth.password_mismatch', 'Passwords do not match. Please verify.'));
      return;
    }

    setLoading(true);

    try {
      // 1. Call Register API
      const res = await fetch('/api/register', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          org_name: orgName.trim(),
          name: name.trim(),
          email: email.trim().toLowerCase(),
          password,
        }),
      });

      const data = await res.json().catch(() => null);

      if (!res.ok) {
        const errCode = data?.code || 'E030001';
        const localized = errorsDict[errCode] || data?.message || t('auth.register_failed', 'Failed to register workspace.');
        setErrorMessage(localized);
        setLoading(false);
        return;
      }

      // 2. Persist tenant ID cookie and name immediately for the new workspace
      const newTenantId = data?.data?.tenant_id;
      const newOrgName = orgName.trim();
      if (newTenantId && typeof document !== 'undefined') {
        document.cookie = `x-tenant-id=${encodeURIComponent(newTenantId)}; path=/; max-age=604800; SameSite=Lax`;
        try {
          localStorage.setItem('core_saas_tenant_id', newTenantId);
          localStorage.setItem('core_saas_tenant_name', newOrgName);
        } catch {
          // ignore
        }
      }

      // 3. Automatically sign in with NextAuth upon registration
      const loginRes = await signIn('credentials', {
        email: email.trim().toLowerCase(),
        password,
        redirect: false,
        callbackUrl: '/',
      });

      if (loginRes?.error) {
        // If auto-login fails, redirect to login page with pre-filled email
        router.push(`/login?email=${encodeURIComponent(email.trim().toLowerCase())}`);
        return;
      }

      router.push('/');
      router.refresh();
    } catch (err: any) {
      setErrorMessage(t('auth.register_error', 'An unexpected error occurred. Please try again.'));
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
        <p className="text-sm text-neutral-500 mt-1">Create an enterprise workspace</p>
      </div>

      <Card className="w-full max-w-md shadow-sm border-neutral-200 bg-white">
        <CardHeader className="space-y-1 pb-4">
          <CardTitle className="text-xl font-semibold text-neutral-900">
            {t('auth.create_account', 'Create your organization')}
          </CardTitle>
          <CardDescription className="text-xs text-neutral-500">
            {t('auth.create_account_desc', 'Set up your workspace and initial Owner account.')}
          </CardDescription>
        </CardHeader>

        <form onSubmit={handleSubmit}>
          <CardContent className="space-y-3.5 pt-0">
            {errorMessage && (
              <div
                role="alert"
                className="flex items-start gap-2.5 p-3 rounded-md bg-red-50 border border-red-200 text-red-700 text-xs"
              >
                <AlertCircle className="h-4 w-4 shrink-0 mt-0.5 text-red-600" />
                <span className="leading-relaxed">{errorMessage}</span>
              </div>
            )}

            <div className="space-y-1">
              <Label htmlFor="orgName" className="text-xs font-medium text-neutral-700">
                {t('auth.org_name', 'Organization / Workspace Name')}
              </Label>
              <div className="relative">
                <Input
                  id="orgName"
                  required
                  placeholder="Acme Corp"
                  value={orgName}
                  onChange={(e) => setOrgName(e.target.value)}
                  className="pl-9 text-xs h-9"
                  disabled={loading}
                />
                <Building2 className="h-4 w-4 text-neutral-400 absolute left-3 top-2.5 pointer-events-none" />
              </div>
            </div>

            <div className="space-y-1">
              <Label htmlFor="name" className="text-xs font-medium text-neutral-700">
                {t('auth.full_name', 'Your Full Name')}
              </Label>
              <div className="relative">
                <Input
                  id="name"
                  required
                  placeholder="Alice Smith"
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  className="pl-9 text-xs h-9"
                  disabled={loading}
                />
                <User className="h-4 w-4 text-neutral-400 absolute left-3 top-2.5 pointer-events-none" />
              </div>
            </div>

            <div className="space-y-1">
              <Label htmlFor="email" className="text-xs font-medium text-neutral-700">
                {t('auth.work_email', 'Work Email')}
              </Label>
              <div className="relative">
                <Input
                  id="email"
                  type="email"
                  autoComplete="email"
                  required
                  placeholder="alice@company.com"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  className="pl-9 text-xs h-9"
                  disabled={loading}
                />
                <Mail className="h-4 w-4 text-neutral-400 absolute left-3 top-2.5 pointer-events-none" />
              </div>
            </div>

            <div className="space-y-1">
              <Label htmlFor="password" className="text-xs font-medium text-neutral-700">
                {t('auth.password', 'Password')}
              </Label>
              <div className="relative">
                <Input
                  id="password"
                  type={showPassword ? 'text' : 'password'}
                  autoComplete="new-password"
                  required
                  placeholder="Minimum 8 characters"
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

            <div className="space-y-1">
              <Label htmlFor="confirmPassword" className="text-xs font-medium text-neutral-700">
                {t('auth.confirm_password', 'Confirm Password')}
              </Label>
              <div className="relative">
                <Input
                  id="confirmPassword"
                  type={showPassword ? 'text' : 'password'}
                  autoComplete="new-password"
                  required
                  placeholder="Re-type password"
                  value={confirmPassword}
                  onChange={(e) => setConfirmPassword(e.target.value)}
                  className="pl-9 text-xs h-9"
                  disabled={loading}
                />
                <Lock className="h-4 w-4 text-neutral-400 absolute left-3 top-2.5 pointer-events-none" />
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
                  {t('auth.creating_org', 'Creating organization...')}
                </>
              ) : (
                t('auth.create_org_button', 'Create Workspace & Account')
              )}
            </Button>
          </CardContent>

          <CardFooter className="flex flex-col space-y-2 border-t border-neutral-100 p-4 bg-neutral-50/50 rounded-b-lg text-center">
            <p className="text-xs text-neutral-600">
              {t('auth.already_have_account', 'Already have an account?')}{' '}
              <Link
                href="/login"
                className="font-medium text-neutral-900 hover:underline transition-colors"
              >
                {t('auth.sign_in_link', 'Sign in')}
              </Link>
            </p>
          </CardFooter>
        </form>
      </Card>
    </div>
  );
}
