package types

// Templates is a list of drupal templates available for download
type Templates struct {
	Template []Template `yaml:"templates"`
}

// Template ...
type Template struct {
	Description     string   `yaml:"description"`
	DrupalVersion   int      `yaml:"drupal_version"`
	Maintainer      string   `yaml:"maintainer"`
	Name            string   `yaml:"name"`
	PackageFilename string   `yaml:"package_filename"`
	PostUpCommands  []string `yaml:"post_up_commands,omitempty"`
}
