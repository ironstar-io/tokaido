package nightwatch

import "github.com/ironstar-io/tokaido/system/fs"

func generateEnvFile(path string) {
	fs.TouchByteArray(path, buildEnvFileContents())
}

func buildEnvFileContents() []byte {
	return []byte(`
#############################
# General Test Environment #
#############################
# This is the URL that Drupal can be accessed by. You don't need an installed
# site here, just make sure you can at least access the installer screen. If you
# don't already have one running, e.g. Apache, you can use PHP's built-in web
# server by running the following command in your Drupal root folder:
# php -S localhost:8888 .ht.router.php
DRUPAL_TEST_BASE_URL=https://nginx:8443

# Tests need to be executed with a user in the same group as the web server
# user.
#DRUPAL_TEST_WEBSERVER_USER=www-data

# By default we use sqlite as database. Use
# mysql://username:password@localhost/databasename#table_prefix for mysql.
#DRUPAL_TEST_DB_URL=sqlite://localhost/sites/default/files/db.sqlite
DRUPAL_TEST_DB_URL=mysql://tokaido:tokaido@mysql/tests

#############
# Webdriver #
#############

# If Chromedriver is running as a service elsewhere, set it here.
# When using DRUPAL_TEST_CHROMEDRIVER_AUTOSTART leave this at the default settings.
DRUPAL_TEST_WEBDRIVER_HOSTNAME=chromedriver
DRUPAL_TEST_WEBDRIVER_PORT=9515

# If using Selenium, override the path prefix here.
# See http://nightwatchjs.org/gettingstarted#browser-drivers-setup
#DRUPAL_TEST_WEBDRIVER_PATH_PREFIX=/wd/hub

################
# Chromedriver #
################

# Automatically start chromedriver for local development. Set to false when you
# use your own webdriver or chromedriver setup.
# Also set it to false when you use a different browser for testing.
DRUPAL_TEST_CHROMEDRIVER_AUTOSTART=false

# A list of arguments to pass to Chrome, separated by spaces
# e.g. '--disable-gpu --headless --no-sandbox'.
DRUPAL_TEST_WEBDRIVER_CHROME_ARGS="--disable-gpu --headless --no-sandbox"

##############
# Nightwatch #
##############

# Nightwatch generates output files. Use this to specify the location where these
# files need to be stored. The default location is ignored by git, if you modify
# the location you will probably want to add this location to your .gitignore.
DRUPAL_NIGHTWATCH_OUTPUT=reports/nightwatch

# The path that Nightwatch searches for assumes the same directory structure as
# when you download Drupal core. If you have Drupal installed into a docroot
# folder, you can use the following folder structure to add integration tests
# for your project, outside of tests specifically for custom modules/themes/profiles.
#
# .
# ├── docroot
# │   ├── core
# ├── tests
# │   ├── Nightwatch
# │   │   ├── Tests
# │   │   │   ├── myTest.js
#
# and then set DRUPAL_NIGHTWATCH_SEARCH_DIRECTORY=../
#
#DRUPAL_NIGHTWATCH_SEARCH_DIRECTORY=

# Filter directories to look for tests. This uses minimatch syntax.
# Separate folders with a comma.
DRUPAL_NIGHTWATCH_IGNORE_DIRECTORIES=node_modules,vendor,.*,sites/*/files,sites/*/private,sites/simpletest

`)
}
