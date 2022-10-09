#!/bin/bash
set -e
echo "Printing environment variables:"
echo "DB USER: ${DBUSER}"
echo "DB NAME: ${DB_NAME}"
echo "DB HOST: ${DB_HOST}"
echo "POSTGRES USER: ${DB_USER}"
echo "Done printing environment variables:"
echo "Wait for Postgres to start"
until pg_isready -h ${DB_HOST} -p 5432 -U ${DB_USER}
  do
  echo "Waiting for postgres to start at ${DB_HOST}..."
  sleep 2;
done
echo "Postgres has started"
echo "Setting up the Postgres user"
POSTGRES="psql --username ${DB_USER}"
echo "Creating database: ${DB_NAME}"
$POSTGRES <<EOSQL
CREATE DATABASE "${DB_NAME}" OWNER ${DB_USER};
EOSQL
echo "Initializing database..."
psql -d ${DB_NAME} -a -U${DB_USER} -f ./db/init/0001_CREATE_TABLES.sql -f ./db/init/0002_FILL_DATA.sql &
SQL_PID=$!
wait $SQL_PID
echo "Finished initializing database"
echo "Running Go Tests"
go test -v ./... -coverprofile coverage/cover.out &
GO_PID=$!
wait $GO_PID
echo "Finished with go tests"