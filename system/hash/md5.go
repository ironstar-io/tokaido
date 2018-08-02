package hash

import (
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"log"
)

// FileMD5 ...
func FileMD5(path string) (string, error) {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("There was an issue reading in your file: %v", err)
	}

	return BytesMD5(f)
}

// BytesMD5 ...
func BytesMD5(body []byte) (string, error) {
	// Open a new hash interface to write to
	h := md5.New()

	h.Write(body)

	// Get the 16 bytes hash
	hb := h.Sum(nil)[:16]

	return hex.EncodeToString(hb), nil
}
