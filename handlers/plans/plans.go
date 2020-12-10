package plans

import (
	"database/sql"
	"encoding/json"
	"files-back/dbase/dbplans"
	"files-back/handlers"
	"files-back/handlers/params"
	"github.com/go-playground/validator/v10"
	"net/http"
	"time"
)

var validate *validator.Validate

type IncomingStruct struct {
	Name        string  `json:"name" validate:"required,alphanum,min=2,max=15,lowercase"`
	DomainName  string  `json:"domainName" validate:"required"`
	FromDate    *string `json:"fromDate" validate:"omitempty,datetime=2006-01-02"`
	DueDate     *string `json:"dueDate" validate:"omitempty,datetime=2006-01-02"`
	Type        string  `json:"type"  validate:"required,oneof=personal group"`
	Description *string `json:"description"`
}

func (NewPlan *IncomingStruct) toDB() *dbplans.DBStruct {
	now := time.Now()
	planResponse := dbplans.DBStruct{
		Name:        NewPlan.Name,
		DomainName:  &NewPlan.DomainName,
		Type:        &NewPlan.Type,
		FromDate:    &now,
		Description: NewPlan.Description,
	}
	if NewPlan.FromDate != nil {
		FromDate, err := time.Parse("2006-01-02", *NewPlan.FromDate)
		if err == nil {
			planResponse.FromDate = &FromDate
		}
	}
	if NewPlan.DueDate != nil {
		DueDate, err := time.Parse("2006-01-02", *NewPlan.DueDate)
		if err == nil {
			planResponse.DueDate = &DueDate
		}
	}
	return &planResponse
}

func Get(w http.ResponseWriter, r *http.Request) {
	plans, err := dbplans.Query(params.GetQueryParams(r))
	if err != nil {
		handlers.StatusBadData(err, w)
		return
	}
	if len(plans) == 1 {
		params.ResponseJSON(w, plans[0])
	} else {
		params.ResponseJSON(w, plans)
	}
}

func Create(w http.ResponseWriter, r *http.Request) {
	n, err := extract(r)
	if err != nil {
		handlers.StatusBadData(err, w)
		return
	}
	err = dbplans.Insert(n.toDB())
	if err != nil {
		handlers.ReturnError(w, err)
		return
	}
	responseDomain, err := dbplans.Query(params.QueryParams{DomainName: &n.DomainName, PlanName: &n.Name})
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
	err = dbplans.Update(n.toDB(), params.GetQueryParams(r))
	if err != nil {
		handlers.StatusBadData(err, w)
		return
	}
	responseDomain, err := dbplans.Query(params.QueryParams{DomainName: &n.DomainName, PlanName: &n.Name})
	if err != nil {
		handlers.StatusBadData(err, w)
		return
	}
	params.ResponseJSON(w, responseDomain)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	err := dbplans.Delete(params.GetQueryParams(r))
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
	params.ResponseJSON(w, handlers.Status{
		Code:    200,
		Message: "Deleted",
	})
}

func extract(r *http.Request) (*IncomingStruct, error) {
	var NewPlan IncomingStruct
	err := json.NewDecoder(r.Body).Decode(&NewPlan)
	if err != nil {
		return nil, err
	}
	validate = validator.New()
	err = validate.Struct(&NewPlan)
	if err != nil {
		return nil, err
	}
	return &NewPlan, nil
}
