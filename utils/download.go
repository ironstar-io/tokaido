package utils

import (
	"io"
	"log"
	"net/http"
	"runtime"
	"errors"
	"os"
	"strconv"
)

// DownloadFile - download a url to a local file.
func DownloadFile(filepath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	if resp.StatusCode == 404 {
		return errors.New("The requested release was not found (404)")
	}
	if resp.StatusCode != 200 {
		return errors.New("GitHub responded with error code " + strconv.Itoa(resp.StatusCode))
	}

	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Change file permission bit
	err = os.Chmod(filepath, 0755)
	if err != nil {
		log.Println(err)
	}

	if runtime.GOOS != "windows" {
		// Change file ownership.
		err = os.Chown(filepath, os.Getuid(), os.Getgid())
		if err != nil {
			log.Println(err)
		}
	}

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
