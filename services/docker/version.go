package docker

import (
	"strconv"
	"fmt"
	"log"

	"github.com/docker/docker/client"
)

// CheckClientVersion ensures that the client running on this system is supported
func CheckClientVersion() {
	cli, err := client.NewEnvClient()
	if err != nil {
		fmt.Println("Error while checking Docker client version. You must have Docker version 18.02 or higher")
		log.Fatal(err)
	}

	v, _ := strconv.ParseFloat(cli.ClientVersion(), 32)

	if v < 1.36 {
		panic("Error while checking Docker client version. You must have Docker version 18.02 or higher")
	}
}
