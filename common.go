package main

import (
	"net/url"
	"strconv"
)

type QueryParams struct {
	Limit  int    `db:"limit"`
	Offset int    `db:"offset"`
	Search string `db:"search"`
}

func getQueryParams(URL *url.URL) QueryParams {
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
