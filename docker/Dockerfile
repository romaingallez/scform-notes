# Build frontend assets
FROM node:20-alpine AS frontend-builder
WORKDIR /app

# Install pnpm
RUN npm install -g pnpm@latest

# Copy frontend-related files
COPY package.json pnpm-lock.yaml ./
COPY postcss.config.js tailwind.config.js ./
COPY assets/ ./assets/
COPY views/ ./views/

# Install dependencies and build frontend
RUN pnpm install --frozen-lockfile
# Create dist directory for Tailwind build
RUN mkdir -p ./assets/dist
RUN pnpm run build

# Build Go backend
FROM golang:1.24-alpine AS backend-builder
WORKDIR /app

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Copy Go files
COPY go.mod go.sum ./
RUN go mod download

COPY main.go ./
COPY internals/ ./internals/
COPY views/ ./views/

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Final stage
FROM alpine:3.19
WORKDIR /app

# Install runtime dependencies if needed
RUN apk add --no-cache ca-certificates

# Copy assets source files (dist will be created by InitAssets)
COPY --from=frontend-builder /app/assets/src /app/assets/src
COPY --from=frontend-builder /app/views /app/views

# Copy the Go binary
COPY --from=backend-builder /app/main /app/

# Expose the port your application runs on
EXPOSE 3000

# Run the application
CMD ["./main"]
