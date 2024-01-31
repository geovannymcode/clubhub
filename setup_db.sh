#!/bin/bash
set -e

# Parámetros de conexión a PostgreSQL
export PGPASSWORD=postgres
HOST="clubhub-db"
DB="clubhub"
USER="postgres"
PORT="5432"

# Conectándose a PostgreSQL y ejecutando comandos SQL
psql -h $HOST -U $USER -p $PORT -d $DB <<-EOSQL
    CREATE USER postgres WITH PASSWORD 'postgres';
    CREATE DATABASE clubhub;
    GRANT ALL PRIVILEGES ON DATABASE clubhub TO clubhub;
EOSQL