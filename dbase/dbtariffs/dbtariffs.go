package dbtariffs

import (
	"files-back/dbase"
	"files-back/dbase/dbdomains"
	"files-back/dbase/dbplans"
	"files-back/handlers/params"
	"github.com/jmoiron/sqlx"
	"log"
)

type DBStruct struct {
	Name        string              `db:"name"`
	OldName     *string             `db:"old_name"`
	Domain      *dbdomains.DBStruct `db:"domain"`
	Plan        *dbplans.DBStruct   `db:"plan"`
	Type        *string             `db:"type"`
	Description *string             `db:"description"`
	DiskQuota   *int                `db:"disk_quota"`
	Office      *bool               `db:"office"`
	Price       *int                `db:"price"`
	Regularity  *string             `db:"regularity"`
}

func (dbTariff *DBStruct) toJSON() *JSONStruct {
	unknown := "unknown"
	domainResponse := JSONStruct{
		Name:       &dbTariff.Name,
		DomainName: &dbTariff.Domain.Name,
		Type:       dbTariff.Type,
		PlanName:   &dbTariff.Plan.Name,
		DiskQuota:  dbTariff.DiskQuota,
		Office:     dbTariff.Office,
		Price:      dbTariff.Price,
		Regularity: dbTariff.Regularity,
	}

	if dbTariff.Description != nil {
		domainResponse.Description = dbTariff.Description
	} else {
		domainResponse.Description = &unknown
	}
	return &domainResponse
}

type JSONStruct struct {
	Name        *string `json:"name,omitempty"`
	DomainName  *string `json:"domainName,omitempty"`
	Type        *string `json:"type,omitempty"`
	Description *string `json:"description,omitempty"`
	PlanName    *string `json:"plan"`
	DiskQuota   *int    `json:"diskQuota"`
	Office      *bool   `json:"office"`
	Price       *int    `json:"price"`
	Regularity  *string `json:"regularity"`
}

func Query(p params.QueryParams) ([]*JSONStruct, error) {
	var res []*JSONStruct
	var rows *sqlx.Rows
	var err error
	var sqlWhere string
	if p.Search != nil {
		sqlWhere = dbase.AppendWhere(sqlWhere) + "(t.name LIKE :search OR p.name LIKE :search OR t.description LIKE :search) "
	}
	if p.TariffName != nil {
		sqlWhere = dbase.AppendWhere(sqlWhere) + "(t.name = :tariff_name"
	}
	if p.DomainName != nil && p.PlanName != nil {
		sqlWhere = dbase.AppendWhere(sqlWhere) + "(d.name = :domain_name AND p.name = :plan_name) "
	}
	if !p.ShowDeleted && !p.ShowDisabled {
		sqlWhere = dbase.AppendWhere(sqlWhere) + "(p.type NOT IN ('disabled', 'deleted')) "
	}
	if p.ShowDisabled && !p.ShowDeleted {
		sqlWhere = dbase.AppendWhere(sqlWhere) + "(p.type NOT IN ('deleted')) "
	}
	sqlQuery := `
		SELECT
		       t.name as name, d.name as domain_name, p.name as plan_name, t.type as type, t.description as description, t.disk_quota as disk_quota, t.office as office, t.price as price, t.regularity as regularity 
		FROM tariffs t
		JOIN plans p ON p.id = t.plan_id
		JOIN domains d ON d.id = p.domain_id ` + sqlWhere + `
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
	err := dbase.ExecWithChekOne(p, `
		UPDATE plans SET type = :delete WHERE id IN 
			(SELECT
		          t.id
		        FROM tariffs t
				JOIN plans p ON p.id = t.domain_id
	    	    JOIN domains d ON d.id = p.domain_id
			    WHERE
					t.name = :old_name AND d.name = :domain.name AND p.name = :plan.name)
			RETURNING id`,
	)
	if err != nil {
		return err
	}
	return nil
}

func Insert(plan *DBStruct) error {
	err := dbase.ExecWithChekOne(plan,
		`INSERT INTO plans
				(name,  description, disk_quota, office, price, type, regularity, domain_id, plan_id)
			SELECT
			    :name, :description, :disk_quota, :office, :price, CAST (:type AS plan_type), regularity, d.id, p.id
			FROM plans p
			JOIN domains d ON p.domain_id = d.id
			WHERE d.name = :domain_name AND p.name = :plan_name
			RETURNING id`)
	if err != nil {
		return err
	}
	return nil
}

func Update(plan *DBStruct) error {
	err := dbase.ExecWithChekOne(plan,
		`UPDATE plans
			SET 
			    name = :name,
			    description = :description,
			    disk_quota = :disk_quota,
			    office = :office,
				price = :price,
				regularity = :regularity
			    type = CAST (:type AS plan_type)
			WHERE id IN 
			   (SELECT
		          t.id
		        FROM tariffs t
				JOIN plans p ON p.id = t.domain_id
	    	    JOIN domains d ON d.id = p.domain_id
			    WHERE
					t.name = :old_name AND d.name = :domain.name AND p.name = :plan.name)
			RETURNING id`,
	)
	if err != nil {
		return err
	}
	return nil
}
