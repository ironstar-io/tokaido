package system

import (
	"bitbucket.org/ironstar/tokaido-cli/system/linux"
	"bitbucket.org/ironstar/tokaido-cli/system/osx"
	"bitbucket.org/ironstar/tokaido-cli/utils"
	"fmt"
	"log"

	"os"
	"strconv"
	"strings"
)

// GetPermissionsMask - Return the permissions mask for a file/directory
func GetPermissionsMask(path string) os.FileMode {
	var permissionString string
	if utils.CheckOS() == "osx" {
		permissionString = osx.GetPermissionsMask(path)
	} else {
		permissionString = linux.GetPermissionsMask(path)
	}

	maskString := strings.Split(permissionString, " ")[0]
	fmt.Println(permissionString)
	fmt.Println(maskString)
	if len(maskString) == 3 {
		maskString = "0" + maskString
	}
	maskInt, err := strconv.ParseUint(maskString, 10, 32)
	if err != nil {
		log.Fatal("A fatal error has occurred when processing file permissions", err)
	}

	return os.FileMode(maskInt)
}
