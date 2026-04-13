#!/usr/bin/env bash
# ctm-batch.sh - Batch execution wrapper

CONFIG_FILE="${1:-batch.json}"
ITERATIONS="${2:-}"

if [[ -n "$ITERATIONS" ]]; then
    sed -i "s/\"iterations\": [0-9]*/\"iterations\": $ITERATIONS/" "$CONFIG_FILE"
fi

exec ctm batch --config "$CONFIG_FILE"
