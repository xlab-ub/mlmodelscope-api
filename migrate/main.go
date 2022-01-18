package main

import (
	"api/db"
	"fmt"
	"os"
)

func main() {
	database, err := db.OpenDb()
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
		os.Exit(1)
	}

	err = database.Migrate()

	if err != nil {
		fmt.Printf("Error: %s", err.Error())
		os.Exit(1)
	}

	println("Done!")
}
