package params

import (
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

var (
	defaultLimit  = 10
	defaultOffset = 0
	deleted       = "deleted"
	disabled      = "disabled"
)

const incomingTrue = "true"

type QueryParams struct {
	Limit        *int    `db:"limit"`
	Offset       *int    `db:"offset"`
	Search       *string `db:"search"`
	DeleteType   *string `db:"delete"`
	ShowDeleted  bool
	ShowDisabled bool
	DomainName   *string `db:"domain_name"`
	TenantName   *string `db:"tenant_name"`
	PlanName     *string `db:"plan_name"`
}

func GetQueryParams(r *http.Request) QueryParams {
	resp := QueryParams{
		Limit:        getLimit(r),
		Offset:       getOffset(r),
		Search:       getSearchLine(r),
		DomainName:   getDomain(r),
		TenantName:   getTenant(r),
		PlanName:     getPlan(r),
		DeleteType:   getDeleteType(r),
		ShowDeleted:  getDeleted(r),
		ShowDisabled: getDisabled(r),
	}
	return resp
}

func getSearchLine(r *http.Request) *string {
	switch temp := r.URL.Query().Get("search"); temp {
	case "":
		return nil
	default:
		temp = "%" + temp + "%"
		return &temp
	}
}

func getLimit(r *http.Request) *int {
	resp, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		return &defaultLimit
	}
	return &resp
}

func getOffset(r *http.Request) *int {
	resp, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		return &defaultOffset
	}
	return &resp
}

func getDomain(r *http.Request) *string {
	respURL := mux.Vars(r)["domainName"]
	respQuery := r.URL.Query().Get("domainName")
	switch {
	case respURL != "":
		return &respURL
	case respQuery != "":
		return &respQuery
	default:
		return nil
	}
}

func getTenant(r *http.Request) *string {
	switch resp := mux.Vars(r)["tenantName"]; resp {
	case "":
		return nil
	default:
		return &resp
	}
}

func getPlan(r *http.Request) *string {
	switch resp := mux.Vars(r)["planName"]; resp {
	case "":
		return nil
	default:
		return &resp
	}
}

func getDeleteType(r *http.Request) *string {
	switch r.URL.Query().Get("forced") {
	case incomingTrue:
		return &deleted
	default:
		return &disabled
	}
}

func getDeleted(r *http.Request) bool {
	switch r.URL.Query().Get("deleted") {
	case incomingTrue:
		return true
	default:
		return false
	}
}

func getDisabled(r *http.Request) bool {
	switch r.URL.Query().Get("deleted") {
	case incomingTrue:
		return true
	default:
		return false
	}
}
