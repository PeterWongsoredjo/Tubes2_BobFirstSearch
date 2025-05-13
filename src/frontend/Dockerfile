# Install dependencies and build
FROM node:20-alpine AS builder
WORKDIR /app

# Copy package files first for better caching
COPY package*.json ./
RUN npm install -g npm@latest
RUN npm install

# Copy the rest of the application code
# (The .dockerignore should NOT exclude the source code!)
COPY . .

# Build the Next.js application
RUN npm run build

# Production image
FROM node:20-alpine
WORKDIR /app

# Copy package files for production dependencies
COPY --from=builder /app/package*.json ./

# Copy the complete .next directory (including build ID files)
COPY --from=builder /app/.next ./.next

# Copy public assets
COPY --from=builder /app/public ./public

# Copy necessary config files
COPY --from=builder /app/next.config.ts ./next.config.ts

# Install only production dependencies
RUN npm install --production

EXPOSE 3000
CMD ["npm", "start"]