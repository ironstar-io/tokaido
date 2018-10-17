package drupaltmpl

// Settings - Collection of variables required for building Drupal settings files
type Settings struct {
	HashSalt          string
	ProjectName       string
	FilePublicPath    string
	FilePrivatePath   string
	FileTemporaryPath string
}

// SettingsD7Tok - docroot/sites/default/settings.tok.php for Drupal 7
func (s *Settings) SettingsD7Tok() []byte {
	return []byte(`<?php

    /**
     * @file
     * Configuration file for Tokaido local dev environments. Add this to .gitignore
     *
     * Check out https://tokaido.io/docs for help managing your Tokaido environment
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

    ` + s.generateD7Paths() + `

    $drupal_hash_salt = '` + s.HashSalt + `';

    $base_url = 'https://` + s.ProjectName + `.local.tokaido.io:5154';

    /*
    * END Generated by Tokaido
    */

  `)
}

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

// SettingsD8Tok - docroot/sites/default/settings.tok.php for Drupal 8
func (s *Settings) SettingsD8Tok() []byte {
	return []byte(`<?php

/**
 * @file
 * Configuration file for Tokaido local dev environments. Add this to .gitignore
 *
 * Check out https://tokaido.io/docs for help managing your Tokaido environment
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

` + s.generateD8Paths() + `

$settings['hash_salt'] = '` + s.HashSalt + `';

if (\Drupal::hasService('cache.backend.memcache')) {
  $settings['cache']['default'] = 'cache.backend.memcache';
  $settings['memcache_storage']['key_prefix'] = '';
  $settings['memcache_storage']['memcached_servers'] = ['memcache:11211' => 'default'];
}

/*
* END Generated by Tokaido
*/
`)
}

func (s *Settings) generateD7Paths() string {
	var ps string

	ps = ps + "$conf['file_private_path'] = '/tokaido/" + s.FilePrivatePath + "';"
	if s.FilePublicPath != "" {
		ps = ps + "\n$conf['file_public_path'] = '/tokaido/" + s.FilePublicPath + "';"
	}
	ps = ps + "\n$conf['file_temporary_path'] = '" + s.FileTemporaryPath + "';"

	return ps
}

func (s *Settings) generateD8Paths() string {
	var ps string

	ps = ps + "$settings['file_private_path'] = '/tokaido/" + s.FilePrivatePath + "';"
	if s.FilePublicPath != "" {
		ps = ps + "\n$settings['file_public_path'] = '/tokaido/" + s.FilePublicPath + "';"
	}
	ps = ps + "\n$settings['file_temporary_path'] = '" + s.FileTemporaryPath + "';"

	return ps
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
