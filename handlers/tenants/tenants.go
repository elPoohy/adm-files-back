package tenants

import (
	"files-back/dbase/dbtenants"
	"files-back/handlers"
	"files-back/handlers/incoming"
	"files-back/handlers/params"
	"net/http"
)

func Get(w http.ResponseWriter, r *http.Request) {
	tenants, err := dbtenants.Query(params.GetQueryParams(r))
	if err != nil {
		handlers.StatusBadData(err, w)
		return
	}
	if len(tenants) == 1 {
		handlers.ResponseJSON(w, tenants[0])
	} else {
		handlers.ResponseJSON(w, tenants)
	}
}

func Create(w http.ResponseWriter, r *http.Request) {
	var n incoming.Tenant
	if err := incoming.Extract(r, &n); err != nil {
		handlers.StatusBadData(err, w)
		return
	}
	if err := dbtenants.Insert(n.ToDB(r)); err != nil {
		handlers.ReturnError(w, err)
		return
	}
	handlers.StatusInserted(w)
}

func Update(w http.ResponseWriter, r *http.Request) {
	var n incoming.Tenant
	if err := incoming.Extract(r, &n); err != nil {
		handlers.StatusBadData(err, w)
		return
	}
	if err := dbtenants.Update(n.ToDB(r)); err != nil {
		handlers.StatusBadData(err, w)
		return
	}
	handlers.StatusInserted(w)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	if err := dbtenants.Delete(params.GetQueryParams(r)); err != nil {
		handlers.StatusBadData(err, w)
		return
	}
	handlers.StatusDeleted(w)
}
