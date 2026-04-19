import type { NextConfig } from 'next'

const nextConfig: NextConfig = {
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
