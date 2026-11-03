package postgres

import (
	"context"
	"database/sql"
	"errors"
	"net"
	"net/url"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type DB struct{ SQL *sql.DB }

func ConnectSQL(host, port, user, pass, dbname string) (*DB, error) {
	if host == "" || port == "" || user == "" || pass == "" || dbname == "" {
		return nil, errors.New("database configuration is incomplete")
	}
	sslmode := os.Getenv("DB_SSLMODE")
	if sslmode == "" {
		sslmode = "verify-full"
	}
	if sslmode != "verify-full" && sslmode != "disable" {
		return nil, errors.New("DB_SSLMODE must be verify-full or disable for a local isolated database")
	}
	u := &url.URL{Scheme: "postgresql", Host: net.JoinHostPort(host, port), User: url.UserPassword(user, pass), Path: "/" + dbname}
	q := url.Values{"sslmode": {sslmode}, "connect_timeout": {"5"}}
	u.RawQuery = q.Encode()
	db, err := sql.Open("postgres", u.String())
	if err != nil {
		return nil, errors.New("cannot open database")
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, errors.New("cannot connect to database")
	}
	return &DB{SQL: db}, nil
}
