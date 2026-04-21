# Stage 1: install dependencies
FROM node:20-alpine AS deps

WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci

# Stage 2: build
FROM node:20-alpine AS builder

WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY . .

# Build args for env vars needed at build time (none currently use NEXT_PUBLIC_)
# Runtime env vars (API_URL, DEV_AUTH_TOKEN) are injected at container start.
RUN npm run build

# Stage 3: lean runtime using Next.js standalone output
FROM node:20-alpine AS runner

WORKDIR /app
ENV NODE_ENV=production

# Copy only what standalone output needs
COPY --from=builder /app/.next/standalone ./
COPY --from=builder /app/.next/static ./.next/static

EXPOSE 3000

ENV PORT=3000
ENV HOSTNAME=0.0.0.0

CMD ["node", "server.js"]
