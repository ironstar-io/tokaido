#!/bin/bash

if [[ -f /app/site/.env ]]; then
    printf "Importing environment variables from /app/site/.env\n"
    set -o allexport
    source /app/site/.env || true
    set +o allexport
fi

/usr/local/php/sbin/php-fpm -F -c /app/config/php/runtime/php.ini --fpm-config /app/config/php/runtime/php-fpm.conf
