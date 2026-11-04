\set ON_ERROR_STOP on
-- Read credentials without placing them in command arguments or source files.
\getenv app_password POSTGRES_APP_PASSWORD
\getenv admin_password POSTGRES_PASSWORD
BEGIN;
SELECT 'CREATE ROLE clubhub' WHERE NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'clubhub') \gexec
ALTER ROLE clubhub WITH LOGIN NOSUPERUSER NOCREATEDB NOCREATEROLE NOREPLICATION PASSWORD :'app_password';
ALTER ROLE postgres PASSWORD :'admin_password';
GRANT CONNECT ON DATABASE clubhub TO clubhub;
GRANT USAGE, CREATE ON SCHEMA public TO clubhub;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO clubhub;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO clubhub;
COMMIT;
