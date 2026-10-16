import type { Metadata } from 'next';
import './globals.css';
import { PlatformProvider } from '@/core/providers/platform-provider';
import { AdminShell } from '@/components/admin-shell';
import { appConfig } from '@/../config/app.config';

export const metadata: Metadata = {
  title: `${appConfig.branding.appName} - Admin Console`,
  description: appConfig.branding.tagline,
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className="bg-neutral-50/50 text-neutral-900 antialiased">
        <PlatformProvider modules={appConfig.enabledModules}>
          <AdminShell>{children}</AdminShell>
        </PlatformProvider>
      </body>
    </html>
  );
}
