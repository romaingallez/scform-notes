# Multi-stage Dockerfile for scform-notes

# Stage 1: Build frontend assets (Tailwind CSS)
FROM node:18-alpine AS frontend-builder

# Set working directory
WORKDIR /app

# Copy package files for frontend dependencies
COPY package.json pnpm-lock.yaml ./

# Install pnpm and frontend dependencies
RUN npm install -g pnpm && pnpm install

# Copy frontend source files
COPY tailwind.config.js postcss.config.js ./
COPY assets/src/ ./assets/src/
COPY views/ ./views/

# Create dist directory for Tailwind build
RUN mkdir -p ./assets/dist

# Build Tailwind CSS
RUN pnpm run build

# Stage 2: Build Go application
FROM golang:1.24-alpine AS go-builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Copy built frontend assets from previous stage
COPY --from=frontend-builder /app/assets/dist/ ./assets/dist/

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Stage 3: Final runtime image
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Copy the binary from builder stage
COPY --from=go-builder /app/main .

# Copy views and other static files
COPY --from=go-builder /app/views/ ./views/
COPY --from=go-builder /app/assets/dist/ ./assets/dist/
COPY --from=frontend-builder /app/assets/src/ ./assets/src/

# Change ownership to non-root user
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 3000

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:3000/ || exit 1

# Run the application
CMD ["./main"]
