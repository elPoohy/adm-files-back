package database

import (
	"fmt"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"log"
	"os"
)

//var Db *sql.DB
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
