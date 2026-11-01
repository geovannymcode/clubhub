#!/bin/sh
set -eu
: "${POSTGRES_APP_PASSWORD:?Set POSTGRES_APP_PASSWORD}"
: "${POSTGRES_PASSWORD:?Set POSTGRES_PASSWORD}"
# The local socket is used inside the database container. Values are SQL-quoted
# by psql; no password is interpolated into SQL or written to a log.
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" --set=app_password="$POSTGRES_APP_PASSWORD" --set=admin_password="$POSTGRES_PASSWORD" <<'SQL'
SELECT 'CREATE ROLE clubhub' WHERE NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'clubhub') \gexec
ALTER ROLE clubhub WITH LOGIN NOSUPERUSER NOCREATEDB NOCREATEROLE NOREPLICATION PASSWORD :'app_password';
ALTER ROLE postgres PASSWORD :'admin_password';
GRANT CONNECT ON DATABASE clubhub TO clubhub;
GRANT USAGE, CREATE ON SCHEMA public TO clubhub;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO clubhub;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO clubhub;
SQL
