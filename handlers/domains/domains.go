package domains

import (
	"database/sql"
	"encoding/json"
	"files-back/dbase/dbdomain"
	"files-back/handlers"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"net/http"
)

var validate *validator.Validate

type IncomingStruct struct {
	Name         string `json:"name" validate:"required,alphanum,min=2,max=15,lowercase"`
	Organisation string `json:"organisation" validate:"required"`
	PrimaryURL   string `json:"primaryUrl" validate:"required"`
	AdminURL     string `json:"adminUrl" validate:"required"`
	DataPath     string `json:"data_path"`
	UserName     string `json:"user_name" validate:"required,alphanum,min=2,max=15,lowercase"`
	Password     string `json:"password" validate:"required,min=12"`
	Type         string `json:"type" validate:"required,oneof=primary wholesale premium"`
	Description  string `json:"description" validate:"required"`
}

func (NewDomain *IncomingStruct) toDB() *dbdomain.DBStruct {
	domainResponse := dbdomain.DBStruct{
		Password:     NewDomain.Password,
		Description:  &NewDomain.Description,
		Name:         NewDomain.Name,
		Organisation: NewDomain.Organisation,
		PrimaryURL:   NewDomain.PrimaryURL,
		AdminURL:     NewDomain.AdminURL,
		DataPath:     NewDomain.DataPath,
		UserName:     NewDomain.UserName,
		Type:         NewDomain.Type,
	}
	return &domainResponse
}

func GetAll(w http.ResponseWriter, r *http.Request) {
	domains, err := dbdomain.QueryAll(handlers.GetQueryParams(r.URL))
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			handlers.StatusDBNotFound(err, w)
			return
		default:
			handlers.StatusDBError(err, w)
			return
		}
	}
	handlers.ResponseJSON(w, err, domains)
}

func GetOne(w http.ResponseWriter, r *http.Request) {
	responseDomain, err := dbdomain.QueryOne(extractName(r))
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			handlers.StatusDBNotFound(err, w)
			return
		default:
			handlers.StatusDBError(err, w)
			return
		}
	}
	handlers.ResponseJSON(w, err, responseDomain)
}

func Create(w http.ResponseWriter, r *http.Request) {
	newDomain, err := extract(r)
	if err != nil {
		handlers.StatusBadData(err, w)
	}
	err = dbdomain.Insert(newDomain.toDB())
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			handlers.StatusDBNotFound(err, w)
			return
		default:
			handlers.StatusDBError(err, w)
			return
		}
	}
	responseDomain, err := dbdomain.QueryOne(newDomain.Name)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			handlers.StatusDBNotFound(err, w)
			return
		default:
			handlers.StatusDBError(err, w)
			return
		}
	}
	handlers.ResponseJSON(w, err, responseDomain)
}

func Update(w http.ResponseWriter, r *http.Request) {
	newDomain, err := extract(r)
	if err != nil {
		handlers.StatusBadData(err, w)
	}
	err = dbdomain.Update(newDomain.toDB(), extractName(r))
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			handlers.StatusDBNotFound(err, w)
			return
		default:
			handlers.StatusDBError(err, w)
			return
		}
	}
	responseDomain, err := dbdomain.QueryOne(newDomain.Name)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			handlers.StatusDBNotFound(err, w)
			return
		default:
			handlers.StatusDBError(err, w)
			return
		}
	}
	handlers.ResponseJSON(w, err, responseDomain)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	err := dbdomain.Delete(extractName(r))
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			handlers.StatusDBNotFound(err, w)
			return
		default:
			handlers.StatusDBError(err, w)
			return
		}
	}
	handlers.ResponseJSON(w, err, handlers.Status{
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

func extractName(r *http.Request) string {
	return mux.Vars(r)["domainName"]
}
