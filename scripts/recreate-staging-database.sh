cd `git rev-parse --show-toplevel`

DB_DRIVER=postgres
export DB_DRIVER
DB_DBNAME=mlmodelscope
export DB_DBNAME
DB_PORT=5432
export DB_PORT

if [ -z "$1" ]; then
  echo "You must pass an import file name as the first argument to this script!"
  exit 1
fi

echo "Staging database host: "
read DB_HOST

echo "Staging database user: "
read DB_USER

echo "Staging database password: "
read PGPASSWORD

DB_PASSWORD="$PGPASSWORD"

if [ -z "$DB_HOST" ]; then
  echo "You must provide a database host!"
  exit 1
fi
export DB_HOST

if [ -z "$DB_USER" ]; then
  echo "You must provide a database user!"
  exit 1
fi
export DB_USER

if [ -z "$PGPASSWORD" ]; then
  echo "You must provide a database password!"
  exit 1
fi
export DB_PASSWORD
export PGPASSWORD

cd import
go build .

psql -h $DB_HOST \
     -p $DB_PORT \
     -U $DB_USER \
     -d postgres \
     -c "DROP DATABASE mlmodelscope;"

psql -h $DB_HOST \
     -p $DB_PORT \
     -U $DB_USER \
     -d postgres \
     -c "CREATE DATABASE mlmodelscope WITH OWNER mlmodelscope;"

./import "$1"
