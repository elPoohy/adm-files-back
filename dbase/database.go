package dbase

import (
	"fmt"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"log"
	"os"
)

var (
	DB *sqlx.DB
)

func InitDB(dbuser, dbpwd, dbname, dbhost, dbport string) {
	dataSourceName := fmt.Sprintf("postgres://%v:%v@%v:%v/%v", dbuser, dbpwd, dbhost, dbport, dbname)

	db, err := sqlx.Connect("pgx", dataSourceName)
	if err != nil {
		log.Printf("Unable to connect to dbase: %v", err)
		os.Exit(1)
	}
	DB = db
}
