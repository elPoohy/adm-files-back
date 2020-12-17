package incoming

import (
	"encoding/json"
	"files-back/dbase/dbdomains"
	"files-back/dbase/dbgroups"
	"files-back/dbase/dbplans"
	"files-back/dbase/dbtariffs"
	"files-back/dbase/dbtenants"
	"files-back/dbase/dbusers"
	"files-back/handlers/params"
	"github.com/go-playground/validator/v10"
	"net/http"
	"time"
)

var validate *validator.Validate

func Extract(r *http.Request, new interface{}) error {
	err := json.NewDecoder(r.Body).Decode(new)
	if err != nil {
		return err
	}
	validate = validator.New()
	err = validate.Struct(new)
	if err != nil {
		return err
	}
	return nil
}

type Domain struct {
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

func (incoming *Domain) ToDB(r *http.Request) *dbdomains.DBStruct {
	p := params.GetQueryParams(r)
	res := dbdomains.DBStruct{
		Password:     &incoming.Password,
		Description:  &incoming.Description,
		Name:         incoming.Name,
		OldName:      p.DomainName,
		Organisation: &incoming.Organisation,
		PrimaryURL:   &incoming.PrimaryURL,
		AdminURL:     &incoming.AdminURL,
		DataPath:     &incoming.DataPath,
		UserName:     &incoming.UserName,
		Type:         &incoming.Type,
	}
	return &res
}

type Plan struct {
	Name        string  `json:"name" validate:"required,alphanum,min=2,max=15,lowercase"`
	DomainName  string  `json:"domainName" validate:"required"`
	FromDate    *string `json:"fromDate" validate:"omitempty,datetime=2006-01-02"`
	DueDate     *string `json:"dueDate" validate:"omitempty,datetime=2006-01-02"`
	Type        string  `json:"type"  validate:"required,oneof=personal group"`
	Description *string `json:"description"`
}

func (incoming *Plan) ToDB(r *http.Request) *dbplans.DBStruct {
	now := time.Now()
	p := params.GetQueryParams(r)
	resp := dbplans.DBStruct{
		Name:        incoming.Name,
		DomainName:  p.DomainName,
		Type:        &incoming.Type,
		FromDate:    &now,
		Description: incoming.Description,
		OldName:     p.PlanName,
	}
	if incoming.FromDate != nil {
		FromDate, err := time.Parse("2006-01-02", *incoming.FromDate)
		if err == nil {
			resp.FromDate = &FromDate
		}
	}
	if incoming.DueDate != nil {
		DueDate, err := time.Parse("2006-01-02", *incoming.DueDate)
		if err == nil {
			resp.DueDate = &DueDate
		}
	}
	return &resp
}

type Tariff struct {
	Name        string  `json:"name" validate:"required,alphanum,min=2,max=15,lowercase"`
	Description *string `json:"description"`
	PlanName    string  `json:"planName" validate:"required"`
	DomainName  string  `json:"domainName" validate:"required"`
	DiskQuota   int     `json:"diskQuota" validate:"required"`
	Office      bool    `json:"office" validate:"omitempty"`
	Price       int     `json:"price" validate:"required,omitempty"`
	Type        string  `json:"type"  validate:"required,oneof=personal group options"`
	Regularity  string  `json:"regularity"  validate:"required,oneof=daily monthly"`
}

func (incoming *Tariff) ToDB(r *http.Request) *dbtariffs.DBStruct {
	p := params.GetQueryParams(r)
	return &dbtariffs.DBStruct{
		Name:        incoming.Name,
		OldName:     p.TariffName,
		Type:        &incoming.Type,
		Description: incoming.Description,
		DiskQuota:   &incoming.DiskQuota,
		Office:      &incoming.Office,
		Price:       &incoming.Price,
		Regularity:  &incoming.Regularity,
		Domain: &dbdomains.DBStruct{
			Name: *p.DomainName,
		},
		Plan: &dbplans.DBStruct{
			Name: *p.PlanName,
		},
	}
}

type Tenant struct {
	Name         string `json:"name" validate:"required,alphanum,min=2,max=15,lowercase"`
	Organisation string `json:"organisation" validate:"required"`
	OrderForm    string `json:"orderForm" validate:"required"`
	OrderLink    string `json:"orderLink" validate:"required,url"`
	Type         string `json:"type" validate:"required,oneof=primary regular premium"`
	Plan         string `json:"planName" validate:"required,alphanum,min=2,max=15,lowercase"`
	Domain       string
	Description  *string `json:"description"`
}

func (incoming *Tenant) ToDB(r *http.Request) *dbtenants.DBStruct {
	p := params.GetQueryParams(r)
	resp := dbtenants.DBStruct{
		Description:  incoming.Description,
		Name:         incoming.Name,
		Organisation: &incoming.Organisation,
		OrderForm:    &incoming.OrderForm,
		OrderLink:    &incoming.OrderLink,
		Type:         &incoming.Type,
		OldName:      p.TenantName,
		Plan: &dbplans.DBStruct{
			Name: incoming.Plan,
		},
		Domain: &dbdomains.DBStruct{
			Name: *p.DomainName,
		},
	}
	return &resp
}

type Users struct {
	Email       string `json:"email" validate:"required,email,lowercase"`
	DisplayName string `json:"name" validate:"required"`
	Type        string `json:"type" validate:"required,oneof=full_Admin domain_admin tenant_admin regular"`
	Tariff      string `json:"tariff" validate:"required"`
}

func (incoming *Users) ToDB(r *http.Request) *dbusers.DBStruct {
	p := params.GetQueryParams(r)
	return &dbusers.DBStruct{
		Email:       incoming.Email,
		DisplayName: &incoming.DisplayName,
		Type:        &incoming.Type,
		OldEmail:    p.Email,
		Tariff: &dbtariffs.DBStruct{
			Name: incoming.Tariff,
		},
		Tenant: &dbtenants.DBStruct{
			Name: *p.TenantName,
		},
		Domain: &dbdomains.DBStruct{
			Name: *p.DomainName,
		},
	}
}

type Groups struct {
	Name string `json:"name" validate:"required,alphanum,min=2,max=15,lowercase"`
	Type string `json:"type" validate:"required,oneof=regular office access"`
}

func (incoming *Groups) ToDB(r *http.Request) *dbgroups.DBStruct {
	p := params.GetQueryParams(r)
	resp := dbgroups.DBStruct{
		Name:    incoming.Name,
		Type:    &incoming.Type,
		OldName: p.GroupName,
		Tenant: &dbtenants.DBStruct{
			Name: *p.TenantName,
		},
		Domain: &dbdomains.DBStruct{
			Name: *p.DomainName,
		},
	}
	return &resp
}
