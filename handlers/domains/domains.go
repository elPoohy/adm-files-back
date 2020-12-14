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

func (newDomain *IncomingStruct) toDB() *dbdomains.DBStruct {
	domainResponse := dbdomains.DBStruct{
		Password:     &newDomain.Password,
		Description:  &newDomain.Description,
		Name:         newDomain.Name,
		Organisation: &newDomain.Organisation,
		PrimaryURL:   &newDomain.PrimaryURL,
		AdminURL:     &newDomain.AdminURL,
		DataPath:     &newDomain.DataPath,
		UserName:     &newDomain.UserName,
		Type:         &newDomain.Type,
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
		handlers.ResponseJSON(w, domains[0])
	} else {
		handlers.ResponseJSON(w, domains)
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
		handlers.ReturnError(w, err)
		return
	}
	handlers.StatusInserted(w)
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
	handlers.ResponseJSON(w, responseDomain)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	err := dbdomains.Delete(params.GetQueryParams(r))
	if err != nil {
		handlers.StatusBadData(err, w)
		return
	}
	handlers.StatusDeleted(w)
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
