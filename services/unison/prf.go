package unison

import (
	"bitbucket.org/ironstar/tokaido-cli/conf"
	"bitbucket.org/ironstar/tokaido-cli/services/unison/templates"
	"bitbucket.org/ironstar/tokaido-cli/system/fs"

	"bufio"
	"bytes"
	"log"
	"os"
	"strings"
	"text/template"
)

type prf struct {
	UnisonPort  string
	ProjectPath string
}

// GetPrfPath ...
func GetPrfPath() string {
	config := conf.GetConfig()
	return fs.HomeDir() + "/.unison/" + config.Project
}

func checkDotUnison() {
	dotUnison := fs.HomeDir() + "/.unison"
	var _, err = os.Stat(dotUnison)

	if os.IsNotExist(err) {
		os.Mkdir(dotUnison, os.ModePerm)
	}
}

// CreateOrUpdatePrf - Create or Update a `.prf` file in `~/.unison/`
func CreateOrUpdatePrf() {
	// detect if file exists
	var _, err = os.Stat(GetPrfPath() + ".prf")

	// create file if not exists
	if os.IsNotExist(err) {
		generatePrf()
	} else {
		updatePrf()
	}
}

// generatePrf - Generate a `.prf` file for unison
func generatePrf() {
	config := conf.GetConfig()
	s := prf{UnisonPort: LocalPort(), ProjectPath: config.Path}

	tmpl := template.New("unison.prf")
	tmpl, err := tmpl.Parse(unisontmpl.PRFTemplateStr)

	if err != nil {
		log.Fatal("Parse: ", err)
		return
	}

	var tpl bytes.Buffer
	if err := tmpl.Execute(&tpl, s); err != nil {
		log.Fatal("Parse: ", err)
		return
	}

	createPrf(tpl.String(), GetPrfPath()+".prf")
}

// createPrf - Write generated `.prf` file to `~/.unison/`
func createPrf(body string, path string) {
	checkDotUnison()

	var file, err = os.Create(path)
	if err != nil {
		log.Fatal("Create: ", err)
	}

	_, _ = file.WriteString(body)

	defer file.Close()
}

// updatePrf - Update a `.prf` file in `~/.unison/`
func updatePrf() {
	f, err := os.Open(GetPrfPath() + ".prf")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	portLine := "root = socket://localhost:"
	// Splits on newlines by default.
	scanner := bufio.NewScanner(f)

	var buffer bytes.Buffer
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), portLine) {
			buffer.Write([]byte(portLine + LocalPort() + "\n"))
		} else {
			buffer.Write([]byte(scanner.Text() + "\n"))
		}
	}

	createPrf(buffer.String(), GetPrfPath()+".tmp.prf")
	replacePrf()
}

// replacePrf - Replace `.tmp.prf` with `.prf` file in `~/.unison/`
func replacePrf() {
	rootPrf := GetPrfPath()
	mainPrf := rootPrf + ".prf"
	tmpPrf := rootPrf + ".tmp.prf"

	// Remove the original .prf file
	os.Remove(mainPrf)

	// Rename `.tmp.prf` to be the new `.prf` file
	os.Rename(tmpPrf, mainPrf)
}
