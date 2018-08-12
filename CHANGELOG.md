## Unreleased

- Remove version and customcompose as persistent global flags
- Audit of command short and long descriptions
- Fix missing Drupal version detection
- Tokaido will now prompt for a Drupal version and path if it can't detect one automatically
- Add the /private and sites/*/files folders to gitignore automatically

- Add Adminer service (enable with: `tok config-set service adminer enabled true`) 
- Add Redis service (enable with: `tok config-set service redis enabled true`)
- Add MailHog support (enable with: `tok config-set service mailhog enabled true`)
