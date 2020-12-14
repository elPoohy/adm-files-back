package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx"
	"log"
	"net/http"
)

type Status struct {
	Code    int
	Message string
}

func ReturnError(w http.ResponseWriter, err error) {
	var pgError *pgx.PgError
	switch {
	case errors.Is(err, sql.ErrNoRows):
		StatusDBNotFound(err, w)
		return
	case errors.As(err, &pgError):
		if pgError.Code == "23505" {
			StatusDBAlreadyExist(err, w)
			return
		}
	}
	StatusDBError(err, w)
}

func StatusDeleted(w http.ResponseWriter) {
	ResponseJSON(w, Status{
		Code:    http.StatusOK,
		Message: "Deleted",
	})
}

func StatusInserted(w http.ResponseWriter) {
	ResponseJSON(w, Status{
		Code:    http.StatusOK,
		Message: "Inserted",
	})
}

func StatusError(err error, w http.ResponseWriter) {
	responseError(w, err, Status{
		Code:    http.StatusBadRequest,
		Message: "Internal error",
	})
}

func StatusDBError(err error, w http.ResponseWriter) {
	responseError(w, err, Status{
		Code:    http.StatusBadRequest,
		Message: "Database error",
	})
}

func StatusDBAlreadyExist(err error, w http.ResponseWriter) {
	responseError(w, err, Status{
		Code:    http.StatusBadRequest,
		Message: "Already exist",
	})
}

func StatusDBNotFound(err error, w http.ResponseWriter) {
	responseError(w, err, Status{
		Code:    http.StatusNotFound,
		Message: "Not found",
	})
}

func StatusBadData(err error, w http.ResponseWriter) {
	responseError(w, err, Status{
		Code:    http.StatusBadRequest,
		Message: "Bad incoming data",
	})
}

func StatusInvalidCredentials(err error, w http.ResponseWriter) {
	responseError(w, err, Status{
		Code:    http.StatusUnauthorized,
		Message: "Invalid Credentials",
	})
}

func StatusUnauthorized(err error, w http.ResponseWriter) {
	responseError(w, err, Status{
		Code:    http.StatusUnauthorized,
		Message: "Unauthorized",
	})
}

func responseError(w http.ResponseWriter, err error, resp interface{}) {
	log.Println(err)
	ResponseJSON(w, resp)
}

func ResponseJSON(w http.ResponseWriter, resp interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Println(err)
	}
}
