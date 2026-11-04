package main

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	ht "github.com/Geovanny0401/clubhub/internal/adapter/handler/http"
	"github.com/Geovanny0401/clubhub/internal/adapter/storage/postgres"
	"github.com/Geovanny0401/clubhub/internal/security"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// @title					Go CLubhub API
// @version					1.0
// @summary                 Clubhub
// @description             This is a simple technical test for senior backend position at clubhub
// @contact.name			Geovanny Mendoza
// @contact.url				https://github.com/Geovanny0401/clubhub
// @contact.email			me@geovannycode.com
// @BasePath                /
func main() {
	token := os.Getenv("API_TOKEN")
	if len(token) < 32 || strings.TrimSpace(token) != token {
		log.Fatal("API_TOKEN must contain at least 32 characters and no surrounding whitespace")
	}

	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbPort := os.Getenv("DB_PORT")

	connection, err := postgres.ConnectSQL(dbHost, dbPort, dbUser, dbPassword, dbName)
	if err != nil {
		log.Fatal("database connection failed")
	}
	defer connection.SQL.Close()

	if err := initDataBase(connection.SQL); err != nil {
		log.Fatal("database initialization failed")
	}

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(25 * time.Second))

	dHandler := ht.NewServerHandler(connection)
	router.Route("/", func(rt chi.Router) {
		rt.Mount("/clubhub", domainRouter(dHandler))
	})

	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = "127.0.0.1:8001"
	}
	server := &http.Server{
		Addr: addr, Handler: security.Protect(token, router),
		ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 10 * time.Second,
		WriteTimeout: 30 * time.Second, IdleTimeout: 60 * time.Second,
		MaxHeaderBytes: 16 << 10,
	}
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal("HTTP server failed")
	}

}

func domainRouter(dHandler *ht.Domain) http.Handler {
	r := chi.NewRouter()
	r.Get("/", dHandler.GetAllAddress)
	r.Get("/address={address}", dHandler.GetByAddress)

	return r
}

func initDataBase(connection *sql.DB) error {

	if _, err := connection.Exec(
		"CREATE TABLE IF NOT EXISTS domain (id SERIAL PRIMARY KEY, address varchar(100) NOT NULL, last_consult TIMESTAMP NOT NULL)"); err != nil {
		return err
	}

	if _, err := connection.Exec(
		"CREATE TABLE IF NOT EXISTS detail_domain (id serial PRIMARY KEY, id_domain serial NOT NULL, ipaddress varchar(100) NOT NULL, servername varchar(200) NULL, grade varchar(10) NOT NULL, date TIMESTAMP NOT NULL, CONSTRAINT detail_domain_fk FOREIGN KEY (id_domain) REFERENCES domain(id))"); err != nil {
		return err
	}
	return nil
}
