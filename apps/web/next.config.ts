import type { NextConfig } from 'next'

const nextConfig: NextConfig = {
  // Produces a self-contained output in .next/standalone — required for the
  // lean Docker image. Only the files Next.js needs to run are included;
  // node_modules are not copied into the final container layer.
  output: 'standalone',

  // The Go API runs on :8080 locally. All /api/* requests from the frontend
  // are proxied here so the browser never needs to know the API's address.
  // In production, the ingress-nginx controller handles this routing instead.
  async rewrites() {
    return [
      {
        source: '/api/:path*',
        destination: `${process.env.API_URL ?? 'http://localhost:8080'}/:path*`,
      },
    ]
  },
}

export default nextConfig
