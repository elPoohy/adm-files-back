package groups

import (
	"files-back/dbase/dbgroups"
	"files-back/handlers"
	"files-back/handlers/incoming"
	"files-back/handlers/params"
	"net/http"
)

func Get(w http.ResponseWriter, r *http.Request) {
	tenants, err := dbgroups.Query(params.GetQueryParams(r))
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
	var n incoming.Groups
	if err := incoming.Extract(r, &n); err != nil {
		handlers.StatusBadData(err, w)
		return
	}
	if err := dbgroups.Insert(n.ToDB(r)); err != nil {
		handlers.ReturnError(w, err)
		return
	}
	handlers.StatusInserted(w)
}

func Update(w http.ResponseWriter, r *http.Request) {
	var n incoming.Groups
	if err := incoming.Extract(r, &n); err != nil {
		handlers.StatusBadData(err, w)
	}
	if err := dbgroups.Update(n.ToDB(r)); err != nil {
		handlers.ReturnError(w, err)
		return
	}
	handlers.StatusInserted(w)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	if err := dbgroups.Delete(params.GetQueryParams(r)); err != nil {
		handlers.StatusBadData(err, w)
		return
	}
	handlers.StatusDeleted(w)
}

func Add(w http.ResponseWriter, r *http.Request) {
	if err := dbgroups.Delete(params.GetQueryParams(r)); err != nil {
		handlers.StatusBadData(err, w)
		return
	}
	handlers.StatusDeleted(w)
}
