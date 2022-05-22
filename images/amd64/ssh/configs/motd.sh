#!/usr/bin/env bash
drupal_root=${DRUPAL_ROOT:-docroot}
if [[ $(tput colors) == 256 ]]; then
  green=$(tput setaf 42)
  blue=$(tput setaf 32)
  red=$(tput setaf 124)
  cyan=$(tput setaf 6)
  reset=$(tput sgr0)
else
  green=$(tput setaf 2)
  blue=$(tput setaf 4)
  red=$(tput setaf 1)
  cyan=$(tput setaf 6)
  reset=$(tput sgr0)
fi

tok_provider=${TOK_PROVIDER:-}
php_version=$(php -v 2> /dev/null| head -n 1 | awk -F ' ' '{print $2}' | awk -F '-' '{print $1}')
composer_version=$(composer --version 2> /dev/null | awk -F ' ' '{print $3}')
drush_version=$(cd /tokaido/site/${drupal_root} 2> /dev/null && drush version 2> /dev/null | awk -F ' ' '{print $4}')
database_status=$(cd /tokaido/site/${drupal_root} 2> /dev/null && drush status 2> /dev/null | grep "Database" | grep "Connected" | awk -F ' ' '{print $3}')

echo "${blue}"
cat <<"EOF"
                                          __   __
                                          /'   `\
                                         Y.     .Y
                               _______    \`. .'/
                ,-------------'======="--""""-""""---.
          __,=+'-------------------------------------|p
       .-/__|_]_]  :"/:""""""""""""""""""""""""""""""|'
    ,-'__________[];/_;______________________________|
  ,".../_|___________________________________________|
 (_>        ,-------.                     ,-------.  |
  `-._____.'(_)`='(_)\_7___7___7___7__7_.'(_)`='(_)\_/ hjw

        __________  __ __ ___    ________  ____
       /_  __/ __ \/ //_//   |  /  _/ __ \/ __ \
        / / / / / / ,<  / /| |  / // / / / / / /
       / / / /_/ / /| |/ ___ |_/ // /_/ / /_/ /
      /_/  \____/_/ |_/_/  |_/___/_____/\____/

EOF
echo "${reset}"
echo "      PHP Version      : $php_version"
echo "      Composer Version : $composer_version"
echo "      Drush Version    : $drush_version"

if [[ "$database_status" == "Connected" ]]; then
  echo "      Database Status  : ${green}OK${reset}"
else
  echo "      Database Status  : ${red}Failed${reset}"
fi

if ! [[ -z "$tok_provider" ]]; then
  cat <<"EOF"

UNAUTHORISED ACCESS TO THIS SYSTEM IS PROHIBITED
You must have explicit permission to access or configure this device.

Unauthorised attempts and actions to access or use this system may result in
civil and/or criminal penalties. All activities performed on this device are
logged and monitored.

---

The Tokaido Drush container is ephemeral. Each time you deploy your
application, anything you save here will be lost forever.

Only the Drupal Public and Private Files directories are writeable.

EOF

else
  cat <<"EOF"

In Tokaido, the /tokaido/site folder is synchronised between this environment
and your local system. Content not saved in the /tokaido/site may be lost
each time the Tokaido environment is restarted.

Be sure to check out https://docs.tokaido.io for help, or come and talk to
us in the #Tokaido channel in the official Drupal Slack: https://www.drupal.org/slack

EOF
fi
