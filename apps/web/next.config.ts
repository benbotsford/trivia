import type { NextConfig } from 'next'

const nextConfig: NextConfig = {
  // Produces a self-contained output in .next/standalone — required for the
  // lean Docker image. Only the files Next.js needs to run are included;
  // node_modules are not copied into the final container layer.
  output: 'standalone',

  // @auth0/nextjs-auth0 uses bare-string subpath exports ("./client", "./edge")
  // which Turbopack (Next.js 15 default bundler) doesn't resolve automatically.
  // transpilePackages forces Next.js to run the package through its bundler
  // pipeline, fixing the "Module not found" error for the /client and /edge imports.
  transpilePackages: ['@auth0/nextjs-auth0'],

  // Proxy /api/* to the Go backend for browser-side requests in local dev.
  // In production, ingress handles this routing instead.
  //
  // /api/auth/* is explicitly excluded — those routes are handled by the
  // Next.js Auth0 route handler (app/api/auth/[auth0]/route.ts) and must
  // never be forwarded to the Go API or the OAuth callback will break.
  async rewrites() {
    return [
      {
        source: '/api/((?!auth/).*)',
        destination: `${process.env.API_URL ?? 'http://localhost:8080'}/:path*`,
      },
    ]
  },
}

export default nextConfig
