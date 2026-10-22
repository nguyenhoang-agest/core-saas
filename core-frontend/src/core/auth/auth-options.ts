import type { AuthOptions, Session, User } from 'next-auth';
import type { JWT } from 'next-auth/jwt';
import CredentialsProvider from 'next-auth/providers/credentials';

export interface UserWithPermissions extends User {
  id: string;
  role?: string;
  permissions?: string[];
  tenantId?: string;
  tenantName?: string;
  isHost?: boolean;
  accessToken?: string;
  refreshToken?: string;
}

export interface CustomSession extends Session {
  user: {
    id?: string;
    name?: string | null;
    email?: string | null;
    image?: string | null;
    role?: string;
    permissions?: string[];
    tenantId?: string;
    tenantName?: string;
    isHost?: boolean;
  };
  accessToken?: string;
}

const BACKEND_URL =
  process.env.BACKEND_INTERNAL_URL ||
  (process.env.NEXT_PUBLIC_API_URL && process.env.NEXT_PUBLIC_API_URL.startsWith('http')
    ? process.env.NEXT_PUBLIC_API_URL
    : 'http://localhost:8080');

/**
 * Reusable NextAuth configuration options connecting to Core-SaaS Auth API.
 */
export const defaultAuthOptions: AuthOptions = {
  providers: [
    CredentialsProvider({
      name: 'Credentials',
      credentials: {
        email: { label: 'Email', type: 'email' },
        password: { label: 'Password', type: 'password' },
      },
      async authorize(credentials) {
        if (!credentials?.email || !credentials?.password) return null;

        try {
          // 1. Call Backend POST /auth/login
          const loginRes = await fetch(`${BACKEND_URL}/auth/login`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
              email: credentials.email,
              password: credentials.password,
            }),
          });

          if (!loginRes.ok) {
            const errData = await loginRes.json().catch(() => null);
            throw new Error(errData?.code || 'E010004');
          }

          const loginPayload = await loginRes.json();
          const authData = loginPayload.data;
          const accessToken = authData.access_token;
          const refreshToken = authData.refresh_token;

          // 2. Fetch real dynamic permissions and tenant details via GET /auth/me
          let permissions: string[] = [];
          let isHost = false;
          let tenantName = '';
          try {
            const meRes = await fetch(`${BACKEND_URL}/auth/me`, {
              headers: {
                Authorization: `Bearer ${accessToken}`,
              },
            });
            if (meRes.ok) {
              const mePayload = await meRes.json();
              permissions = mePayload.data?.permissions || [];
              isHost = !!mePayload.data?.tenant?.is_host;
              tenantName = mePayload.data?.tenant?.name || '';
            }
          } catch (e) {
            console.warn('Failed to hydrate permissions from /auth/me:', e);
          }

          return {
            id: authData.user.id,
            name: authData.user.name,
            email: authData.user.email,
            role: authData.user.role,
            permissions,
            tenantId: authData.tenant_id,
            tenantName,
            isHost,
            accessToken,
            refreshToken,
          } as UserWithPermissions;
        } catch (error: any) {
          // If offline / dev fallback is requested when backend is unreachable
          if (
            credentials.email === 'demo@example.com' &&
            credentials.password === 'demo12345'
          ) {
            return {
              id: 'user-demo-1',
              name: 'Demo Admin',
              email: credentials.email,
              role: 'owner',
              permissions: ['*'],
              tenantId: 'default',
              tenantName: 'Platform Host',
              isHost: true,
            } as UserWithPermissions;
          }

          throw error;
        }
      },
    }),
  ],
  pages: {
    signIn: '/login',
    newUser: '/register',
    error: '/login',
  },
  session: {
    strategy: 'jwt',
    maxAge: 7 * 24 * 60 * 60, // 7 days
  },
  callbacks: {
    async jwt({ token, user }): Promise<JWT> {
      if (user) {
        const customUser = user as UserWithPermissions;
        token.id = customUser.id;
        token.role = customUser.role || 'member';
        token.permissions = customUser.permissions || [];
        token.tenantId = customUser.tenantId || 'default';
        token.tenantName = customUser.tenantName || '';
        token.isHost = customUser.isHost ?? false;
        token.accessToken = customUser.accessToken;
        token.refreshToken = customUser.refreshToken;
      }
      return token;
    },
    async session({ session, token }): Promise<Session> {
      if (session.user) {
        (session.user as any).id = token.id as string;
        (session.user as any).role = (token.role as string) || 'member';
        (session.user as any).permissions = (token.permissions as string[]) || [];
        (session.user as any).tenantId = (token.tenantId as string) || 'default';
        (session.user as any).tenantName = (token.tenantName as string) || '';
        (session.user as any).isHost = !!token.isHost;
      }
      (session as any).accessToken = token.accessToken as string;
      return session;
    },
  },
  secret: process.env.NEXTAUTH_SECRET || 'development-secret-key-change-in-production',
};
