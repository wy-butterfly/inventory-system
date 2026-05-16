/** @type {import('next').NextConfig} */
const nextConfig = {
  eslint: {
    ignoreDuringBuilds: true,
  },
  typescript: {
    ignoreBuildErrors: true,
  },
};

// 只在本地开发环境启用 API rewrite，避免 Vercel 生产环境转发到 localhost
if (process.env.NODE_ENV === 'development') {
  nextConfig.rewrites = async () => [
    {
      source: '/api/:path*',
      destination: 'http://localhost:8080/api/:path*',
    },
  ];
}

export default nextConfig;
