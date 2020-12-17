package dbusers

import (
	"files-back/dbase"
	"files-back/dbase/dbdomains"
	"files-back/dbase/dbtariffs"
	"files-back/dbase/dbtenants"
	"files-back/handlers/params"
	"github.com/jmoiron/sqlx"
	"log"
)

type DBStruct struct {
	Email       string              `db:"email"`
	OldEmail    *string             `db:"old_email"`
	DisplayName *string             `db:"display_name"`
	Type        *string             `db:"type"`
	Free        *int                `db:"free"`
	Tariff      *dbtariffs.DBStruct `db:"tariff"`
	Tenant      *dbtenants.DBStruct `db:"tenant"`
	Domain      *dbdomains.DBStruct `db:"domain"`
}

func (dbUsers *DBStruct) toJSON() *JSONStruct {
	return &JSONStruct{
		Email:       &dbUsers.Email,
		DisplayName: dbUsers.DisplayName,
		Free:        dbUsers.Free,
		Tenant:      &dbUsers.Tenant.Name,
		Domain:      &dbUsers.Domain.Name,
		Type:        dbUsers.Type,
	}
}

type JSONStruct struct {
	Email       *string `json:"email"`
	DisplayName *string `json:"name"`
	Type        *string `json:"type"`
	Free        *int    `json:"free,omitempty"`
	Tenant      *string `json:"tenant"`
	Domain      *string `json:"domain"`
}

func Query(p params.QueryParams) ([]*JSONStruct, error) {
	var res []*JSONStruct
	var rows *sqlx.Rows
	var err error
	var sqlWhere string
	if p.Search != nil {
		sqlWhere = dbase.AppendWhere(sqlWhere) + "(u.email LIKE :search OR u.display_name LIKE :search) "
	}
	if p.TenantName != nil {
		sqlWhere = dbase.AppendWhere(sqlWhere) + "(t.name = :tenant_name AND d.name = :domain_name) "
	}
	if p.GroupName != nil {
		sqlWhere = dbase.AppendWhere(sqlWhere) + "(u.name = :group_name) "
	}
	sqlQuery := `
		SELECT
				u.email as email, u.display_name as display_name, u.type as type, u.free as free, tf.name as "tariff.name", t.name as "tenant.name", d.name as "domain.name"
		FROM users u
		JOIN tariff tf ON tf.id = u.tariff_id
		JOIN tenants t ON t.id = u.tenant_id
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
		UPDATE users SET type = :delete WHERE id IN 
			(SELECT
		       u.id
		    FROM users u
		    JOIN tenants t ON t.id = u.tenant_id
			JOIN domains d ON d.id = t.domain_id
			WHERE
				t.name=:tenant_name AND u.name=:email AND d.name=:domain_name)
			RETURNING id`,
	)
	if err != nil {
		return err
	}
	return nil
}

func Insert(group *DBStruct) error {
	err := dbase.ExecWithChekOne(group, `INSERT INTO groups
							(email, display_name, type, tariff_id, tenant_id)
						SELECT
							:email, :display_name, , CAST (:type AS group_type), tf.id, t.id
						FROM tariffs tf
						JOIN domains d ON d.id = tf.domain_id
						JOIN tenants t ON t.domain_id = d.id 
						WHERE d.name = :domain.name AND t.name = :tenant.name AND :tariff.name
						RETURNING id`)
	if err != nil {
		return err
	}
	return nil
}

func Update(group *DBStruct) error {
	err := dbase.ExecWithChekOne(group,
		`UPDATE users SET
					SET 
			    		(email, display_name, type, tariff_id, tenant_id)
						(SELECT
							:email, :display_name, , CAST (:type AS group_type), tf.id, t.id
						FROM tariffs tf
						JOIN domains d ON d.id = tf.domain_id
						JOIN tenants t ON t.domain_id = d.id 
						WHERE d.name = :domain.name AND t.name = :tenant.name AND :tariff.name)
					WHERE id IN 
						(SELECT
		       				u.id
						FROM users u
						JOIN tenants t ON t.id = u.tenant_id
						JOIN domains d ON d.id = t.domain_id
						WHERE
							t.name=:tenant_name AND u.name=:old_email AND d.name=:domain_name)
						RETURNING id`)
	if err != nil {
		return err
	}
	return nil
}
