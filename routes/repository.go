package routes

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/godror/godror"
)

// Repository DBを開
func Repository() *sql.DB {

	// this is oracle way
	var ocistring string
	if ocistring = os.Getenv("OCISTRING"); ocistring == "" {
		fmt.Printf("OCI adapter string not specified in 'OCISTRING' environment variable. Program will quit")
		os.Exit(1)
	}

	db, err := sql.Open("godror", ocistring)

	if err != nil {
		panic(err)
	}

	return db
}
