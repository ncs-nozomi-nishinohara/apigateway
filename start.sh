#!/usr/bin/env sh

if [ -e "/app/config.yaml" ]; then
    /app/app
fi

exec "$@"