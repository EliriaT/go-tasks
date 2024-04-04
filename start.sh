#!/bin/sh

set -e

echo "run the database migrations"
/app/migrate -path /app/migrations -database "$DB_SOURCE" -verbose up

echo "start the app"
exec "$@"