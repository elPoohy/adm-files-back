package tenants

import (
	"encoding/json"
	"files-back/dbase/dbdomains"
	"files-back/dbase/dbplans"
	"files-back/dbase/dbtenants"
	"files-back/handlers"
	"files-back/handlers/params"
	"github.com/go-playground/validator/v10"
	"net/http"
)

var validate *validator.Validate

type IncomingStruct struct {
	Name         string  `json:"name" validate:"required,alphanum,min=2,max=15,lowercase"`
	Organisation string  `json:"organisation" validate:"required"`
	OrderForm    string  `json:"order_form" validate:"required"`
	OrderLink    string  `json:"order_link" validate:"required,url"`
	Type         string  `json:"type" validate:"required,oneof=primary regular premium"`
	Plan         string  `json:"planName" validate:"required,alphanum,min=2,max=15,lowercase"`
	Domain       string  `json:"domainName" validate:"required,alphanum,min=2,max=15,lowercase"`
	Description  *string `json:"description"`
}

func (newTenant *IncomingStruct) toDB() *dbtenants.DBStruct {
	tenantResponse := dbtenants.DBStruct{
		Description:  newTenant.Description,
		Name:         newTenant.Name,
		Organisation: &newTenant.Organisation,
		OrderForm:    &newTenant.OrderForm,
		OrderLink:    &newTenant.OrderLink,
		Type:         &newTenant.Type,
		Plan: &dbplans.DBStruct{
			Name: newTenant.Plan,
		},
		Domain: &dbdomains.DBStruct{
			Name: newTenant.Domain,
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
		handlers.ResponseJSON(w, tenants[0])
	} else {
		handlers.ResponseJSON(w, tenants)
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
		handlers.ReturnError(w, err)
		return
	}
	handlers.StatusDone(w)
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
	handlers.StatusDone(w)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	err := dbtenants.Delete(params.GetQueryParams(r))
	if err != nil {
		handlers.StatusBadData(err, w)
		return
	}
	handlers.ResponseJSON(w, handlers.Status{
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
