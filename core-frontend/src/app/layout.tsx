import type { Metadata } from 'next';
import './globals.css';
import { PlatformProvider } from '@/core/providers/platform-provider';
import { AppShell } from '@/components/app-shell';
import { appConfig } from '@/../config/app.config';

export const metadata: Metadata = {
  title: appConfig.branding.appName,
  description: appConfig.branding.tagline,
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body>
        <PlatformProvider modules={appConfig.enabledModules}>
          <AppShell>{children}</AppShell>
        </PlatformProvider>
      </body>
    </html>
  );
}
