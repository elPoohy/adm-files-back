package handlers

import (
	"net/http"
)

type Status struct {
	Code    int
	Message string
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
