package main

import (
	"files-back/auth"
	"files-back/auth/directory"
	"files-back/dbase"
	"files-back/handlers/domains"
	"files-back/handlers/plans"
	"files-back/handlers/tenants"
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

	LDAPConnect()
	DBConnect()
	defer func() {
		err := dbase.DB.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	router := mux.NewRouter()
	router.HandleFunc("/login", auth.Login).Methods(http.MethodGet)
	domainsHandlers(router)
	plansHandlers(router)
	tenantsHandlers(router)

	log.Panic(http.ListenAndServe(":"+port, router))
}

func domainsHandlers(router *mux.Router) {
	router.Handle("/domains", auth.Middleware(http.HandlerFunc(domains.Get))).Methods(http.MethodGet)
	router.Handle("/domains/{domainName}", auth.Middleware(http.HandlerFunc(domains.Get))).Methods(http.MethodGet)
	router.Handle("/domains", auth.Middleware(http.HandlerFunc(domains.Create))).Methods(http.MethodPost)
	router.Handle("/domains/{domainName}", auth.Middleware(http.HandlerFunc(domains.Update))).Methods(http.MethodPut)
	router.Handle("/domains/{domainName}", auth.Middleware(http.HandlerFunc(domains.Delete))).Methods(http.MethodDelete)
}

func plansHandlers(router *mux.Router) {
	router.Handle("/plans", auth.Middleware(http.HandlerFunc(plans.Get))).Methods(http.MethodGet)
	router.Handle("/domains/{domainName}/plans", auth.Middleware(http.HandlerFunc(plans.Get))).Methods(http.MethodGet)
	router.Handle("/domains/{domainName}/plans/{planName}", auth.Middleware(http.HandlerFunc(plans.Get))).Methods(http.MethodGet)
	router.Handle("/domains/{domainName}/plans", auth.Middleware(http.HandlerFunc(plans.Create))).Methods(http.MethodPost)
	router.Handle("/domains/{domainName}/plans/{planName}", auth.Middleware(http.HandlerFunc(plans.Update))).Methods(http.MethodPut)
	router.Handle("/domains/{domainName}/plans/{planName}", auth.Middleware(http.HandlerFunc(plans.Delete))).Methods(http.MethodDelete)
}

func tenantsHandlers(router *mux.Router) {
	router.Handle("/tenants", auth.Middleware(http.HandlerFunc(tenants.Get))).Methods(http.MethodGet)
	router.Handle("/tenants", auth.Middleware(http.HandlerFunc(tenants.Get))).Methods(http.MethodGet)
	router.Handle("/tenants", auth.Middleware(http.HandlerFunc(tenants.Create))).Methods(http.MethodPost)
	router.Handle("/tenants/{tenantName}", auth.Middleware(http.HandlerFunc(tenants.Get))).Methods(http.MethodGet)
	router.Handle("/tenants/{tenantName}", auth.Middleware(http.HandlerFunc(tenants.Update))).Methods(http.MethodPut)
	router.Handle("/tenants/{tenantName}", auth.Middleware(http.HandlerFunc(tenants.Delete))).Methods(http.MethodDelete)
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
	directory.LDAPPassword = os.Getenv("BINDPASSWORD")
	directory.LDAPUsername = os.Getenv("BINDUSERNAME")
	directory.BaseDN = os.Getenv("BASEDN")
	directory.LDAPServer = os.Getenv("BINDADDRESS")
}
