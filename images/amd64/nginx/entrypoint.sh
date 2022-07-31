#!/usr/bin/env bash
set -eo pipefail

################################################################################
#
# Setting Default Variable Values
# Tokaido enables us to set variable values at three levels. These levels
# are as follows, listed from highest-priority to lowest-priority
# 1 - As Environment Variable exposed to the PHP container
# 2 - As Config Variable in .tok/config.yml
# 3 - As default variable defined here
#
# For example, if the variable FASTCGI_BUFFERS is defined in .tok/config.yml
# then any value set as an environment variable (FASTCGI_BUFFERS) won't be used
#
################################################################################

# Colours
RED='\e[31m'
BLUE='\e[34m'
GREEN='\e[32m'
YELLOW='\e[33m'
PURPLE='\e[35m'
CYAN='\e[36m'
NC='\033[0m' # No Color

printf "${GREEN}NGINX container is starting...${NC}\n"

# resolved holds all of our final values to be applied
declare -A settings=()

# defaults holds all of our default values
declare -A defaultSettings

# Default Values
defaultSettings[WORKER_CONNECTIONS]="1024"
defaultSettings[TYPES_HASH_MAX_SIZE]="2048"
defaultSettings[CLIENT_MAX_BODY_SIZE]="1024m"
defaultSettings[KEEPALIVE_TIMEOUT]="300"
defaultSettings[FASTCGI_READ_TIMEOUT]="300"
defaultSettings[FASTCGI_BUFFERS]="16 16k"
defaultSettings[FASTCGI_BUFFER_SIZE]="32k"
defaultSettings[DRUPAL_ROOT]="docroot"
defaultSettings[ALLOWED_HOSTS]="localhost"

# Set default "null" Tokaido values
declare -A tokaidoSettings
tokaidoSettings[WORKER_CONNECTIONS]="null"
tokaidoSettings[TYPES_HASH_MAX_SIZE]="null"
tokaidoSettings[CLIENT_MAX_BODY_SIZE]="null"
tokaidoSettings[KEEPALIVE_TIMEOUT]="null"
tokaidoSettings[FASTCGI_READ_TIMEOUT]="null"
tokaidoSettings[FASTCGI_BUFFERS]="null"
tokaidoSettings[FASTCGI_BUFFER_SIZE]="null"
tokaidoSettings[DRUPAL_ROOT]="null"

# If there is a tokaido config, look up any values
if [ -f /app/site/.tok/config.yml ]; then
    tokaidoSettings[WORKER_CONNECTIONS]="$(yq '.nginx.workerconnections' /app/site/.tok/config.yml)"
    tokaidoSettings[TYPES_HASH_MAX_SIZE]="$(yq '.nginx.hashmaxsize' /app/site/.tok/config.yml)"
    tokaidoSettings[CLIENT_MAX_BODY_SIZE]="$(yq '.nginx.clientmaxbodysize' /app/site/.tok/config.yml)"
    tokaidoSettings[KEEPALIVE_TIMEOUT]="$(yq '.nginx.keepalivetimeout' /app/site/.tok/config.yml)"
    tokaidoSettings[FASTCGI_READ_TIMEOUT]="$(yq '.nginx.fastcgireadtimeout' /app/site/.tok/config.yml)"
    tokaidoSettings[FASTCGI_BUFFERS]="$(yq '.nginx.fastcgibuffers' /app/site/.tok/config.yml)"
    tokaidoSettings[FASTCGI_BUFFER_SIZE]="$(yq '.nginx.fastcgibuffersize' /app/site/.tok/config.yml)"
    tokaidoSettings[DRUPAL_ROOT]="$(yq '.drupal.path' /app/site/.tok/config.yml)"
fi

printf "${BLUE}NGINX will run with the following configuration values and sources:${NC}\n"
for i in "${!defaultSettings[@]}"
do
    if [ -n "${!i}" ]; then
        # An ENV var exists for this setting, so we'll use it
        settings["$i"]="${!i}"
        printf "  ${CYAN}$i${NC}"
        printf "\033[50D\033[43C :: ${YELLOW}Use Env value${NC}"
        printf "\033[50D\033[69C :: ${BLUE}$settings[${!i}]${NC}\n"
        continue
    elif [[ ${tokaidoSettings[$i]} != "null" ]] && [[ ! -z ${tokaidoSettings[$i]} ]]; then
        # No ENV var exists - check if a Tokaido value exists
        settings["$i"]="${tokaidoSettings[$i]}"
        printf "  ${CYAN}$i${NC}"
        printf "\033[50D\033[43C :: ${GREEN}Use Tokaido value${NC}"
        printf "\033[50D\033[65C :: ${BLUE}[${settings[$i]}]${NC}\n"
        continue
    fi

    # No env var or tokaido var exists, so we use the default
    settings["$i"]="${defaultSettings[$i]}"
    printf "  ${CYAN}$i${NC}"
    printf "\033[50D\033[43C :: ${PURPLE}Use Default value${NC}"
    printf "\033[50D\033[65C :: ${BLUE}[${settings[$i]}]${NC}\n"
done

# Strip any forward-slashes out of our resolve drupal root, just in case
settings[DRUPAL_ROOT]=$(echo ${settings[DRUPAL_ROOT]} | sed -e 's/\///g')

# FPM_HOSTNAME is a special value that can only be provided
# as environment variables, not via the .tok/config.yml file.
settings[FPM_HOSTNAME]=${FPM_HOSTNAME:-fpm}

# ALLOWED_HOSTS can only be provided via a base64 encoded environment variable
# this value can be anything matching the nginx server_name directive, including wildcards
if [ ! -z "$ALLOWED_HOSTS" ]; then
    settings[ALLOWED_HOSTS]=$(echo "$ALLOWED_HOSTS" | base64 --decode)
else
    settings[ALLOWED_HOSTS]="${defaultSettings[ALLOWED_HOSTS]}"
fi

# BLOCK_UNKNOWN_HOSTS if "true" will add a default server_name block that prohibits unconfigured
# ALLOWED_HOSTS from being served content
if [ ! -z "${BLOCK_UNKNOWN_HOSTS}" ]; then
    settings[BLOCK_UNKNOWN_HOSTS]="server {\n    listen 8082 default_server;\n    server_name _;\n    return 400;\n}"
else
    settings[BLOCK_UNKNOWN_HOSTS]=""
fi


################################################################################
#
# Setting Config Files Paths
# Tokaido can supply custom Nginx config files to completely override the
# config files included in this Docker image.
#
# For example, if the file .tok/nginx/redirects.conf exists, it will be used
# for all redirect config instead of the default.
#
# These config files can be used in conjunction with the above config values
# as well, so you can both use your own config file and per-environment
# overrides from environment variables, for example.
#
################################################################################

NGINX_CONFIG_FILE="nginx.conf"
HOST_CONFIG_FILE="host.conf"
MIMETYPES_CONFIG_FILE="mimetypes.conf"
REDIRECTS_CONFIG_FILE="redirects.conf"
ADDITIONAL_CONFIG_FILE="additional.conf"
SECURITY_CONFIG_FILE="security.conf"

DEFAULT_CONFIG_PATH="/app/config/nginx"
CUSTOM_CONFIG_PATH="/app/site/.tok/nginx"

declare -A configFiles
configFiles[NGINX_CONFIG]="$DEFAULT_CONFIG_PATH/$NGINX_CONFIG_FILE"
configFiles[HOST_CONFIG]="$DEFAULT_CONFIG_PATH/$HOST_CONFIG_FILE"
configFiles[MIMETYPES_CONFIG]="$DEFAULT_CONFIG_PATH/$MIMETYPES_CONFIG_FILE"

printf "${BLUE}Discovering any custom NGINX configuration files...${NC}\n"
if [ -f "${CUSTOM_CONFIG_PATH}/${NGINX_CONFIG_FILE}" ]; then
    configFiles[NGINX_CONFIG]="${CUSTOM_CONFIG_PATH}/${NGINX_CONFIG_FILE}"
    printf "${YELLOW}Custom config file ${configFiles[NGINX_CONFIG]} will be used${NC}\n"
fi

if [ -f "${CUSTOM_CONFIG_PATH}/${HOST_CONFIG_FILE}" ]; then
    configFiles[HOST_CONFIG]="${CUSTOM_CONFIG_PATH}/${HOST_CONFIG_FILE}"
    printf "${YELLOW}Custom config file ${configFiles[HOST_CONFIG]} will be used${NC}\n"
fi

if [ -f "${CUSTOM_CONFIG_PATH}/${MIMETYPES_CONFIG_FILE}" ]; then
    configFiles[MIMETYPES_CONFIG]="${CUSTOM_CONFIG_PATH}/${MIMETYPES_CONFIG_FILE}"
    printf "${YELLOW}Custom config file ${configFiles[MIMETYPES_CONFIG]} will be used${NC}\n"
fi

if [ -f "${CUSTOM_CONFIG_PATH}/${REDIRECTS_CONFIG_FILE}" ]; then
    configFiles[REDIRECTS_CONFIG]="${CUSTOM_CONFIG_PATH}/${REDIRECTS_CONFIG_FILE}"
    printf "${YELLOW}Custom config file ${configFiles[REDIRECTS_CONFIG]} will be used${NC}\n"
fi

if [ -f "${CUSTOM_CONFIG_PATH}/${ADDITIONAL_CONFIG_FILE}" ]; then
    configFiles[ADDITIONAL_CONFIG]="${CUSTOM_CONFIG_PATH}/${ADDITIONAL_CONFIG_FILE}"
    printf "${YELLOW}Custom config file ${configFiles[ADDITIONAL_CONFIG]} will be used${NC}\n"
fi

# Place all our defined variables into their respective config files
printf "Provisioning NGINX config files...\n"
sed -i "s/{{.WORKER_CONNECTIONS}}/${settings[WORKER_CONNECTIONS]}/g" "${configFiles[NGINX_CONFIG]}"
sed -i "s/{{.TYPES_HASH_MAX_SIZE}}/${settings[TYPES_HASH_MAX_SIZE]}/g" "${configFiles[NGINX_CONFIG]}"
sed -i "s/{{.CLIENT_MAX_BODY_SIZE}}/${settings[CLIENT_MAX_BODY_SIZE]}/g" "${configFiles[NGINX_CONFIG]}"
sed -i "s/{{.KEEPALIVE_TIMEOUT}}/${settings[KEEPALIVE_TIMEOUT]}/g" "${configFiles[NGINX_CONFIG]}"
sed -i "s/{{.FASTCGI_READ_TIMEOUT}}/${settings[FASTCGI_READ_TIMEOUT]}/g" "${configFiles[HOST_CONFIG]}"
sed -i "s/{{.FASTCGI_BUFFERS}}/${settings[FASTCGI_BUFFERS]}/g" "${configFiles[HOST_CONFIG]}"
sed -i "s/{{.FASTCGI_BUFFER_SIZE}}/${settings[FASTCGI_BUFFER_SIZE]}/g" "${configFiles[HOST_CONFIG]}"
sed -i "s/{{.DRUPAL_ROOT}}/${settings[DRUPAL_ROOT]}/g" "${configFiles[HOST_CONFIG]}"
sed -i "s/{{.FPM_HOSTNAME}}/${settings[FPM_HOSTNAME]}/g" "${configFiles[HOST_CONFIG]}"
sed -i "s/{{.STATUS_TOKEN}}/${settings[STATUS_TOKEN]}/g" "${configFiles[HOST_CONFIG]}"
sed -i "s/{{.ALLOWED_HOSTS}}/${settings[ALLOWED_HOSTS]}/g" "${configFiles[HOST_CONFIG]}"
sed -i "s/{{.BLOCK_UNKNOWN_HOSTS}}/${settings[BLOCK_UNKNOWN_HOSTS]}/g" "${configFiles[HOST_CONFIG]}"

sed -i "s/{{.HOST_CONFIG}}/${configFiles[HOST_CONFIG]//\//\\\/}/g" "${configFiles[NGINX_CONFIG]}"
sed -i "s/{{.MIMETYPES_CONFIG}}/${configFiles[MIMETYPES_CONFIG]//\//\\\/}/g" "${configFiles[NGINX_CONFIG]}"

printf "${GREEN}Starting NGINX...${NC}\n"
nginx -c "${configFiles[NGINX_CONFIG]}"
