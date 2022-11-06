#!/usr/bin/env bash
set -euxo pipefail

# This file performs runtime preparation for the SSH image. It is called by start.sh, but should
# also by called in any CI system before performing any actions. Without it, things like PHP
# variables and FNM won't be configured correctly and will not work

echo "Preparing PHP config files"
ep /app/config/php/php.ini
ep /app/config/php/www.conf

echo "Configuring home directory permissions" # mostly for FNM
chown app:app /home/app -R
find /home/app -type d -print0 | xargs -P0 -0 chmod 2770 -f
find /home/app -wholename \*installation/bin/node -print0 | xargs -P0 -0 chmod u+x
find /home/app -name npm-cli.js -print0 | xargs -P0 -0 chmod u+x
