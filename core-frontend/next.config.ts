import type { NextConfig } from 'next';

const nextConfig: NextConfig = {
  reactStrictMode: true,
  transpilePackages: ['antd', '@ant-design/cssinjs', '@ant-design/nextjs-registry'],
  async redirects() {
    return [
      {
        source: '/iam',
        destination: '/users',
        permanent: true,
      },
    ];
  },
  async rewrites() {
    const backendUrl = process.env.BACKEND_INTERNAL_URL || 'http://localhost:8080';
    return {
      beforeFiles: [],
      afterFiles: [],
      fallback: [
        {
          source: '/api/:path*',
          destination: `${backendUrl}/:path*`,
        },
      ],
    };
  },
};

export default nextConfig;
