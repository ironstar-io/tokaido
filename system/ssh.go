package system

import (
	"bitbucket.org/ironstar/tokaido-cli/system/fs"

	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/crypto/ssh"
)

var sshPriv = fs.HomeDir() + "/.ssh/tok_ssh.key"
var sshPub = fs.HomeDir() + "/.ssh/tok_ssh.pub"
var tokDir = fs.WorkDir() + "/.tok"

// GenerateSSHKeys ...
func GenerateSSHKeys() {
	var _, err = os.Stat(sshPub)

	// create file if not exists
	if os.IsNotExist(err) {
		fmt.Println("Generating a new set of SSH keys")
		generateAndCopyPub()
	} else {
		copyPub()
	}
}

func copyPub() {
	CheckAndCreateFolder(tokDir)
	CheckAndCreateFolder(tokDir + "/local")

	fs.Copy(sshPub, tokDir+"/local/tok_ssh.pub-copy")
	replacePub()
}

// replacePub - Replace `.pub-copy` with `.pub` file in `./.tok/local`
func replacePub() {
	mainPub := tokDir + "/local/tok_ssh.pub"
	copyPub := mainPub + "-copy"

	// Remove the original .pub key
	os.Remove(mainPub)

	// Rename `.pub-copy` to be the new `.pub` key
	os.Rename(copyPub, mainPub)
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
