# Database Migrator

This tool will run database migrations against the provided database.

## Building

`go build .`

## Running

Database details are provided by a set of environment variables:

* DB_DRIVER (i.e.: postgres)
* DB_HOST (i.e.: localhost)
* DB_PORT
* DB_USER - a user with permission to insert rows into the `models` table
* DB_PASSWORD
* DB_DBNAME

`./migrate`
