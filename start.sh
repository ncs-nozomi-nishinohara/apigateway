#!/usr/bin/env sh

if [ -e "/app/config/config.yaml" ]; then
    /app/app
fi

exec "$@"