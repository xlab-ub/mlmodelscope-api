package api_db

import "api/db"

var database db.Db

func GetDatabase() (db.Db, error) {
	var err error
	if database == nil {
		database, err = db.OpenDb()
	}

	return database, err
}

func CloseDatabase() {
	database = nil
}
