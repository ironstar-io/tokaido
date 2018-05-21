package main

import (
	"log"

	"bitbucket.org/ironstar/tokaido-cli/cmd"
)

func main() {
	if err := cmd.RootCmd().Execute(); err != nil {
		log.Fatal(err)
	}
}
