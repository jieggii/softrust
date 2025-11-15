# 1. Install dependencies
FROM node:20-alpine AS deps
WORKDIR /app

# Install OS deps (optional but recommended for sharp and other native mods)
RUN apk add --no-cache libc6-compat

COPY package.json package-lock.json* yarn.lock* pnpm-lock.yaml* ./
# Automatically picks the right package manager
RUN \
  if [ -f yarn.lock ]; then yarn install --frozen-lockfile; \
  elif [ -f pnpm-lock.yaml ]; then npm install -g pnpm && pnpm install --frozen-lockfile; \
  else npm ci; \
  fi


# 2. Build stage
FROM node:20-alpine AS builder
WORKDIR /app

COPY --from=deps /app/node_modules ./node_modules
COPY . .

RUN npm run build


# 3. Production runtime
FROM node:20-alpine AS runner
WORKDIR /app

ENV NODE_ENV=production

# Needed for Next.js
RUN addgroup --system nextjs && adduser --system nextjs -G nextjs

COPY --from=builder /app/public ./public
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/package.json ./package.json

USER nextjs

EXPOSE 3000
CMD ["npm", "start"]