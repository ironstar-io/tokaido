## 1.0.0

- Remove version and customcompose as persistent global flags
- Audit of command short and long descriptions
- Fix missing Drupal version detection
- Tokaido will now prompt for a Drupal version and path if it can't detect one automatically
- Add the /private and sites/\*/files folders to gitignore automatically
- Add Adminer service (enable with: `tok config-set service adminer enabled true`)
- Add Redis service (enable with: `tok config-set service redis enabled true`)
- Add MailHog support (enable with: `tok config-set service mailhog enabled true`)
- 'Open' can now be used to open secondary services like mailhog, adminer (eg `tok open adminer`)
- 'ports' now dynamically returns ports only for services which are enabled
- Add dependency checks for Linux
- MySQL max_allowed_packet is now set to 64MB by default. 
- Fixed a bug with Compose file used 'cmd' instead of 'command' parameters. 
- Add the Tokaido Proxy service and HTTPS certificate support

