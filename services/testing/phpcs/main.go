package phpcs

import (
	"log"
	"path/filepath"

	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/services/composer"
	"github.com/ironstar-io/tokaido/system/console"
	"github.com/ironstar-io/tokaido/system/fs"
	"github.com/ironstar-io/tokaido/system/ssh"
	"github.com/ironstar-io/tokaido/utils"
)

var phpCSCInstaller = "dealerdirect/phpcodesniffer-composer-installer"
var drupalCoder = "drupal/coder"
var phpCSXML = "phpcs.xml.dist"

func installReq(packageName string) {
	console.Println("🐘  Tokaido was unable to detect an installation of '"+packageName+"', which is required for this command\n", "")
	confirmUpdate := utils.ConfirmationPrompt("Would you like us to attempt to install it for you automatically with composer?", "y")
	if confirmUpdate == false {
		log.Fatalf("Unable to run this command without '" + packageName + "' being installed")
		return
	}

	composer.RequirePackage([]string{"--dev", packageName})
}

// CheckReq ...
func CheckReq(packageName string) {
	v, err := composer.FindPackageVersion(packageName)
	if err != nil {
		installReq(packageName)
		return
	}
	if v == "" {
		installReq(packageName)
		return
	}

	return
}

func checkPHPCSXML() {
	f := filepath.Join(conf.GetProjectPath(), phpCSXML)
	if fs.CheckExists(f) == false {
		fs.TouchByteArray(f, phpCSXMLTemplate(conf.GetConfig().Drupal.Path))
	}
}

// CheckReqs - Ensure installation has the required dependencies to run phpcs
func CheckReqs() {
	CheckReq(phpCSCInstaller)
	CheckReq(drupalCoder)

	checkPHPCSXML()
}

// RunLinter - Run the PHPCodeSniffer linter (phpcs)
func RunLinter() {
	console.Println("⛑  Starting phpcs linting. This can take some time depending on your phpcs.xml.dist configuration.", "")

	ssh.StreamConnectCommand([]string{"/tokaido/site/vendor/bin/phpcs", "/tokaido/site/" + conf.GetConfig().Drupal.Path})
}

// RunLinterFix - Run the PHPCodeSniffer linter and attempt to fix automatically (phpcbf)
func RunLinterFix() {
	console.Println("⛑  Starting phpcbf linting and automated fixes. This can take some time depending on your phpcs.xml.dist configuration.", "")

	ssh.StreamConnectCommand([]string{"/tokaido/site/vendor/bin/phpcbf", "/tokaido/site/" + conf.GetConfig().Drupal.Path})
}
