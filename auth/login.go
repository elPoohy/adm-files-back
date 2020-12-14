package auth

import (
	"encoding/json"
	"files-back/auth/directory"
	"files-back/handlers"
	"github.com/go-ldap/ldap/v3"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	var IncomeAuth incomingJSON
	err := json.NewDecoder(r.Body).Decode(&IncomeAuth)
	if err != nil {
		handlers.StatusBadData(err, w)
		return
	}

	err = directory.CheckAuth(IncomeAuth.Username, IncomeAuth.Password)
	if err != nil {
		switch {
		case ldap.IsErrorWithCode(err, ldap.LDAPResultInvalidCredentials):
			handlers.StatusInvalidCredentials(err, w)
		default:
			handlers.StatusError(err, w)
		}
		return
	}

	token, err := GenerateToken(IncomeAuth.Username)
	if err != nil {
		handlers.ReturnError(w, err)
	}
	handlers.ResponseJSON(w, Token{
		Token:   token,
		Expires: "",
	})
}
