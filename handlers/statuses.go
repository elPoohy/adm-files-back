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

func StatusDone(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(Status{
		Code:    http.StatusOK,
		Message: "Inserted",
	})
	if err != nil {
		log.Println(err)
	}
}

func StatusError(err error, w http.ResponseWriter) {
	log.Println(err)
	w.WriteHeader(http.StatusBadRequest)
	err = json.NewEncoder(w).Encode(
		Status{
			Code:    http.StatusBadRequest,
			Message: "Internal error",
		})
	if err != nil {
		log.Println(err)
	}
}

func StatusDBError(err error, w http.ResponseWriter) {
	log.Println(err)
	w.WriteHeader(http.StatusBadRequest)
	err = json.NewEncoder(w).Encode(
		Status{
			Code:    http.StatusBadRequest,
			Message: "Database error",
		})
	if err != nil {
		log.Println(err)
	}
}

func StatusDBAlreadyExist(err error, w http.ResponseWriter) {
	log.Println(err)
	w.WriteHeader(http.StatusBadRequest)
	err = json.NewEncoder(w).Encode(
		Status{
			Code:    http.StatusBadRequest,
			Message: "Already exist",
		})
	if err != nil {
		log.Println(err)
	}
}

func StatusDBNotFound(err error, w http.ResponseWriter) {
	log.Println(err)
	w.WriteHeader(http.StatusNotFound)
	err = json.NewEncoder(w).Encode(
		Status{
			Code:    http.StatusNotFound,
			Message: "Not found",
		})
	if err != nil {
		log.Println(err)
	}
}

func StatusBadData(err error, w http.ResponseWriter) {
	log.Println(err)
	w.WriteHeader(http.StatusBadRequest)
	err = json.NewEncoder(w).Encode(Status{
		Code:    http.StatusBadRequest,
		Message: "Bad incoming data",
	})
	if err != nil {
		log.Println(err)
	}
}

func StatusInvalidCredentials(err error, w http.ResponseWriter) {
	log.Println(err)
	w.WriteHeader(http.StatusUnauthorized)
	err = json.NewEncoder(w).Encode(Status{
		Code:    http.StatusUnauthorized,
		Message: "Invalid Credentials",
	})
	if err != nil {
		log.Println(err)
	}
}

func StatusUnauthorized(err error, w http.ResponseWriter) {
	log.Println(err)
	w.WriteHeader(http.StatusUnauthorized)
	err = json.NewEncoder(w).Encode(Status{
		Code:    http.StatusUnauthorized,
		Message: "Unauthorized",
	})
	if err != nil {
		log.Println(err)
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
