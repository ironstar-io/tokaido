package ssh

import (
	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/services/docker"
	"github.com/ironstar-io/tokaido/system/console"
	"github.com/ironstar-io/tokaido/system/fs"
	"github.com/ironstar-io/tokaido/utils"

	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/crypto/ssh"
)

var sshPriv = filepath.Join(fs.HomeDir(), "/.ssh/tok_ssh.key")
var sshPub = filepath.Join(fs.HomeDir(), "/.ssh/tok_ssh.pub")
var tokDir = filepath.Join(fs.WorkDir(), "/.tok")

// CheckKey ...
func CheckKey() {
	localPort := docker.LocalPort("drush", "22")
	cmdStr := `ssh ` + conf.GetConfig().Tokaido.Project.Name + `.tok -q -p ` + localPort + ` -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no -C "echo 1" | echo $?`

	keyResult := utils.BashStringCmd(cmdStr)
	if keyResult == "0" {
		console.Println("✅  SSH access is configured", "√")
		return
	}

	console.Println(`😓  SSH access not configured

Tokaido is running but your SSH access to the Drush container looks broken.
Make sure you have an SSH public key uploaded in './.tok/local/ssh_key.pub'.

You should be able to run 'tok repair' to attempt to fix this automatically
	`, "×")
	os.Exit(1)
}

// GenerateKeys ...
func GenerateKeys() {
	var _, err = os.Stat(sshPub)

	// create file if not exists
	if os.IsNotExist(err) {
		fmt.Println("Generating a new set of SSH keys")
		fs.Mkdir(filepath.Join(fs.HomeDir(), "/.ssh"))
		generateAndCopyPub()
	} else {
		copyPub()
	}
}

func copyPub() {
	fs.Mkdir(tokDir)
	fs.Mkdir(filepath.Join(tokDir, "/local"))

	fs.Copy(sshPub, filepath.Join(tokDir, "/local/ssh_key.pub"))
}

func generateAndCopyPub() {
	bitSize := 4096

	privateKey, err := generatePrivateKey(bitSize)
	if err != nil {
		log.Fatal(err.Error())
	}

	publicKeyBytes, err := generatePublicKey(&privateKey.PublicKey)
	if err != nil {
		log.Fatal(err.Error())
	}

	privateKeyBytes := encodePrivateKeyToPEM(privateKey)

	err = writeKeyToFile(privateKeyBytes, sshPriv)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = writeKeyToFile([]byte(publicKeyBytes), sshPub)
	if err != nil {
		log.Fatal(err.Error())
	}

	copyPub()
}

// generatePrivateKey creates a RSA Private Key of specified byte size
func generatePrivateKey(bitSize int) (*rsa.PrivateKey, error) {
	// Private Key generation
	privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return nil, err
	}

	// Validate Private Key
	err = privateKey.Validate()
	if err != nil {
		return nil, err
	}

	log.Println("Private Key generated")
	return privateKey, nil
}

// encodePrivateKeyToPEM encodes Private Key from RSA to PEM format
func encodePrivateKeyToPEM(privateKey *rsa.PrivateKey) []byte {
	// Get ASN.1 DER format
	privDER := x509.MarshalPKCS1PrivateKey(privateKey)

	// pem.Block
	privBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privDER,
	}

	// Private key in PEM format
	privatePEM := pem.EncodeToMemory(&privBlock)

	return privatePEM
}

// generatePublicKey take a rsa.PublicKey and return bytes suitable for writing to .pub file
// returns in the format "ssh-rsa ..."
func generatePublicKey(privatekey *rsa.PublicKey) ([]byte, error) {
	publicRsaKey, err := ssh.NewPublicKey(privatekey)
	if err != nil {
		return nil, err
	}

	pubKeyBytes := ssh.MarshalAuthorizedKey(publicRsaKey)

	log.Println("Public key generated")
	return pubKeyBytes, nil
}

// writePemToFile writes keys to a file
func writeKeyToFile(keyBytes []byte, saveFileTo string) error {
	err := ioutil.WriteFile(saveFileTo, keyBytes, 0600)
	if err != nil {
		return err
	}

	log.Printf("Key saved to: %s", saveFileTo)
	return nil
}
