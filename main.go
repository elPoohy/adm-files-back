package main

import (
	"files-back/auth"
	"files-back/auth/directory"
	"files-back/dbase"
	"files-back/handlers/domains"
	"files-back/handlers/plans"
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
	auth.SecretKey = []byte(os.Getenv("SECRET"))

	DBConnect()
	defer dbase.DB.Close()

	LDAPConnect()
	defer directory.LDAP.Close()

	router := mux.NewRouter()

	router.HandleFunc("/login", auth.Login).Methods(http.MethodGet)

	domainsHandlers(router)
	plansHandlers(router)

	log.Fatal(http.ListenAndServe(":"+port, router))
}

func plansHandlers(router *mux.Router) {
	router.HandleFunc("/plans", plans.Get).Methods(http.MethodGet)
	router.HandleFunc("/domains/{domainName}/plans", plans.Get).Methods(http.MethodGet)
	router.HandleFunc("/domains/{domainName}/plans/{planName}", plans.Get).Methods(http.MethodGet)
	router.HandleFunc("/domains/{domainName}/plans", plans.Create).Methods(http.MethodPost)
	router.HandleFunc("/domains/{domainName}/plans/{planName}", plans.Update).Methods(http.MethodPut)
	router.HandleFunc("/domains/{domainName}/plans/{planName}", plans.Delete).Methods(http.MethodDelete)
}

func domainsHandlers(router *mux.Router) {
	router.HandleFunc("/domains", domains.Get).Methods(http.MethodGet)
	router.HandleFunc("/domains/{domainName}", domains.Get).Methods(http.MethodGet)
	router.HandleFunc("/domains", domains.Create).Methods(http.MethodPost)
	router.HandleFunc("/domains/{domainName}", domains.Update).Methods(http.MethodPut)
	router.HandleFunc("/domains/{domainName}", domains.Delete).Methods(http.MethodDelete)
}

func DBConnect() {
	dbuser := os.Getenv("DBUSER")
	dbpwd := os.Getenv("DBPWD")
	dbname := os.Getenv("DB")
	dbhost := os.Getenv("DBHOST")
	dbport := os.Getenv("DBPORT")
	dbase.InitDB(dbuser, dbpwd, dbname, dbhost, dbport)
}

func LDAPConnect() {
	directory.BindPassword = os.Getenv("BINDPASSWORD")
	directory.BindUsername = os.Getenv("BINDUSERNAME")
	directory.BaseDN = os.Getenv("BASEDN")
	ldapAddr := os.Getenv("BINDADDRESS")
	ldapPort := os.Getenv("BINDPORT")
	directory.InitLDAP(ldapAddr, ldapPort)
}
