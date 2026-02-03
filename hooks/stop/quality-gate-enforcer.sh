#!/bin/bash
# Quality Gate Enforcer - Stop Hook
# Runs quality checks before allowing session to end

set -e

cd "$(dirname "$0")/../.."

echo "=== Quality Gate Check ==="

# Gate 1: Build
echo "Checking build..."
cargo build --release 2>/dev/null || { echo "BUILD FAILED"; exit 1; }

# Gate 2: Tests
echo "Running tests..."
cargo test 2>/dev/null || { echo "TESTS FAILED"; exit 1; }

# Gate 3: Clippy
echo "Running clippy..."
cargo clippy -- -D warnings 2>/dev/null || { echo "CLIPPY FAILED"; exit 1; }

# Gate 4: Format
echo "Checking format..."
cargo fmt --check 2>/dev/null || { echo "FORMAT CHECK FAILED"; exit 1; }

echo "=== All Quality Gates Passed ==="
exit 0
