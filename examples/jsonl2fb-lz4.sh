#!/bin/bash
# Convert JSONL to FlatBuffer LZ4 format
# Usage: ./jsonl2fb-lz4.sh input.jsonl [output.fb.lz4]

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TOOL="$SCRIPT_DIR/jsonl2fb-lz4"

# Build the tool if it doesn't exist or source is newer
if [ ! -f "$TOOL" ] || [ "$SCRIPT_DIR/jsonl2fb-lz4.go" -nt "$TOOL" ]; then
    echo "Building jsonl2fb-lz4..."
    (cd "$SCRIPT_DIR" && go build -o jsonl2fb-lz4 jsonl2fb-lz4.go)
fi

# Run the tool
"$TOOL" "$@"
