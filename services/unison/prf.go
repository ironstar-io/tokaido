package unison

import (
	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/services/unison/templates"
	"github.com/ironstar-io/tokaido/system/fs"

	"bufio"
	"bytes"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type prf struct {
	UnisonPort  string
	ProjectPath string
}

// CreateOrUpdatePrf - Create or Update a `.prf` file in `~/.unison/`
func CreateOrUpdatePrf() {
	// detect if file exists
	var _, err = os.Stat(getPrfPath())

	// create file if not exists
	if os.IsNotExist(err) {
		generatePrf()
	} else {
		updatePrf()
	}
}

func getPrfPath() string {
	return filepath.Join(fs.HomeDir(), "/.unison/", conf.GetConfig().Tokaido.Project.Name+".prf")
}

// generatePrf - Generate a `.prf` file for unison
func generatePrf() {
	config := conf.GetConfig()
	s := prf{UnisonPort: LocalPort(), ProjectPath: config.Tokaido.Project.Path}

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

	createPrf(tpl.Bytes(), getPrfPath())
}

// createPrf - Write generated `.prf` file to `~/.unison/`
func createPrf(body []byte, path string) {
	fs.Mkdir(filepath.Join(fs.HomeDir(), "/.unison"))

	fs.TouchByteArray(path, body)
}

// updatePrf - Update a `.prf` file in `~/.unison/`
func updatePrf() {
	f, err := os.Open(getPrfPath())
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

	fs.Replace(getPrfPath(), buffer.Bytes())
}
