package users

import (
	"files-back/dbase/dbusers"
	"files-back/handlers"
	"files-back/handlers/incoming"
	"files-back/handlers/params"
	"net/http"
)

func Get(w http.ResponseWriter, r *http.Request) {
	plans, err := dbusers.Query(params.GetQueryParams(r))
	switch {
	case err != nil:
		handlers.StatusBadData(err, w)
		return
	case len(plans) == 1:
		handlers.ResponseJSON(w, plans[0])
	default:
		handlers.ResponseJSON(w, plans)
	}
}

func Create(w http.ResponseWriter, r *http.Request) {
	var n incoming.Users
	if err := incoming.Extract(r, &n); err != nil {
		handlers.StatusBadData(err, w)
		return
	}
	if err := dbusers.Insert(n.ToDB(r)); err != nil {
		handlers.ReturnError(w, err)
		return
	}
	handlers.StatusInserted(w)
}

func Update(w http.ResponseWriter, r *http.Request) {
	var n incoming.Users
	if err := incoming.Extract(r, &n); err != nil {
		handlers.StatusBadData(err, w)
	}

	if err := dbusers.Update(n.ToDB(r)); err != nil {
		handlers.StatusBadData(err, w)
		return
	}
	handlers.StatusInserted(w)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	if err := dbusers.Delete(params.GetQueryParams(r)); err != nil {
		handlers.ReturnError(w, err)
	}
	handlers.StatusDeleted(w)
}

func Add(w http.ResponseWriter, r *http.Request) {
	if err := dbusers.Delete(params.GetQueryParams(r)); err != nil {
		handlers.ReturnError(w, err)
	}
	handlers.StatusDeleted(w)
}
