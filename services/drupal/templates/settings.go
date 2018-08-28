package drupaltmpl

// SettingsD7Append - (Append) docroot/sites/default/settings.php for Drupal 7
var SettingsD7Append = []byte(`/*
  * Generated by Tokaido
  */

  if (file_exists(DRUPAL_ROOT . '/sites/default/settings.tok.php')) {
    include DRUPAL_ROOT . '/sites/default/settings.tok.php';
  }

  /*
  * END Generated by Tokaido
  */

`)

// SettingsD7Tok - docroot/sites/default/settings.tok.php for Drupal 7
func SettingsD7Tok(projectName string) []byte {
	return []byte(`<?php

    /**
     * @file
     * Configuration file for Tokaido local dev environments. Add this to .gitignore
     *
     * Check out https://docs.tokaido.io for help managing your Tokaido environment
     *
     * Generated by Tokaido
     */
  
     $databases['default']['default'] = [
      'host' => 'mysql',
      'database' => 'tokaido',
      'username' => 'tokaido',
      'password' => 'tokaido',
      'port' => 3306,
      'driver' => 'mysql',
      'prefix' => '',
    ];
  
    $conf['file_private_path'] = '/tokaido/site/private/default';
    $conf['file_temporary_path'] = '/tmp';
  
    $base_url = 'https://` + projectName + `.tokaido.local:5154';
        
    /*
    * END Generated by Tokaido
    */
  
  `)
}

// SettingsD8Append - (Append) docroot/sites/default/settings.php for Drupal 8
var SettingsD8Append = []byte(`/*
* Generated by Tokaido
*/

if (file_exists($app_root . '/' . $site_path . '/settings.tok.php')) {
  include $app_root . '/' . $site_path . '/settings.tok.php';
}

/*
* END Generated by Tokaido
*/

`)

// SettingsD8Tok - docroot/sites/default/settings.tok.php for Drupal 8
var SettingsD8Tok = []byte(`<?php

/**
 * @file
 * Configuration file for Tokaido local dev environments. Add this to .gitignore
 *
 * Check out https://docs.tokaido.io for help managing your Tokaido environment
 *
 * Generated by Tokaido
 */

$databases['default']['default'] = [
  'host' => 'mysql',
  'database' => 'tokaido',
  'username' => 'tokaido',
  'password' => 'tokaido',
  'port' => 3306,
  'driver' => 'mysql',
  'namespace' => 'Drupal\\Core\\Database\\Driver\\mysql',
  'prefix' => '',
];

$settings['file_private_path'] = '/tokaido/site/private/default/';
$settings['file_temporary_path'] = '/tmp';

/*
* END Generated by Tokaido
*/

?>
`)
