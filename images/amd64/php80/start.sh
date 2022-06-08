#!/bin/bash
set -e -o pipefail -o errexit

if [[ -f /app/site/.env ]]; then
    printf "Importing environment variables from /app/site/.env\n"
    set -o allexport
    source /app/site/.env || true
    set +o allexport
fi

ep /app/config/php/php.ini
ep /app/config/php/www.conf

/usr/local/php/sbin/php-fpm -F -c /app/config/php/php.ini --fpm-config /app/config/php/www.conf
