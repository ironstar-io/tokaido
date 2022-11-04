#!/usr/bin/env bash
set -euxo pipefail

drupal_root=${DRUPAL_ROOT:-web}

if [[ -f /app/site/.env ]]; then
    cat /app/site/.env >> /home/app/.ssh/environment || true
fi

ep /app/config/php/php.ini
ep /app/config/php/www.conf

echo "Adding your local SSH key to the 'tok' user"
cp /app/site/.tok/local/ssh_key.pub /home/app/.ssh/authorized_keys

# Set up environment variables for the user
echo "PATH=$PATH:/usr/local/bin" > /home/app/.ssh/environment
echo "APP_ENV=${APP_ENV:-unknown}" >> /home/app/.ssh/environment
echo "PROJECT_NAME=${PROJECT_NAME:-}" >> /home/app/.ssh/environment
echo "DRUPAL_ROOT=${drupal_root}" >> /home/app/.ssh/environment

chown app:app /home/app -R
find /home/app -type d -print0 | xargs -P0 -0 chmod 2770 -f
find /home/app -wholename \*installation/bin/node -print0 | xargs -P0 -0 chmod u+x
find /home/app -name npm-cli.js -print0 | xargs -P0 -0 chmod u+x
chown app:root /home/app/.ssh -R
chmod 600 /home/app/.ssh/environment
chmod 600 /home/app/.ssh/authorized_keys

# Start SSH server
/usr/sbin/sshd -D -e
