package main

import (
	"files-back/handlers/domains"

	"files-back/dbase"
	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

const defaultPort = "8080"

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	dbuser := os.Getenv("DBUSER")
	dbpwd := os.Getenv("DBPWD")
	dbname := os.Getenv("DB")
	dbhost := os.Getenv("DBHOST")
	dbport := os.Getenv("DBPORT")

	dbase.InitDB(dbuser, dbpwd, dbname, dbhost, dbport)

	router := mux.NewRouter()

	router.HandleFunc("/domains", domains.GetAll).Methods(http.MethodGet)
	router.HandleFunc("/domains/{domainName}", domains.GetOne).Methods(http.MethodGet)
	router.HandleFunc("/domains", domains.Create).Methods(http.MethodPost)
	router.HandleFunc("/domains/{domainName}", domains.Update).Methods(http.MethodPut)
	router.HandleFunc("/domains/{domainName}", domains.Delete).Methods(http.MethodDelete)

	log.Fatal(http.ListenAndServe(":"+port, router))
}
