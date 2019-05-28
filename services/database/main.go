package database

import (
	"database/sql"
	"fmt"

	"github.com/ironstar-io/tokaido/services/docker"

	_ "github.com/go-sql-driver/mysql" // Official docs recommend blank import
)

// Connect opens a connection to the named database
func Connect(dbname string) (*sql.DB, error) {

	port := docker.LocalPort("mysql", "3306")

	db, err := sql.Open("mysql", "tokaido:tokaido@tcp(127.0.0.1:"+port+")/"+dbname)
	if err != nil {
		fmt.Println("err not nil")
		return &sql.DB{}, err
	}
	// defer db.Close()

	return db, nil
}

// ConnectRoot opens a connection to the named database as root
func ConnectRoot(dbname string) (*sql.DB, error) {

	port := docker.LocalPort("mysql", "3306")

	db, err := sql.Open("mysql", "root:tokaido@tcp(127.0.0.1:"+port+")/"+dbname)
	if err != nil {
		fmt.Println("err not nil")
		return &sql.DB{}, err
	}
	// defer db.Close()

	return db, nil
}
