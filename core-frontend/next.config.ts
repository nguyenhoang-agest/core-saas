import type { NextConfig } from 'next';

const nextConfig: NextConfig = {
  reactStrictMode: true,
  transpilePackages: ['antd', '@ant-design/cssinjs', '@ant-design/nextjs-registry'],
};

export default nextConfig;
