package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
)

type Domain struct {
	Name         string  `db:"name"`
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

type NewDomainJSON struct {
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

func queryDomains(params QueryParams) ([]*DomainJSON, error) {
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
	rows, err = Db.NamedQuery(sqlQuery, params)
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
		domainResponse := convertDomainsDBToJSON(&resultDomain)
		resultDomains = append(resultDomains, domainResponse)
	}
	return resultDomains, nil
}

func queryDomain(domainName string) (*DomainJSON, error) {
	var resultDomain *DomainJSON
	var DBDomain Domain
	err := Db.QueryRowx("SELECT name, organisation, primary_url, admin_url, data_path, user_name, type FROM domains WHERE name=$1", domainName).StructScan(&DBDomain)
	if err != nil {
		return resultDomain, err
	}
	resultDomain = convertDomainsDBToJSON(&DBDomain)
	return resultDomain, nil
}

func convertDomainsDBToJSON(dbDomain *Domain) *DomainJSON {
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

func getDomains(w http.ResponseWriter, r *http.Request) {
	domains, err := queryDomains(getQueryParams(r.URL))
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			statusDBNotFound(err, w)
			return
		default:
			statusDBError(err, w)
			return
		}
	}
	responseJSON(w, err, domains)
}

func responseJSON(w http.ResponseWriter, err error, domains interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(domains)
	if err != nil {
		log.Println(err)
	}
}

func getDomain(w http.ResponseWriter, r *http.Request) {
	responseDomain, err := queryDomain(mux.Vars(r)["domainName"])
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			statusDBNotFound(err, w)
			return
		default:
			statusDBError(err, w)
			return
		}
	}
	responseJSON(w, err, responseDomain)
}

func createDomain(w http.ResponseWriter, r *http.Request) {
	newDomain, err := extractDomain(r)
	if err != nil {
		statusBadData(err, w)
	}
	err = insertNewDomain(convertNewDomainsJSONToDB(newDomain))
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			statusDBNotFound(err, w)
			return
		default:
			statusDBError(err, w)
			return
		}
	}
	responseDomain, err := queryDomain(newDomain.Name)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			statusDBNotFound(err, w)
			return
		default:
			statusDBError(err, w)
			return
		}
	}
	responseJSON(w, err, responseDomain)
}

func convertNewDomainsJSONToDB(NewDomain *NewDomainJSON) *Domain {
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

func insertNewDomain(domain *Domain) error {
	result, err := Db.NamedExec(`
			INSERT INTO domains
				(name, organisation, admin_url, primary_url, data_path, password, user_name, type, description)
			VALUES
			    (:name, :organisation, :admin_url, :primary_url, :data_path, :password, :user_name, :type, :description)`,
		domain)
	fmt.Println(result)
	if err != nil {
		return err
	}
	return nil
}

func extractDomain(r *http.Request) (*NewDomainJSON, error) {
	var NewDomain NewDomainJSON
	err := json.NewDecoder(r.Body).Decode(&NewDomain)
	if err != nil {
		return nil, err
	}
	err = validate.Struct(&NewDomain)
	if err != nil {
		return nil, err
	}
	return &NewDomain, nil
}

func updateDomain(w http.ResponseWriter, r *http.Request) {

}

func deleteDomain(w http.ResponseWriter, r *http.Request) {

}
