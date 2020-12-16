package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx"
	"log"
	"net/http"
)

func responseError(w http.ResponseWriter, err error, resp interface{}) {
	log.Println(err)
	ResponseJSON(w, resp)
}

func ReturnError(w http.ResponseWriter, err error) {
	var pgError *pgx.PgError
	switch {
	case errors.Is(err, sql.ErrNoRows):
		StatusDBNotFound(err, w)
		return
	case errors.As(err, &pgError):
		switch pgError.Code {
		case "42P01":
			StatusDBError(err, w)
			return
		case "23505":
			StatusDBAlreadyExist(err, w)
			return
		}
	default:
		StatusDBError(err, w)
	}
}

func ResponseJSON(w http.ResponseWriter, resp interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Println(err)
	}
}
