package auth

import (
	"encoding/json"
	"files-back/auth/directory"
	"files-back/handlers"
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
		handlers.StatusError(err, w)
	}

	token, err := GenerateToken(IncomeAuth.Username)

	handlers.ResponseJSON(w, Token{
		Token:   token,
		Expires: "",
	})
}
