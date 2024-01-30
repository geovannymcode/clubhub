package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	ht "github.com/Geovanny0401/clubhub/internal/adapter/handler/http"
	"github.com/Geovanny0401/clubhub/internal/adapter/storage/postgres"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

func main() {
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPort := os.Getenv("DB_PORT")

	connection, err := postgres.ConnectSQL(dbHost, dbPort, dbUser, "", dbName)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	defer connection.SQL.Close()

	initDataBase(connection.SQL)

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(setDefaultHeaders().Handler)

	dHandler := ht.NewServerHandler(connection)
	router.Route("/", func(rt chi.Router) {
		rt.Mount("/scraping", domainRouter(dHandler))
	})

	http.ListenAndServe(":8005", router)
}

func domainRouter(dHandler *ht.Domain) http.Handler {
	r := chi.NewRouter()
	r.Get("/", dHandler.GetAllAddress)
	r.Get("/address={address}", dHandler.GetByAddress)

	return r
}

func setDefaultHeaders() *cors.Cors {
	headers := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	return headers
}

func initDataBase(connection *sql.DB) {

	if _, err := connection.Exec(
		"CREATE TABLE IF NOT EXISTS domain (id SERIAL PRIMARY KEY, address varchar(100) NOT NULL, last_consult TIMESTAMP NOT NULL)"); err != nil {
		log.Println(err)
	}

	if _, err := connection.Exec(
		"CREATE TABLE IF NOT EXISTS detail_domain (id serial PRIMARY KEY, id_domain serial NOT NULL, ipaddress varchar(100) NOT NULL, servername varchar(200) NULL, grade varchar(10) NOT NULL, date TIMESTAMP NOT NULL, CONSTRAINT detail_domain_fk FOREIGN KEY (id_domain) REFERENCES domain(id))"); err != nil {
		log.Println(err)
	}
}
