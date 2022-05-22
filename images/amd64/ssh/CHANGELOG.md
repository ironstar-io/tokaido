[0.0.9]

- Improved dockerfile for faster builds
- Make sure drush for 'tok' user is executable

[0.0.8]

- Include some additional env settings for production environments
- Automatically copy any SSH host keys from /tokaido/config/host_ssh_keys
- Improved Dockerfile layout to enable faster builds in future

[0.0.7]

- Remove 'tok tip' about running tok watch, it's not needed anymore. 
- Update base PHP image to run 7.1.20
- Set non-login shells to automatically cd to drupal site root
- Update included Drush to 9.4.0

[0.0.6]

- Install Redis CLI 4.0.11
