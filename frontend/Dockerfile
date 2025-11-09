FROM node:20-alpine
WORKDIR /app

# Enable corepack/pnpm
RUN corepack enable

# Copy only package manager files first (cache-friendly)
COPY frontend/package.json frontend/pnpm-lock.yaml ./
RUN pnpm install

# Now copy the source
COPY frontend ./

EXPOSE 5173

# Auto-install only if vite binary is missing (volume vacío)
CMD ["sh","-lc","[ -x node_modules/.bin/vite ] || pnpm install && pnpm dev --host 0.0.0.0"]
