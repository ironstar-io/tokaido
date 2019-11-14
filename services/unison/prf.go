package unison

import (
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
	UnisonPort string
	SyncPath   string
}

// CreateOrUpdatePrf - Create or Update a `.prf` file in `~/.unison/`
func CreateOrUpdatePrf(unisonPort, prfName, syncPath string) {
	f := getPrfPath(prfName)

	var _, err = os.Stat(f)

	if os.IsNotExist(err) {
		generatePrf(unisonPort, syncPath, f)
	} else {
		updatePrf(unisonPort, f)
	}
}

func getPrfPath(filename string) string {
	return filepath.Join(fs.HomeDir(), ".unison", filename+".prf")
}

// generatePrf - Generate a `.prf` file for unison
func generatePrf(unisonPort, syncPath, prfPath string) {
	s := prf{UnisonPort: unisonPort, SyncPath: syncPath}

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

	createPrf(tpl.Bytes(), prfPath)
}

// createPrf - Write generated `.prf` file to `~/.unison/`
func createPrf(body []byte, path string) {
	fs.Mkdir(filepath.Join(fs.HomeDir(), ".unison"))

	fs.TouchByteArray(path, body)
}

// updatePrf - Update a `.prf` file in `~/.unison/`
func updatePrf(unisonPort, prfPath string) {
	f, err := os.Open(prfPath)
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
			buffer.Write([]byte(portLine + unisonPort + "\n"))
		} else {
			buffer.Write([]byte(scanner.Text() + "\n"))
		}
	}

	fs.Replace(prfPath, buffer.Bytes())
}
