package dbdomain

import (
	"files-back/dbase"
	"files-back/handlers"
	"github.com/jmoiron/sqlx"
	"log"
)

type DBStruct struct {
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

func (dbDomain *DBStruct) toJSON() *JSONStruct {
	unknown := "unknown"
	domainResponse := JSONStruct{
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

type JSONStruct struct {
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

func QueryAll(params handlers.QueryParams) ([]*JSONStruct, error) {
	var resultDomains []*JSONStruct
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
		var resultDomain DBStruct
		err := rows.StructScan(&resultDomain)
		if err != nil {
			return resultDomains, err
		}
		domainResponse := resultDomain.toJSON()
		resultDomains = append(resultDomains, domainResponse)
	}
	return resultDomains, nil
}

func QueryOne(domainName string) (*JSONStruct, error) {
	var resultDomain *JSONStruct
	var DBDomain DBStruct
	err := dbase.DB.QueryRowx("SELECT name, organisation, primary_url, admin_url, data_path, user_name, type FROM domains WHERE name=$1", domainName).StructScan(&DBDomain)
	if err != nil {
		return resultDomain, err
	}
	resultDomain = DBDomain.toJSON()
	return resultDomain, nil
}

func Delete(domainName string) error {
	_, err := dbase.DB.Exec("DELETE FROM domains WHERE name=$1", domainName)
	if err != nil {
		return err
	}
	return nil
}

func Insert(domain *DBStruct) error {
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

func Update(domain *DBStruct, domainName string) error {

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
