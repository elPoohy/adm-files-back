package tenants

import (
	"encoding/json"
	"files-back/dbase/dbplans"
	"files-back/dbase/dbtenants"
	"files-back/handlers"
	"files-back/handlers/params"
	"github.com/go-playground/validator/v10"
	"net/http"
)

var validate *validator.Validate

type IncomingStruct struct {
	Name         string `json:"name" validate:"required,alphanum,min=2,max=15,lowercase"`
	Organisation string `json:"organisation" validate:"required"`
	OrderForm    string `json:"order_form" validate:"required"`
	OrderLink    string `json:"order_link" validate:"required,url"`
	Type         string `json:"type" validate:"required,oneof=primary regular premium"`
	Plan         string `json:"planName" validate:"required"`
	Description  string `json:"description"`
}

func (NewTenant *IncomingStruct) toDB() *dbtenants.DBStruct {
	tenantResponse := dbtenants.DBStruct{
		Description:  &NewTenant.Description,
		Name:         NewTenant.Name,
		Organisation: &NewTenant.Organisation,
		OrderForm:    &NewTenant.OrderForm,
		OrderLink:    &NewTenant.OrderLink,
		Type:         &NewTenant.Type,
		Plan: &dbplans.DBStruct{
			Name: NewTenant.Plan,
		},
	}
	return &tenantResponse
}

func Get(w http.ResponseWriter, r *http.Request) {
	tenants, err := dbtenants.Query(params.GetQueryParams(r))
	if err != nil {
		handlers.StatusBadData(err, w)
		return
	}
	if len(tenants) == 1 {
		params.ResponseJSON(w, tenants[0])
	} else {
		params.ResponseJSON(w, tenants)
	}
}

func Create(w http.ResponseWriter, r *http.Request) {
	n, err := extract(r)
	if err != nil {
		handlers.StatusBadData(err, w)
		return
	}
	err = dbtenants.Insert(n.toDB())
	if err != nil {
		handlers.StatusBadData(err, w)
		return
	}
	responseDomain, err := dbtenants.Query(params.QueryParams{TenantName: &n.Name})
	if err != nil {
		handlers.StatusBadData(err, w)
		return
	}
	params.ResponseJSON(w, responseDomain)
}

func Update(w http.ResponseWriter, r *http.Request) {
	n, err := extract(r)
	if err != nil {
		handlers.StatusBadData(err, w)
	}
	err = dbtenants.Update(n.toDB(), params.GetQueryParams(r))
	if err != nil {
		handlers.StatusBadData(err, w)
		return
	}
	responseDomain, err := dbtenants.Query(params.QueryParams{TenantName: &n.Name})
	if err != nil {
		handlers.StatusBadData(err, w)
		return
	}
	params.ResponseJSON(w, responseDomain)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	err := dbtenants.Delete(params.GetQueryParams(r))
	if err != nil {
		handlers.StatusBadData(err, w)
		return
	}
	params.ResponseJSON(w, handlers.Status{
		Code:    200,
		Message: "Deleted",
	})
}

func extract(r *http.Request) (*IncomingStruct, error) {
	var NewTenant IncomingStruct
	err := json.NewDecoder(r.Body).Decode(&NewTenant)
	if err != nil {
		return nil, err
	}
	validate = validator.New()
	err = validate.Struct(&NewTenant)
	if err != nil {
		return nil, err
	}
	return &NewTenant, nil
}
