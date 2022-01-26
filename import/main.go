package main

import (
	"api/db"
	"api/db/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type ImportJson struct {
	Manifests []models.Model `json:"manifests"`
}

func main() {
	if len(os.Args) < 2 {
		println("Usage: import /path/to/file-to-import.json")
		os.Exit(1)
	}

	fileToImport := os.Args[1]
	bytes, err := ioutil.ReadFile(fileToImport)

	if err != nil {
		fmt.Printf("failed to read import file: %s", err.Error())
		os.Exit(2)
	}

	m := ImportJson{}
	if err = json.Unmarshal(bytes, &m); err != nil {
		fmt.Printf("failed to parse import JSON: %s", err.Error())
		os.Exit(3)
	}

	database, err := db.OpenDb()

	if err != nil {
		fmt.Printf("failed to open database: %s", err.Error())
		os.Exit(4)
	}

	if err = database.Migrate(); err != nil {
		fmt.Printf("failed to migrate database: %s", err.Error())
		os.Exit(5)
	}

	for _, model := range m.Manifests {
		if f, err := database.QueryFrameworks(&models.Framework{Name: model.Framework.Name, Version: model.Framework.Version}); err == nil {
			model.Framework = nil
			model.FrameworkID = int(f.ID)
		} else {
			model.Framework.Architectures = []models.Architecture{
				{
					Name:        "amd64",
					FrameworkID: model.Framework.ID,
				},
			}
		}

		err = database.CreateModel(&model)

		if err != nil {
			fmt.Printf("failed to create model: %s", err.Error())
			os.Exit(6)
		}
	}

	println("Done!")
}
