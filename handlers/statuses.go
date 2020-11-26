package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

type Status struct {
	Code    int
	Message string
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
	err = json.NewEncoder(w).Encode(Status{
		Code:    http.StatusBadRequest,
		Message: "Bad incoming data",
	})
	if err != nil {
		log.Println(err)
	}
}
