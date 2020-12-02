package dbplans

import (
	"files-back/dbase"
	"files-back/handlers"
	"github.com/jmoiron/sqlx"
	"log"
	"time"
)

type DBStruct struct {
	Name        string     `db:"name"`
	OldName     *string    `db:"old_name"`
	DomainName  string     `db:"domain_name"`
	FromDate    time.Time  `db:"from_date"`
	DueDate     *time.Time `db:"due_date"`
	Type        string     `db:"type"`
	Description *string    `db:"description"`
}

func (dbDomain *DBStruct) toJSON() *JSONStruct {
	unknown := "unknown"
	domainResponse := JSONStruct{
		Name:       &dbDomain.Name,
		DomainName: &dbDomain.DomainName,
		FromDate:   &dbDomain.FromDate,
		DueDate:    dbDomain.DueDate,
		Type:       &dbDomain.Type,
	}

	if dbDomain.Description != nil {
		domainResponse.Description = dbDomain.Description
	} else {
		domainResponse.Description = &unknown
	}
	return &domainResponse
}

type JSONStruct struct {
	Name        *string    `json:"name"`
	DomainName  *string    `json:"domainName"`
	FromDate    *time.Time `json:"fromDate"`
	DueDate     *time.Time `json:"dueDate"`
	Type        *string    `json:"type"`
	Description *string    `json:"description"`
}

func Query(p handlers.QueryParams) ([]*JSONStruct, error) {
	var res []*JSONStruct
	var rows *sqlx.Rows
	var err error
	var sqlWhere string
	if p.Search != nil {
		sqlWhere = dbase.AppendWhere(sqlWhere) + "(p.name LIKE :search OR d.name LIKE :search OR d.organisation LIKE :search) "
	}
	if p.DomainName != nil {
		sqlWhere = dbase.AppendWhere(sqlWhere) + "(d.name = :domain_name) "
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
		       p.name as name, d.name as domain_name, p.from_date as from_date, p.due_date as due_date, p.type as type, p.description as description
		FROM plans p
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

func Delete(p handlers.QueryParams) error {
	_, err := dbase.DB.NamedExec(`
		UPDATE plans SET type = :delete WHERE id IN 
			(SELECT
		       p.id
		    FROM plans p
		    JOIN domains d ON d.id = p.domain_id
			WHERE
				p.name=:domain_name AND d.name=:plan_name)`,
		p,
	)
	if err != nil {
		return err
	}
	return nil
}

func Insert(plan *DBStruct) error {
	_, err := dbase.DB.NamedExec(`
			INSERT INTO plans
				(name, from_date, description, due_date, type, domain_id)
			SELECT
			    :name, :from_date, :description, :due_date, :type, domains.id
			FROM domains
			WHERE domains.name = :domain_name`,
		plan)
	if err != nil {
		return err
	}
	return nil
}

func Update(plan *DBStruct, p handlers.QueryParams) error {
	plan.OldName = p.PlanName
	_, err := dbase.DB.NamedExec(`
			UPDATE plans
			SET 
			    name = :name,
			    from_date = :from_date,
			    description = :description,
			    due_date =  :due_date,
			    type = :type
			WHERE id IN 
			   (SELECT
		          p.id
		        FROM plans p
	    	    JOIN domains d ON d.id = p.domain_id
			    WHERE
					p.name = :old_name AND d.name = :domain_name)`,
		plan)
	if err != nil {
		return err
	}
	return nil
}
