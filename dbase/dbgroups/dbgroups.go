package dbgroups

import (
	"files-back/dbase"
	"files-back/dbase/dbdomains"
	"files-back/dbase/dbtenants"
	"files-back/handlers/params"
	"github.com/jmoiron/sqlx"
	"log"
)

type DBStruct struct {
	Name    string              `db:"name"`
	OldName *string             `db:"old_name"`
	Type    *string             `db:"type"`
	Tenant  *dbtenants.DBStruct `db:"tenant"`
	Domain  *dbdomains.DBStruct `db:"domain"`
}

func (dbGroups *DBStruct) toJSON() *JSONStruct {
	return &JSONStruct{
		Name:   &dbGroups.Name,
		Tenant: &dbGroups.Tenant.Name,
		Type:   dbGroups.Type,
	}
}

type JSONStruct struct {
	Name   *string `json:"name"`
	Type   *string `json:"type,omitempty"`
	Tenant *string `json:"tenant,omitempty"`
}

func Query(p params.QueryParams) ([]*JSONStruct, error) {
	var res []*JSONStruct
	var rows *sqlx.Rows
	var err error
	var sqlWhere string
	if p.Search != nil {
		sqlWhere = dbase.AppendWhere(sqlWhere) + "(g.name LIKE :search) "
	}
	if p.TenantName != nil {
		sqlWhere = dbase.AppendWhere(sqlWhere) + "(t.name = :tenant_name AND d.name = :domain_name) "
	}
	if p.GroupName != nil {
		sqlWhere = dbase.AppendWhere(sqlWhere) + "(g.name = :group_name) "
	}
	sqlQuery := `
		SELECT
				g.name as name, g.type as type, t.name as "tenant.name", d.name as "domain.name"
		FROM groups g
		JOIN tenants t ON t.id = g.tenant_id
		JOIN domains d ON d.id = t.domain_id ` + sqlWhere + `
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
		UPDATE groups SET type = :delete WHERE id IN 
			(SELECT
		       g.id
		    FROM groups g
		    JOIN tenants t ON t.id = g.tenant_id
			JOIN domains d ON d.id = t.domain_id
			WHERE
				t.name=:tenant_name AND g.name=:group_name AND d.name=:domain_name)
			RETURNING id`,
	)
	if err != nil {
		return err
	}
	return nil
}

func Insert(group *DBStruct) error {
	err := dbase.ExecWithChekOne(group, `INSERT INTO groups
							(name, type, tenant_id)
						SELECT
							:name, CAST (:type AS group_type), t.id
						FROM tenants t
						JOIN domains d on d.id = t.domain_id
						WHERE d.name = :domain.name AND t.name = :tenant.name
						RETURNING id`)
	if err != nil {
		return err
	}
	return nil
}

func Update(group *DBStruct) error {
	err := dbase.ExecWithChekOne(group,
		`UPDATE groups
			SET 
			    name = :name,
				type = CAST (:type AS group_type)
			WHERE
				id IN (
				SELECT
		       		g.id
				FROM groups g
		    		JOIN tenants t ON t.id = g.tenant_id
					JOIN domains d ON d.id = t.domain_id
					WHERE
						t.name=:tenant.name AND g.name=:old_name AND d.name=:domain.name)
			RETURNING id`)
	if err != nil {
		return err
	}
	return nil
}
