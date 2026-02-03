# Build stage
FROM docker.io/library/rust:1-slim AS builder

WORKDIR /app

# Install build dependencies and nightly toolchain
RUN apt-get update && apt-get install -y \
    pkg-config \
    libssl-dev \
    curl \
    && rm -rf /var/lib/apt/lists/* \
    && rustup default nightly

# Copy manifests and toolchain config
COPY Cargo.toml Cargo.lock* rust-toolchain.toml ./

# Copy source and config
COPY src ./src
COPY settings ./settings
COPY migrations ./migrations

# Build release
RUN cargo build --release

# Runtime stage
FROM docker.io/library/debian:bookworm-slim

# Install runtime dependencies including curl for healthcheck
RUN apt-get update && apt-get install -y \
    ca-certificates \
    curl \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy binaries from builder
COPY --from=builder /app/target/release/hailango /usr/local/bin/
COPY --from=builder /app/target/release/manage /usr/local/bin/

# Copy settings for runtime
COPY --from=builder /app/settings ./settings

# Create non-root user
RUN useradd -r -s /bin/false hailango && chown -R hailango:hailango /app
USER hailango

# Expose port
EXPOSE 8080

# Environment defaults
ENV APP_ENV=production
ENV RUST_LOG=info

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

# Run the application
CMD ["hailango"]
