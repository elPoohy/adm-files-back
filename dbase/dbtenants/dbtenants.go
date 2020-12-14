package dbtenants

import (
	"files-back/dbase"
	"files-back/dbase/dbdomains"
	"files-back/dbase/dbplans"
	"files-back/handlers/params"
	"github.com/jmoiron/sqlx"
	"log"
)

type DBStruct struct {
	Name         string              `db:"tenant_name"`
	OldName      *string             `db:"old_name"`
	Organisation *string             `db:"organisation"`
	OrderForm    *string             `db:"order_form"`
	OrderLink    *string             `db:"order_link"`
	Type         *string             `db:"type"`
	Description  *string             `db:"description"`
	Domain       *dbdomains.DBStruct `db:"domain"`
	Plan         *dbplans.DBStruct   `db:"plan"`
}

func (dbTenant *DBStruct) toJSON() *JSONStruct {
	unknown := "unknown"
	domainResponse := JSONStruct{
		Name:         &dbTenant.Name,
		Organisation: dbTenant.Organisation,
		OrderForm:    dbTenant.OrderForm,
		OrderLink:    dbTenant.OrderLink,
		Domain:       dbTenant.Domain,
		Plan:         dbTenant.Plan,
		Type:         dbTenant.Type,
	}

	if dbTenant.Description != nil {
		domainResponse.Description = dbTenant.Description
	} else {
		domainResponse.Description = &unknown
	}
	return &domainResponse
}

type JSONStruct struct {
	Name         *string             `json:"name"`
	Organisation *string             `json:"organisation,omitempty"`
	OrderForm    *string             `json:"orderForm,omitempty"`
	OrderLink    *string             `json:"orderLink,omitempty"`
	Type         *string             `json:"type,omitempty"`
	Description  *string             `json:"description,omitempty"`
	Domain       *dbdomains.DBStruct `json:"domain,omitempty"`
	Plan         *dbplans.DBStruct   `json:"plan,omitempty"`
}

func Query(p params.QueryParams) ([]*JSONStruct, error) {
	var res []*JSONStruct
	var rows *sqlx.Rows
	var err error
	var sqlWhere string
	if p.DomainName != nil && p.PlanName != nil {
		sqlWhere = dbase.AppendWhere(sqlWhere) + "(p.name = :plan_name AND d.name = :domain_name) "
	}
	if p.Search != nil {
		sqlWhere = dbase.AppendWhere(sqlWhere) + "(t.name LIKE :search OR t.organisation LIKE :search) "
	}
	if p.TenantName != nil {
		sqlWhere = dbase.AppendWhere(sqlWhere) + "(t.name = :tenant_name) "
	}
	if p.DomainName != nil && p.TenantName != nil {
		sqlWhere = dbase.AppendWhere(sqlWhere) + "(d.name = :domain_name AND t.name = :tenant_name) "
	}
	if !p.ShowDeleted && !p.ShowDisabled {
		sqlWhere = dbase.AppendWhere(sqlWhere) + "(t.type NOT IN ('disabled', 'deleted')) "
	}
	if p.ShowDisabled && !p.ShowDeleted {
		sqlWhere = dbase.AppendWhere(sqlWhere) + "(t.type NOT IN ('deleted')) "
	}
	sqlQuery := `
		SELECT
				t.name as tenant_name, t.organisation as organisation, t.order_form as order_form, t.order_link as order_link, t.description as description, t.type as type, d.name as "domain.name", p.name as "plan.name"
		FROM tenants t
		JOIN domains d ON d.id = t.domain_id 
		JOIN plans p ON d.id = p.domain_id ` + sqlWhere + `
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
		UPDATE tenants SET type = :delete WHERE id IN 
			(SELECT
		       t.id
		    FROM tenants t
		    JOIN domains d ON d.id = t.domain_id
			WHERE
				t.name=:tenant_name AND d.name=:domain_name)
			RETURNING id`,
	)
	if err != nil {
		return err
	}
	return nil
}

func Insert(tenant *DBStruct) error {
	sqlQuery := `INSERT INTO tenants
							(name, organisation, order_form, order_link, description, type, domain_id, plan_id)
						SELECT
							:tenant_name, :organisation, :order_form, :order_link, :description, :type, d.id, p.id
						FROM domains d
						JOIN plans p on d.id = p.domain_id
						WHERE d.name = :domain.name AND p.name = :plan.name
						RETURNING id`
	err := dbase.ExecWithChekOne(tenant, sqlQuery)
	if err != nil {
		return err
	}
	return nil
}

func Update(tenant *DBStruct, p params.QueryParams) error {
	tenant.OldName = p.TenantName
	err := dbase.ExecWithChekOne(tenant,
		`UPDATE tenants
			SET 
			    (name, organisation, order_form, order_link, description, type, domain_id, plan_id) = 
			    (SELECT
			    	:tenant_name, :organisation, :order_form, :order_link, :description, CAST (:type AS tenant_type), d.id, p.id
				FROM domains d
				JOIN plans p on d.id = p.domain_id
				WHERE d.name = :domain.name AND p.name = :plan.name)
			WHERE
				name = :old_name
			RETURNING id`)
	if err != nil {
		return err
	}
	return nil
}
