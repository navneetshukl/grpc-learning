#!/bin/bash
# generate.sh — Compiles greet.proto into Go code
# Run this from the stage-06-first-grpc-server directory:
#   ./generate.sh

set -e

PROTO_DIR="proto"
GEN_DIR="gen"
PROTO_FILE="$PROTO_DIR/greet.proto"

echo "=== Step 1: Cleaning old generated files ==="
rm -rf "$GEN_DIR"/*.go 2>/dev/null || true

echo "=== Step 2: Generating Go code from proto ==="
protoc \
    --go_out="$GEN_DIR" \
    --go_opt=paths=source_relative \
    --go-grpc_out="$GEN_DIR" \
    --go-grpc_opt=paths=source_relative \
    -I "$PROTO_DIR" \
    "$PROTO_FILE"

echo "=== Step 3: Generated files ==="
ls -la "$GEN_DIR/"
echo ""
echo "Done! Generated files are in $GEN_DIR/"
