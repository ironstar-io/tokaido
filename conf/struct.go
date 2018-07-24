package conf

// Config the application's configuration
type Config struct {
	Tokaido struct {
		Force            bool   `yaml:"force,omitempty"`
		CustomCompose    bool   `yaml:"customcompose,omitempty"`
		Debug            bool   `yaml:"debug,omitempty"`
		Config           string `yaml:"config,omitempty"`
		EnableEmoji      bool   `yaml:"enableemoji,omitempty"`
		BetaContainers   bool   `yaml:"betacontainers,omitempty"`
		DependencyChecks bool   `yaml:"dependencychecks"`
		Project          struct {
			Name string `yaml:"name,omitempty"`
			Path string `yaml:"path,omitempty"`
		} `yaml:"project,omitempty"`
	} `yaml:"tokaido,omitempty"`
	Drupal struct {
		Path         string `yaml:"path,omitempty"`
		MajorVersion string `yaml:"majorversion,omitempty"`
	} `yaml:"drupal,omitempty"`
	System struct {
		Xdebug struct {
			Port      string `yaml:"port,omitempty"`
			LogPath   string `yaml:"logpath,omitempty"`
			Enabled   bool   `yaml:"enabled,omitempty"`
			Autostart bool   `yaml:"autostart,omitempty"`
		} `yaml:"xdebug,omitempty"`
		SyncSvc struct {
			SystemdPath string `yaml:"systemdpath,omitempty"`
			LaunchdPath string `yaml:"launchdpath,omitempty"`
			Enabled     bool   `yaml:"enabled"`
		} `yaml:"syncsvc,omitempty"`
	} `yaml:"system,omitempty"`
	Services Services `yaml:"services,omitempty"`
}

// Services ...
type Services struct {
	Unison struct {
		Image       string   `yaml:"image,omitempty"`
		Hostname    string   `yaml:"hostname,omitempty"`
		Ports       []string `yaml:"ports,omitempty"`
		Entrypoint  []string `yaml:"entrypoint,omitempty"`
		User        string   `yaml:"user,omitempty"`
		Cmd         string   `yaml:"cmd,omitempty"`
		VolumesFrom []string `yaml:"volumes_from,omitempty"`
		DependsOn   []string `yaml:"depends_on,omitempty"`
		Environment []string `yaml:"environment,omitempty"`
		Volumes     []string `yaml:"volumes,omitempty"`
	} `yaml:"unison,omitempty"`
	Syslog struct {
		Image       string            `yaml:"image,omitempty"`
		Hostname    string            `yaml:"hostname,omitempty"`
		Ports       []string          `yaml:"ports,omitempty"`
		Entrypoint  []string          `yaml:"entrypoint,omitempty"`
		User        string            `yaml:"user,omitempty"`
		Cmd         string            `yaml:"cmd,omitempty"`
		VolumesFrom []string          `yaml:"volumes_from,omitempty"`
		DependsOn   []string          `yaml:"depends_on,omitempty"`
		Environment map[string]string `yaml:"environment,omitempty"`
		Volumes     []string          `yaml:"volumes,omitempty"`
	} `yaml:"syslog,omitempty"`
	Haproxy struct {
		Image       string            `yaml:"image,omitempty"`
		Hostname    string            `yaml:"hostname,omitempty"`
		Ports       []string          `yaml:"ports,omitempty"`
		Entrypoint  []string          `yaml:"entrypoint,omitempty"`
		User        string            `yaml:"user,omitempty"`
		Cmd         string            `yaml:"cmd,omitempty"`
		VolumesFrom []string          `yaml:"volumes_from,omitempty"`
		DependsOn   []string          `yaml:"depends_on,omitempty"`
		Environment map[string]string `yaml:"environment,omitempty"`
	} `yaml:"haproxy,omitempty"`
	Varnish struct {
		Image       string            `yaml:"image,omitempty"`
		Hostname    string            `yaml:"hostname,omitempty"`
		Ports       []string          `yaml:"ports,omitempty"`
		Entrypoint  []string          `yaml:"entrypoint,omitempty"`
		User        string            `yaml:"user,omitempty"`
		Cmd         string            `yaml:"cmd,omitempty"`
		VolumesFrom []string          `yaml:"volumes_from,omitempty"`
		DependsOn   []string          `yaml:"depends_on,omitempty"`
		Environment map[string]string `yaml:"environment,omitempty"`
	} `yaml:"varnish,omitempty"`
	Nginx struct {
		Image       string            `yaml:"image,omitempty"`
		Hostname    string            `yaml:"hostname,omitempty"`
		Ports       []string          `yaml:"ports,omitempty"`
		Entrypoint  []string          `yaml:"entrypoint,omitempty"`
		User        string            `yaml:"user,omitempty"`
		Cmd         string            `yaml:"cmd,omitempty"`
		VolumesFrom []string          `yaml:"volumes_from,omitempty"`
		DependsOn   []string          `yaml:"depends_on,omitempty"`
		Environment map[string]string `yaml:"environment,omitempty"`
	} `yaml:"nginx,omitempty"`
	Fpm struct {
		Image       string            `yaml:"image,omitempty"`
		Hostname    string            `yaml:"hostname,omitempty"`
		Ports       []string          `yaml:"ports,omitempty"`
		Entrypoint  []string          `yaml:"entrypoint,omitempty"`
		User        string            `yaml:"user,omitempty"`
		Cmd         string            `yaml:"cmd,omitempty"`
		VolumesFrom []string          `yaml:"volumes_from,omitempty"`
		WorkingDir  string            `yaml:"working_dir,omitempty"`
		DependsOn   []string          `yaml:"depends_on,omitempty"`
		Environment map[string]string `yaml:"environment,omitempty"`
	} `yaml:"fpm,omitempty"`
	Memcache struct {
		Enabled     bool              `yaml:"enabled,omitempty"`
		Image       string            `yaml:"image,omitempty"`
		Hostname    string            `yaml:"hostname,omitempty"`
		Ports       []string          `yaml:"ports,omitempty"`
		Entrypoint  []string          `yaml:"entrypoint,omitempty"`
		User        string            `yaml:"user,omitempty"`
		Cmd         string            `yaml:"cmd,omitempty"`
		Environment map[string]string `yaml:"environment,omitempty"`
	} `yaml:"memcache,omitempty"`
	Mysql struct {
		Image       string            `yaml:"image,omitempty"`
		Hostname    string            `yaml:"hostname,omitempty"`
		Ports       []string          `yaml:"ports,omitempty"`
		Entrypoint  []string          `yaml:"entrypoint,omitempty"`
		User        string            `yaml:"user,omitempty"`
		Cmd         string            `yaml:"cmd,omitempty"`
		Environment map[string]string `yaml:"environment,omitempty"`
	} `yaml:"mysql,omitempty"`
	Drush struct {
		Image       string            `yaml:"image,omitempty"`
		Hostname    string            `yaml:"hostname,omitempty"`
		Ports       []string          `yaml:"ports,omitempty"`
		Entrypoint  []string          `yaml:"entrypoint,omitempty"`
		User        string            `yaml:"user,omitempty"`
		Cmd         string            `yaml:"cmd,omitempty"`
		WorkingDir  string            `yaml:"working_dir,omitempty"`
		VolumesFrom []string          `yaml:"volumes_from,omitempty"`
		Environment map[string]string `yaml:"environment,omitempty"`
	} `yaml:"drush,omitempty"`
	Solr struct {
		Enabled     bool              `yaml:"enabled,omitempty"`
		Image       string            `yaml:"image,omitempty"`
		Hostname    string            `yaml:"hostname,omitempty"`
		Ports       []string          `yaml:"ports,omitempty"`
		Entrypoint  []string          `yaml:"entrypoint,omitempty"`
		User        string            `yaml:"user,omitempty"`
		Cmd         string            `yaml:"cmd,omitempty"`
		Environment map[string]string `yaml:"environment,omitempty"`
	} `yaml:"solr,omitempty"`
}

// ComposeDotTok ...
type ComposeDotTok struct {
	Version  string
	Services Services `yaml:"services,omitempty"`
}
