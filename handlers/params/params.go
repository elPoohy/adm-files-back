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

const (
	incomingTrue = "true"
	noData       = ""
)

type QueryParams struct {
	Limit        *int    `db:"limit"`
	Offset       *int    `db:"offset"`
	Search       *string `db:"search"`
	DeleteType   *string `db:"delete"`
	ShowDeleted  bool
	ShowDisabled bool
	Email        *string `db:"email"`
	DomainName   *string `db:"domain_name"`
	TenantName   *string `db:"tenant_name"`
	PlanName     *string `db:"plan_name"`
	TariffName   *string `db:"tariff_name"`
	GroupName    *string `db:"group_name"`
}

func GetQueryParams(r *http.Request) QueryParams {
	resp := QueryParams{
		Limit:        getLimit(r),
		Offset:       getOffset(r),
		Search:       getSearchLine(r),
		Email:        getEmail(r),
		DomainName:   getDomain(r),
		TenantName:   getTenant(r),
		GroupName:    getGroup(r),
		PlanName:     getPlan(r),
		TariffName:   getTariff(r),
		DeleteType:   getDeleteType(r),
		ShowDeleted:  getDeleted(r),
		ShowDisabled: getDisabled(r),
	}
	return resp
}

func getSearchLine(r *http.Request) *string {
	switch temp := r.URL.Query().Get("search"); temp {
	case noData:
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
	switch resp := mux.Vars(r)["domainName"]; resp {
	case noData:
		return nil
	default:
		return &resp
	}
}

func getTariff(r *http.Request) *string {
	switch resp := mux.Vars(r)["tariffName"]; resp {
	case noData:
		return nil
	default:
		return &resp
	}
}

func getTenant(r *http.Request) *string {
	switch resp := mux.Vars(r)["tenantName"]; resp {
	case noData:
		return nil
	default:
		return &resp
	}
}

func getGroup(r *http.Request) *string {
	switch resp := mux.Vars(r)["groupName"]; resp {
	case noData:
		return nil
	default:
		return &resp
	}
}

func getEmail(r *http.Request) *string {
	switch resp := mux.Vars(r)["email"]; resp {
	case noData:
		return nil
	default:
		return &resp
	}
}

func getPlan(r *http.Request) *string {
	switch resp := mux.Vars(r)["planName"]; resp {
	case noData:
		return nil
	default:
		return &resp
	}
}

func getDeleteType(r *http.Request) *string {
	if r.URL.Query().Get("forced") == incomingTrue {
		return &deleted
	}
	return &disabled
}

func getDeleted(r *http.Request) bool {
	return r.URL.Query().Get("deleted") == incomingTrue
}

func getDisabled(r *http.Request) bool {
	return r.URL.Query().Get("disabled") == incomingTrue
}
