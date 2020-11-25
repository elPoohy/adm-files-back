package main

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

const defaultPort = "8080"

var defaultLimit = 10
var defaultOffset = 0

var validate *validator.Validate
var Db *sqlx.DB

func InitDB(dbuser, dbpwd, dbname, dbhost, dbport string) {
	dataSourceName := fmt.Sprintf("postgres://%v:%v@%v:%v/%v", dbuser, dbpwd, dbhost, dbport, dbname)

	db, err := sqlx.Connect("pgx", dataSourceName)
	if err != nil {
		log.Printf("Unable to connect to database: %v", err)
		os.Exit(1)
	}
	Db = db
}

// init is invoked before main()
func init() {
	// loads values from .env into the system
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

	InitDB(dbuser, dbpwd, dbname, dbhost, dbport)

	validate = validator.New()
	router := mux.NewRouter()

	router.HandleFunc("/domains", getDomains).Methods(http.MethodGet)
	router.HandleFunc("/domains/{domainName}", getDomain).Methods(http.MethodGet)
	router.HandleFunc("/domains", createDomain).Methods(http.MethodPost)
	router.HandleFunc("/domains/{domainName}", updateDomain).Methods(http.MethodPut)
	router.HandleFunc("/domains/{domainName}", deleteDomain).Methods(http.MethodDelete)

	log.Fatal(http.ListenAndServe(":"+port, router))
}
