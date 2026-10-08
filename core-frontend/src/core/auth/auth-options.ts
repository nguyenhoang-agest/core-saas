import type { AuthOptions, Session, User } from 'next-auth';
import type { JWT } from 'next-auth/jwt';
import CredentialsProvider from 'next-auth/providers/credentials';

export interface UserWithPermissions extends User {
  id: string;
  permissions?: string[];
  tenantId?: string;
}

export interface CustomSession extends Session {
  user: {
    id?: string;
    name?: string | null;
    email?: string | null;
    image?: string | null;
    permissions?: string[];
    tenantId?: string;
  };
}

/**
 * Reusable NextAuth configuration options.
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
        if (!credentials?.email) return null;

        // In a real implementation, call domain IAM / Auth service
        // For development/demo, returns a mock user with permissions
        return {
          id: 'user-demo-1',
          name: 'Demo Admin',
          email: credentials.email,
          permissions: ['*'], // wildcard for development admin
          tenantId: 'default',
        } as UserWithPermissions;
      },
    }),
  ],
  session: {
    strategy: 'jwt',
  },
  callbacks: {
    async jwt({ token, user }): Promise<JWT> {
      if (user) {
        const customUser = user as UserWithPermissions;
        token.id = customUser.id;
        token.permissions = customUser.permissions || [];
        token.tenantId = customUser.tenantId || 'default';
      }
      return token;
    },
    async session({ session, token }): Promise<Session> {
      if (session.user) {
        (session.user as any).id = token.id as string;
        (session.user as any).permissions = (token.permissions as string[]) || [];
        (session.user as any).tenantId = (token.tenantId as string) || 'default';
      }
      return session;
    },
  },
  secret: process.env.NEXTAUTH_SECRET || 'development-secret-key-change-in-production',
};
