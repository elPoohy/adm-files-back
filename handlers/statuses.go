package handlers

import (
	"database/sql"
	"encoding/json"
	"github.com/jackc/pgx"
	"log"
	"net/http"
)

type Status struct {
	Code    int
	Message string
}

func ReturnError(w http.ResponseWriter, err error) {
	switch {
	case err == sql.ErrNoRows:
		StatusDBNotFound(err, w)
		return
	case err.(pgx.PgError).Code == "23505":
		StatusDBAlreadyExist(err, w)
		return
	default:
		StatusDBError(err, w)
		return
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
