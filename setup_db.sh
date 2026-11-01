#!/bin/sh
set -eu
# PostgreSQL's image creates DB_NAME on first startup. Never embed credentials.
: "${DB_HOST:?Set DB_HOST}"
: "${DB_PORT:?Set DB_PORT}"
: "${DB_USER:?Set DB_USER}"
: "${DB_PASSWORD:?Set DB_PASSWORD}"
: "${DB_NAME:?Set DB_NAME}"
export PGPASSWORD="$DB_PASSWORD"
exec psql -v ON_ERROR_STOP=1 -h "$DB_HOST" -U "$DB_USER" -p "$DB_PORT" -d "$DB_NAME" -c 'SELECT 1'
