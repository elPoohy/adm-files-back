package domains

import (
	"encoding/json"
	"files-back/dbase/dbdomains"
	"files-back/handlers"
	"files-back/handlers/params"
	"github.com/go-playground/validator/v10"
	"net/http"
)

var validate *validator.Validate

type IncomingStruct struct {
	Name         string `json:"name" validate:"required,alphanum,min=2,max=15,lowercase"`
	Organisation string `json:"organisation" validate:"required"`
	PrimaryURL   string `json:"primaryUrl" validate:"required,url"`
	AdminURL     string `json:"adminUrl" validate:"required,url"`
	DataPath     string `json:"data_path" validate:"required"`
	UserName     string `json:"user_name" validate:"required,alphanum,min=2,max=15,lowercase"`
	Password     string `json:"password" validate:"required,min=12"`
	Type         string `json:"type" validate:"required,oneof=primary wholesale premium"`
	Description  string `json:"description"`
}

func (NewDomain *IncomingStruct) toDB() *dbdomains.DBStruct {
	domainResponse := dbdomains.DBStruct{
		Password:     &NewDomain.Password,
		Description:  &NewDomain.Description,
		Name:         NewDomain.Name,
		Organisation: &NewDomain.Organisation,
		PrimaryURL:   &NewDomain.PrimaryURL,
		AdminURL:     &NewDomain.AdminURL,
		DataPath:     &NewDomain.DataPath,
		UserName:     &NewDomain.UserName,
		Type:         &NewDomain.Type,
	}
	return &domainResponse
}

func Get(w http.ResponseWriter, r *http.Request) {
	domains, err := dbdomains.Query(params.GetQueryParams(r))
	if err != nil {
		handlers.StatusBadData(err, w)
		return
	}
	if len(domains) == 1 {
		params.ResponseJSON(w, domains[0])
	} else {
		params.ResponseJSON(w, domains)
	}
}

func Create(w http.ResponseWriter, r *http.Request) {
	n, err := extract(r)
	if err != nil {
		handlers.StatusBadData(err, w)
		return
	}
	err = dbdomains.Insert(n.toDB())
	if err != nil {
		handlers.StatusBadData(err, w)
		return
	}
	responseDomain, err := dbdomains.Query(params.QueryParams{DomainName: &n.Name})
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
	err = dbdomains.Update(n.toDB(), params.GetQueryParams(r))
	if err != nil {
		handlers.StatusBadData(err, w)
		return
	}
	responseDomain, err := dbdomains.Query(params.QueryParams{DomainName: &n.Name})
	if err != nil {
		handlers.StatusBadData(err, w)
		return
	}
	params.ResponseJSON(w, responseDomain)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	err := dbdomains.Delete(params.GetQueryParams(r))
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
	var NewDomain IncomingStruct
	err := json.NewDecoder(r.Body).Decode(&NewDomain)
	if err != nil {
		return nil, err
	}
	validate = validator.New()
	err = validate.Struct(&NewDomain)
	if err != nil {
		return nil, err
	}
	return &NewDomain, nil
}
