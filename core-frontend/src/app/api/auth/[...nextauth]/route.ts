import NextAuth from 'next-auth';
import { defaultAuthOptions } from '@/core/auth/auth-options';

export const authOptions = defaultAuthOptions;

const handler = NextAuth(authOptions);

export { handler as GET, handler as POST };
