package dbdomains

import (
	"files-back/dbase"
	"files-back/handlers/params"
	"github.com/jmoiron/sqlx"
	"log"
)

type DBStruct struct {
	Name         string  `db:"name"`
	OldName      *string `db:"old_name"`
	Organisation *string `db:"organisation"`
	PrimaryURL   *string `db:"primary_url"`
	AdminURL     *string `db:"admin_url"`
	DataPath     *string `db:"data_path"`
	Password     *string `db:"password"`
	UserName     *string `db:"user_name"`
	Version      *string `db:"version"`
	Type         *string `db:"type"`
	Description  *string `db:"description"`
}

func (dbDomain *DBStruct) toJSON() *JSONStruct {
	unknown := "unknown"
	domainResponse := JSONStruct{
		Name:         &dbDomain.Name,
		Organisation: dbDomain.Organisation,
		PrimaryURL:   dbDomain.PrimaryURL,
		AdminURL:     dbDomain.AdminURL,
		DataPath:     dbDomain.DataPath,
		UserName:     dbDomain.UserName,
		Type:         dbDomain.Type,
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
	Name         *string `json:"name,omitempty"`
	Organisation *string `json:"organisation,omitempty"`
	PrimaryURL   *string `json:"primaryUrl,omitempty"`
	AdminURL     *string `json:"adminUrl,omitempty"`
	Version      *string `json:"version,omitempty"`
	Type         *string `json:"type,omitempty"`
	DataPath     *string `json:"data_path,omitempty"`
	UserName     *string `json:"user_name,omitempty"`
	Description  *string `json:"description,omitempty"`
}

func Query(p params.QueryParams) ([]*JSONStruct, error) {
	var res []*JSONStruct
	var rows *sqlx.Rows
	var err error
	var sqlWhere string
	if p.Search != nil {
		sqlWhere = dbase.AppendWhere(sqlWhere) + "(name LIKE :search OR organisation LIKE :search) "
	}
	if p.DomainName != nil {
		sqlWhere = dbase.AppendWhere(sqlWhere) + "(name = :domain_name) "
	}
	if !p.ShowDeleted && !p.ShowDisabled {
		sqlWhere = dbase.AppendWhere(sqlWhere) + "(type NOT IN ('disabled', 'deleted')) "
	}
	if p.ShowDisabled && !p.ShowDeleted {
		sqlWhere = dbase.AppendWhere(sqlWhere) + "(type NOT IN ('deleted')) "
	}
	sqlQuery := `
		SELECT
		       name as name,  primary_url, admin_url, organisation, version, type, data_path, user_name
		FROM domains ` + sqlWhere + `
		LIMIT :limit
		    OFFSET :offset`
	rows, err = dbase.DB.NamedQuery(sqlQuery, p)
	if err != nil {
		return res, err
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
			return res, err
		}
		domainResponse := resultDomain.toJSON()
		res = append(res, domainResponse)
	}
	return res, nil
}

func Delete(p params.QueryParams) error {
	err := dbase.ExecWithChekOne(p, "UPDATE domains SET type=:delete WHERE name=:name RETURNING id")
	if err != nil {
		return err
	}
	return nil
}

func Insert(domain *DBStruct) error {
	err := dbase.ExecWithChekOne(domain,
		`
			INSERT INTO domains
				(name, organisation, admin_url, primary_url, data_path, password, user_name, type, description)
			VALUES
			    (:name, :organisation, :admin_url, :primary_url, :data_path, :password, :user_name, CAST (:type AS domain_type), :description)
			RETURNING id`)
	if err != nil {
		return err
	}
	return nil
}

func Update(domain *DBStruct) error {
	err := dbase.ExecWithChekOne(domain, `
			UPDATE domains
			SET 
			    name = :name,
			    organisation = :organisation,
			    admin_url = :admin_url,
			    primary_url = :primary_url,
			    data_path = :data_path,
			    password = :password,
			    user_name = :user_name,
			    type = CAST (:type AS domain_type),
			    description = :description
			WHERE
				name = :old_name
			RETURNING id`,
	)
	if err != nil {
		return err
	}
	return nil
}
