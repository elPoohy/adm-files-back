package domains

import (
	"files-back/dbase/dbdomains"
	"files-back/handlers"
	"files-back/handlers/incoming"
	"files-back/handlers/params"
	"net/http"
)

func Get(w http.ResponseWriter, r *http.Request) {
	domains, err := dbdomains.Query(params.GetQueryParams(r))
	switch {
	case err != nil:
		handlers.StatusBadData(err, w)
		return
	case len(domains) == 1:
		handlers.ResponseJSON(w, domains[0])
	default:
		handlers.ResponseJSON(w, domains)
	}
}

func Create(w http.ResponseWriter, r *http.Request) {
	var n incoming.Domain
	if err := incoming.Extract(r, &n); err != nil {
		handlers.StatusBadData(err, w)
		return
	}
	if err := dbdomains.Insert(n.ToDB(r)); err != nil {
		handlers.ReturnError(w, err)
		return
	}
	handlers.StatusInserted(w)
}

func Update(w http.ResponseWriter, r *http.Request) {
	var n incoming.Domain
	if err := incoming.Extract(r, &n); err != nil {
		handlers.StatusBadData(err, w)
		return
	}
	if err := dbdomains.Update(n.ToDB(r)); err != nil {
		handlers.StatusBadData(err, w)
		return
	}
	handlers.StatusInserted(w)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	if err := dbdomains.Delete(params.GetQueryParams(r)); err != nil {
		handlers.StatusBadData(err, w)
		return
	}
	handlers.StatusDeleted(w)
}
