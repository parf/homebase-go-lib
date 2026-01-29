#!/bin/bash

# jsonl2parquet - Convert JSONL to Parquet format
# Supports compressed inputs: .gz, .zst, .lz4, .br, .xz

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TOOL_NAME="jsonl2parquet"
TOOL_PATH="$SCRIPT_DIR/$TOOL_NAME"

# Build tool if not exists or source is newer
if [ ! -f "$TOOL_PATH" ] || [ "$SCRIPT_DIR/$TOOL_NAME.go" -nt "$TOOL_PATH" ]; then
    echo "Building $TOOL_NAME..."
    (cd "$SCRIPT_DIR/.." && go build -o "cmd/$TOOL_NAME" "./cmd/$TOOL_NAME.go")
fi

# Run the tool
exec "$TOOL_PATH" "$@"
