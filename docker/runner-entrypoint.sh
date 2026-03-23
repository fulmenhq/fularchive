#!/bin/sh
set -eu

if [ "$#" -gt 0 ]; then
    exec "$@"
fi

crontab_path="${SUPERCRONIC_CRONTAB:-/etc/refbolt/crontab}"

if [ ! -f "$crontab_path" ]; then
    echo "missing crontab file: $crontab_path" >&2
    exit 1
fi

if [ ! -s "$crontab_path" ]; then
    echo "crontab file is empty: $crontab_path" >&2
    exit 1
fi

echo "starting supercronic with $crontab_path"
exec /usr/local/bin/supercronic "$crontab_path"
