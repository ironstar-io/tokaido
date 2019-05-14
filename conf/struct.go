package conf

// Project is a singular entry of a project name and path used in global config
type Project struct {
	Name string `yaml:"name,omitempty"`
	Path string `yaml:"path,omitempty"`
}

// Global contains all our global config settings that are saved in ~/.tok/global.yml
type Global struct {
	Syncservice string    `yaml:"syncservice,omitempty"`
	Projects    []Project `yaml:"projects,omitempty"`
}

// Config the application's configuration
// IMPORTANT!
// Casing of the `Config` struct properties is important to note
// All properties must be cased as capital letter first, followed by all lowercase
// eg. `Customcompose` (correct), `CustomCompose` (incorrect)
// This is to ensure they both conform to the golang convention
// and that they are able to be properly parsed by the `tok config-x` commands
type Config struct {
	Global  Global `yaml:"global,omitempty"`
	Tokaido struct {
		Config           string `yaml:"config,omitempty"`
		Customcompose    bool   `yaml:"customcompose"`
		Debug            bool   `yaml:"debug,omitempty"`
		Dependencychecks bool   `yaml:"dependencychecks"`
		Enableemoji      bool   `yaml:"enableemoji"`
		Force            bool   `yaml:"force,omitempty"`
		Yes              bool   `yaml:"yes,omitempty"`
		Phpversion       string `yaml:"phpversion"`
		Stability        string `yaml:"stability"`
		Xdebug           bool   `yaml:"xdebug"`
		Xdebugport       string `yaml:"xdebugport"`
		Project          struct {
			Name string `yaml:"name"`
		} `yaml:"project"`
	} `yaml:"tokaido"`
	Drupal struct {
		Path              string `yaml:"path,omitempty"`
		Majorversion      string `yaml:"majorversion,omitempty"`
		FilePublicPath    string `yaml:"filepublicpath,omitempty"`
		FilePrivatePath   string `yaml:"fileprivatepath,omitempty"`
		FileTemporaryPath string `yaml:"filetemporarypath,omitempty"`
	} `yaml:"drupal,omitempty"`
	Nginx struct {
		Workerconnections  string `yaml:"workerconnections,omitempty"`
		Clientmaxbodysize  string `yaml:"clientmaxbodysize,omitempty"`
		Keepalivetimeout   string `yaml:"keepalivetimeout,omitempty"`
		Fastcgireadtimeout string `yaml:"fastcgireadtimeout,omitempty"`
		Fastcgibuffers     string `yaml:"fastcgibuffers,omitempty"`
		Fastcgibuffersize  string `yaml:"fastcgibuffersize,omitempty"`
	} `yaml:"nginx,omitempty"`
	Fpm struct {
		Maxexecutiontime     string `yaml:"maxexecutiontime,omitempty"`
		Phpmemorylimit       string `yaml:"phpmemorylimit,omitempty"`
		Phpdisplayerrors     string `yaml:"phpdisplayerrors,omitempty"`
		Phplogerrors         string `yaml:"phplogerrors,omitempty"`
		Phpreportmemleaks    string `yaml:"phpreportmemleaks,omitempty"`
		Phppostmaxsize       string `yaml:"phppostmaxsize,omitempty"`
		Phpdefaultcharset    string `yaml:"phpdefaultcharset,omitempty"`
		Phpfileuploads       string `yaml:"phpfileuploads,omitempty"`
		Phpuploadmaxfilesize string `yaml:"phpuploadmaxfilesize,omitempty"`
		Phpmaxfileuploads    string `yaml:"phpmaxfileuploads,omitempty"`
		Phpallowurlfopen     string `yaml:"phpallowurlfopen,omitempty"`
	}
	System struct {
		Xdebug struct {
			Port      string `yaml:"port,omitempty"`
			Logpath   string `yaml:"logpath,omitempty"`
			Enabled   bool   `yaml:"enabled"`
			Autostart bool   `yaml:"autostart"`
		} `yaml:"xdebug,omitempty"`
		Syncsvc struct {
			Systemdpath string `yaml:"systemdpath,omitempty"`
			Launchdpath string `yaml:"launchdpath,omitempty"`
			Enabled     bool   `yaml:"enabled"`
		} `yaml:"syncsvc,omitempty"`
		Proxy struct {
			Enabled bool `yaml:"enabled"`
		} `yaml:"proxy,omitempty"`
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
		Command     string   `yaml:"command,omitempty"`
		Volumesfrom []string `yaml:"volumes_from,omitempty"`
		Dependson   []string `yaml:"depends_on,omitempty"`
		Environment []string `yaml:"environment,omitempty"`
		Volumes     []string `yaml:"volumes,omitempty"`
	} `yaml:"unison,omitempty"`
	Sync struct {
		Image       string            `yaml:"image,omitempty"`
		Hostname    string            `yaml:"hostname,omitempty"`
		Ports       []string          `yaml:"ports,omitempty"`
		Entrypoint  []string          `yaml:"entrypoint,omitempty"`
		User        string            `yaml:"user,omitempty"`
		Command     string            `yaml:"command,omitempty"`
		Volumesfrom []string          `yaml:"volumes_from,omitempty"`
		Dependson   []string          `yaml:"depends_on,omitempty"`
		Environment map[string]string `yaml:"environment,omitempty"`
		Volumes     []string          `yaml:"volumes,omitempty"`
		Restart     string            `yaml:"restart,omitempty"`
	} `yaml:"sync,omitempty"`
	Syslog struct {
		Image       string            `yaml:"image,omitempty"`
		Hostname    string            `yaml:"hostname,omitempty"`
		Ports       []string          `yaml:"ports,omitempty"`
		Entrypoint  []string          `yaml:"entrypoint,omitempty"`
		User        string            `yaml:"user,omitempty"`
		Command     string            `yaml:"command,omitempty"`
		Volumesfrom []string          `yaml:"volumes_from,omitempty"`
		Dependson   []string          `yaml:"depends_on,omitempty"`
		Environment map[string]string `yaml:"environment,omitempty"`
		Volumes     []string          `yaml:"volumes,omitempty"`
	} `yaml:"syslog,omitempty"`
	Haproxy struct {
		Image       string            `yaml:"image,omitempty"`
		Hostname    string            `yaml:"hostname,omitempty"`
		Ports       []string          `yaml:"ports,omitempty"`
		Entrypoint  []string          `yaml:"entrypoint,omitempty"`
		User        string            `yaml:"user,omitempty"`
		Command     string            `yaml:"command,omitempty"`
		Volumesfrom []string          `yaml:"volumes_from,omitempty"`
		Dependson   []string          `yaml:"depends_on,omitempty"`
		Environment map[string]string `yaml:"environment,omitempty"`
	} `yaml:"haproxy,omitempty"`
	Varnish struct {
		Image       string            `yaml:"image,omitempty"`
		Hostname    string            `yaml:"hostname,omitempty"`
		Ports       []string          `yaml:"ports,omitempty"`
		Entrypoint  []string          `yaml:"entrypoint,omitempty"`
		User        string            `yaml:"user,omitempty"`
		Command     string            `yaml:"command,omitempty"`
		Volumesfrom []string          `yaml:"volumes_from,omitempty"`
		Dependson   []string          `yaml:"depends_on,omitempty"`
		Environment map[string]string `yaml:"environment,omitempty"`
	} `yaml:"varnish,omitempty"`
	Nginx struct {
		Image       string            `yaml:"image,omitempty"`
		Hostname    string            `yaml:"hostname,omitempty"`
		Ports       []string          `yaml:"ports,omitempty"`
		Entrypoint  []string          `yaml:"entrypoint,omitempty"`
		User        string            `yaml:"user,omitempty"`
		Command     string            `yaml:"command,omitempty"`
		Volumesfrom []string          `yaml:"volumes_from,omitempty"`
		Volumes     []string          `yaml:"volumes,omitempty"`
		Dependson   []string          `yaml:"depends_on,omitempty"`
		Environment map[string]string `yaml:"environment,omitempty"`
	} `yaml:"nginx,omitempty"`
	Fpm struct {
		Image       string            `yaml:"image,omitempty"`
		Hostname    string            `yaml:"hostname,omitempty"`
		Ports       []string          `yaml:"ports,omitempty"`
		Entrypoint  []string          `yaml:"entrypoint,omitempty"`
		User        string            `yaml:"user,omitempty"`
		Command     string            `yaml:"command,omitempty"`
		Volumesfrom []string          `yaml:"volumes_from,omitempty"`
		Volumes     []string          `yaml:"volumes,omitempty"`
		Workingdir  string            `yaml:"working_dir,omitempty"`
		Dependson   []string          `yaml:"depends_on,omitempty"`
		Environment map[string]string `yaml:"environment,omitempty"`
	} `yaml:"fpm,omitempty"`
	Memcache struct {
		Enabled     bool              `yaml:"enabled,omitempty"`
		Image       string            `yaml:"image,omitempty"`
		Hostname    string            `yaml:"hostname,omitempty"`
		Ports       []string          `yaml:"ports,omitempty"`
		Entrypoint  []string          `yaml:"entrypoint,omitempty"`
		User        string            `yaml:"user,omitempty"`
		Command     string            `yaml:"command,omitempty"`
		Environment map[string]string `yaml:"environment,omitempty"`
	} `yaml:"memcache,omitempty"`
	Mysql struct {
		Image       string            `yaml:"image,omitempty"`
		Hostname    string            `yaml:"hostname,omitempty"`
		Ports       []string          `yaml:"ports,omitempty"`
		Entrypoint  []string          `yaml:"entrypoint,omitempty"`
		User        string            `yaml:"user,omitempty"`
		Command     string            `yaml:"command,omitempty"`
		Environment map[string]string `yaml:"environment,omitempty"`
		Volumesfrom []string          `yaml:"volumes_from,omitempty"`
		Volumes     []string          `yaml:"volumes,omitempty"`
	} `yaml:"mysql,omitempty"`
	Drush struct {
		Image       string            `yaml:"image,omitempty"`
		Hostname    string            `yaml:"hostname,omitempty"`
		Ports       []string          `yaml:"ports,omitempty"`
		Entrypoint  []string          `yaml:"entrypoint,omitempty"`
		User        string            `yaml:"user,omitempty"`
		Command     string            `yaml:"command,omitempty"`
		Workingdir  string            `yaml:"working_dir,omitempty"`
		Volumesfrom []string          `yaml:"volumes_from,omitempty"`
		Environment map[string]string `yaml:"environment,omitempty"`
		Volumes     []string          `yaml:"volumes,omitempty"`
	} `yaml:"drush,omitempty"`
	Solr struct {
		Enabled     bool              `yaml:"enabled,omitempty"`
		Image       string            `yaml:"image,omitempty"`
		Hostname    string            `yaml:"hostname,omitempty"`
		Ports       []string          `yaml:"ports,omitempty"`
		Entrypoint  []string          `yaml:"entrypoint,omitempty"`
		User        string            `yaml:"user,omitempty"`
		Command     string            `yaml:"command,omitempty"`
		Environment map[string]string `yaml:"environment,omitempty"`
	} `yaml:"solr,omitempty"`
	Redis struct {
		Enabled     bool              `yaml:"enabled,omitempty"`
		Image       string            `yaml:"image,omitempty"`
		Hostname    string            `yaml:"hostname,omitempty"`
		Ports       []string          `yaml:"ports,omitempty"`
		Entrypoint  []string          `yaml:"entrypoint,omitempty"`
		User        string            `yaml:"user,omitempty"`
		Command     string            `yaml:"command,omitempty"`
		Environment map[string]string `yaml:"environment,omitempty"`
	} `yaml:"redis,omitempty"`
	Mailhog struct {
		Enabled     bool              `yaml:"enabled,omitempty"`
		Image       string            `yaml:"image,omitempty"`
		Hostname    string            `yaml:"hostname,omitempty"`
		Ports       []string          `yaml:"ports,omitempty"`
		Entrypoint  []string          `yaml:"entrypoint,omitempty"`
		User        string            `yaml:"user,omitempty"`
		Command     string            `yaml:"command,omitempty"`
		Environment map[string]string `yaml:"environment,omitempty"`
	} `yaml:"mailhog,omitempty"`
	Adminer struct {
		Enabled     bool              `yaml:"enabled,omitempty"`
		Image       string            `yaml:"image,omitempty"`
		Hostname    string            `yaml:"hostname,omitempty"`
		Ports       []string          `yaml:"ports,omitempty"`
		Entrypoint  []string          `yaml:"entrypoint,omitempty"`
		User        string            `yaml:"user,omitempty"`
		Command     string            `yaml:"command,omitempty"`
		Environment map[string]string `yaml:"environment,omitempty"`
	} `yaml:"adminer,omitempty"`
	Kishu struct {
		Enabled     bool              `yaml:"enabled,omitempty"`
		Image       string            `yaml:"image,omitempty"`
		Environment map[string]string `yaml:"environment,omitempty"`
		Volumes     []string          `yaml:"volumes,omitempty"`
	} `yaml:"kishu,omitempty"`
}

// ComposeDotTok ...
type ComposeDotTok struct {
	Version  string   `yaml:"version,omitempty"`
	Services Services `yaml:"services,omitempty"`
}
