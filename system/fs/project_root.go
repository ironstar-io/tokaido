package fs

import (
	"log"
	"os"
)

// ProjectRoot - Find the root directory path for the project
func ProjectRoot() string {

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	return dir
}

func CheckDir(path string) string {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		fmt.Println(file.Name())
	}

	files, err := ioutil.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		fmt.Println(file.Name())
	}

	dirname, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Current directory: %v\n", dirname)
	dir, err := os.Open(filepath.Join(dirname, "../../../../../../"))
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", dir.Name())
}
