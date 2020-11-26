package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

var defaultLimit = 10
var defaultOffset = 0

type QueryParams struct {
	Limit  int    `db:"limit"`
	Offset int    `db:"offset"`
	Search string `db:"search"`
}

func GetQueryParams(URL *url.URL) QueryParams {
	response := QueryParams{
		Limit:  defaultLimit,
		Offset: defaultOffset,
		Search: "",
	}
	limitString := URL.Query().Get("limit")
	temp, err := strconv.Atoi(limitString)
	if err == nil {
		response.Limit = temp
	}
	offsetString := URL.Query().Get("offset")
	temp, err = strconv.Atoi(offsetString)
	if err == nil {
		response.Offset = temp
	}
	response.Search = "%" + URL.Query().Get("search") + "%"
	return response
}

func ResponseJSON(w http.ResponseWriter, err error, domains interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(domains)
	if err != nil {
		log.Println(err)
	}
}
