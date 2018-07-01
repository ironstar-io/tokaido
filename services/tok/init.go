package tok

import (
	"bitbucket.org/ironstar/tokaido-cli/conf"
	"bitbucket.org/ironstar/tokaido-cli/services/docker"
	"bitbucket.org/ironstar/tokaido-cli/services/drupal"
	"bitbucket.org/ironstar/tokaido-cli/services/git"
	"bitbucket.org/ironstar/tokaido-cli/services/unison"
	// "bitbucket.org/ironstar/tokaido-cli/services/xdebug"
	"bitbucket.org/ironstar/tokaido-cli/system"
	"bitbucket.org/ironstar/tokaido-cli/system/ssh"
)

// Init - The core run sheet of `tok init`
func Init() {
	c := conf.GetConfig()
	system.CheckDependencies()

	docker.FindOrCreateTokCompose()

	ssh.GenerateKeys()

	drupal.CheckSettings()

	git.IgnoreDefaults()

	unison.DockerUp()
	unison.CreateOrUpdatePrf()
	unison.Sync()
	if c.CreateSyncService {
		unison.CreateSyncService()
	}

	docker.Up()

	// xdebug.Configure()

	drupal.ConfigureSSH()
}
