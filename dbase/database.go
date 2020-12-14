package dbase

import (
	"database/sql"
	"fmt"
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

func AppendWhere(where string) string {
	if len(where) == 0 {
		where = "WHERE "
	} else {
		where += "AND "
	}
	return where
}

func ExecWithChekOne(data interface{}, sqlQuery string) error {
	tx, err := DB.Beginx()
	if err != nil {
		return err
	}
	result, err := tx.NamedExec(sqlQuery, data)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	if rows != 1 {
		_ = tx.Rollback()
		return sql.ErrNoRows
	}
	err = tx.Commit()
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	return nil
}
