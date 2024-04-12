#!/bin/sh

set -e

echo "clean up db"
/app/migrate -path /app/migrations -database "$DB_SOURCE" down -all

echo "run the database migrations"
/app/migrate -path /app/migrations -database "$DB_SOURCE" -verbose up

echo "start the app"
exec "$@"