# Model Importer

This tool will import ML model details into a database from a JSON file containing
an array of model manifests.

## Building

`go build .`

## Running

You must define a set of environment variables that the tool will use to connect
to your database:

* DB_DRIVER (i.e.: postgres)
* DB_HOST (i.e.: localhost)
* DB_PORT
* DB_USER - a user with permission to insert rows into the `models` table
* DB_PASSWORD
* DB_DBNAME

`./import /path/to/models-to-import.json`