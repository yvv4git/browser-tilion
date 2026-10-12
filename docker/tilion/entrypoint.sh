#!/usr/bin/env bash
# Thin wrapper over the official Fortress engine entrypoint.
# Optionally activates the Tilion license, then starts the engine.
set -e

if [ -n "${TILION_API_KEY:-}" ]; then
    tilion activate --token "$TILION_API_KEY"
fi

exec /usr/local/bin/docker-entrypoint.sh "$@"
