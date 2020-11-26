package domains

import (
	"database/sql"
	"encoding/json"
	"files-back/dbase"
	"files-back/handlers"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
)

type Domain struct {
	Name         string  `db:"name"`
	OldName      *string `db:"old_name"`
	Organisation string  `db:"organisation"`
	PrimaryURL   string  `db:"primary_url"`
	AdminURL     string  `db:"admin_url"`
	DataPath     string  `db:"data_path"`
	Password     string  `db:"password"`
	UserName     string  `db:"user_name"`
	Version      *string `db:"version"`
	Type         string  `db:"type"`
	Description  *string `db:"description"`
}

type DomainJSON struct {
	Name         *string `json:"name"`
	Organisation *string `json:"organisation"`
	PrimaryURL   *string `json:"primaryUrl"`
	AdminURL     *string `json:"adminUrl"`
	Version      *string `json:"version"`
	Type         *string `json:"type"`
	DataPath     *string `json:"data_path"`
	UserName     *string `json:"user_name"`
	Description  *string `json:"description"`
}

type inDomainJSON struct {
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

var validate *validator.Validate

func queryDomains(params handlers.QueryParams) ([]*DomainJSON, error) {
	var resultDomains []*DomainJSON
	var rows *sqlx.Rows
	var err error
	sqlQuery := `
		SELECT
		       name,  primary_url, admin_url, organisation, version, type
		FROM domains
		LIMIT :limit
		    OFFSET :offset`
	if params.Search != "%%" {
		sqlQuery = `
			SELECT
					name, primary_url, admin_url, organisation, version, type FROM domains 
			WHERE
					name LIKE :search OR organisation LIKE :search
			LIMIT :limit
			OFFSET :offset`
	}
	rows, err = dbase.DB.NamedQuery(sqlQuery, params)
	if err != nil {
		return resultDomains, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Println(err)
		}
	}()
	for rows.Next() {
		var resultDomain Domain
		err := rows.StructScan(&resultDomain)
		if err != nil {
			return resultDomains, err
		}
		domainResponse := resultDomain.toJSON()
		resultDomains = append(resultDomains, domainResponse)
	}
	return resultDomains, nil
}

func queryDomain(domainName string) (*DomainJSON, error) {
	var resultDomain *DomainJSON
	var DBDomain Domain
	err := dbase.DB.QueryRowx("SELECT name, organisation, primary_url, admin_url, data_path, user_name, type FROM domains WHERE name=$1", domainName).StructScan(&DBDomain)
	if err != nil {
		return resultDomain, err
	}
	resultDomain = DBDomain.toJSON()
	return resultDomain, nil
}

func deleteDomainDB(domainName string) error {
	_, err := dbase.DB.Exec("DELETE FROM domains WHERE name=$1", domainName)
	if err != nil {
		return err
	}
	return nil
}

func (dbDomain *Domain) toJSON() *DomainJSON {
	unknown := "unknown"
	domainResponse := DomainJSON{
		Name:         &dbDomain.Name,
		Organisation: &dbDomain.Organisation,
		PrimaryURL:   &dbDomain.PrimaryURL,
		AdminURL:     &dbDomain.AdminURL,
		DataPath:     &dbDomain.DataPath,
		UserName:     &dbDomain.UserName,
		Type:         &dbDomain.Type,
	}

	if dbDomain.Version != nil {
		domainResponse.Version = dbDomain.Version
	} else {
		domainResponse.Version = &unknown
	}
	if dbDomain.Description != nil {
		domainResponse.Description = dbDomain.Description
	} else {
		domainResponse.Description = &unknown
	}
	return &domainResponse
}

func GetAll(w http.ResponseWriter, r *http.Request) {
	domains, err := queryDomains(handlers.GetQueryParams(r.URL))
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
	responseDomain, err := queryDomain(extractDomainName(r))
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
	newDomain, err := extractDomain(r)
	if err != nil {
		handlers.StatusBadData(err, w)
	}
	err = insertDomainToDB(newDomain.toDB())
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
	responseDomain, err := queryDomain(newDomain.Name)
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

func (NewDomain *inDomainJSON) toDB() *Domain {
	domainResponse := Domain{
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

func insertDomainToDB(domain *Domain) error {
	_, err := dbase.DB.NamedExec(`
			INSERT INTO domains
				(name, organisation, admin_url, primary_url, data_path, password, user_name, type, description)
			VALUES
			    (:name, :organisation, :admin_url, :primary_url, :data_path, :password, :user_name, :type, :description)`,
		domain)
	if err != nil {
		return err
	}
	return nil
}

func updateDomainDB(domain *Domain, domainName string) error {
	domain.OldName = &domainName
	_, err := dbase.DB.NamedExec(`
			UPDATE domains
			SET 
			    name = :name,
			    organisation = :organisation,
			    admin_url = :admin_url,
			    primary_url = :primary_url,
			    data_path = :data_path,
			    password = :password,
			    user_name = :user_name,
			    type = :type,
			    description = :description
			WHERE
				name = :old_name
			    `,
		domain)
	if err != nil {
		return err
	}
	return nil
}

func extractDomain(r *http.Request) (*inDomainJSON, error) {
	var NewDomain inDomainJSON
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

func Update(w http.ResponseWriter, r *http.Request) {
	newDomain, err := extractDomain(r)
	if err != nil {
		handlers.StatusBadData(err, w)
	}
	err = updateDomainDB(newDomain.toDB(), extractDomainName(r))
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
	responseDomain, err := queryDomain(newDomain.Name)
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

func extractDomainName(r *http.Request) string {
	return mux.Vars(r)["domainName"]
}

func Delete(w http.ResponseWriter, r *http.Request) {
	err := deleteDomainDB(extractDomainName(r))
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
