package phpcs

import (
	"log"

	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/services/composer"
	"github.com/ironstar-io/tokaido/system/console"
	"github.com/ironstar-io/tokaido/utils"
)

var phpCS = "squizlabs/php_codesniffer=*"

func installPHPCS() {
	console.Println("🐘  Tokaido was unable to detect an installation of PHPCodeSniffer, which is required for the the `test:phpcs` command\n", "")
	confirmUpdate := utils.ConfirmationPrompt("Would you like us to attempt to install it for you automatically with composer?", "y")
	if confirmUpdate == false {
		log.Fatalf("We'll still attempt to run the tests, but they make not behave as expected!")
		return
	}

	composer.RequirePackage([]string{"--dev", phpCS})
}

func RunLinter() {
	v, err := composer.FindPackageVersion(phpCS)
	if err != nil {
		installPHPCS()
	}
	if v == "" {
		installPHPCS()
	}

	utils.StreamOSCmdContext("/tokaido/site", "./vendor/bin/phpcs", "/tokaido/site/"+conf.GetConfig().Drupal.Path)
}
