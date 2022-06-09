#!/usr/bin/env bash
set -euxo pipefail

drupal_root=${DRUPAL_ROOT:-web}


# Invoke the environment variables into this sshd process
while read -r line; do
  export "${line?}"
done < /app/config/.env

echo "Adding your local SSH key to the 'tok' user"
username="tok"
cp /app/site/.tok/local/ssh_key.pub /home/"$username"/.ssh/authorized_keys

# Set up environment variables for the user
echo "Setting up environment variables for $username"
echo "PATH=$PATH:/usr/local/bin" > /home/"$username"/.ssh/environment
cat /app/config/.env >> /home/"$username"/.ssh/environment
echo "APP_ENV=${APP_ENV:-unknown}" >> /home/"$username"/.ssh/environment
echo "PROJECT_NAME=${PROJECT_NAME:-}" >> /home/"$username"/.ssh/environment
echo "DRUPAL_ROOT=${drupal_root}" >> /home/"$username"/.ssh/environment
echo "VARNISH_PURGE_KEY=${VARNISH_PURGE_KEY:-}" >> /home/"$username"/.ssh/environment
chmod 600 /home/"$username"/.ssh/environment
chmod 600 /home/"$username"/.ssh/authorized_keys
chown "$username":root /home/"$username"/.ssh -R

# If we're running in a Tokaido production environment, we'll create multiple users
# and also set up some additional configuration that they'll need
# Give users read access to the environment file
chown tok:web /tokaido/config/.env
chmod 0750 /tokaido/config/.env

# Start SSH server
/usr/sbin/sshd -D -e
