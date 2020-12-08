package directory

import (
	"crypto/tls"
	"fmt"
	"github.com/go-ldap/ldap/v3"
	"log"
	"os"
)

var (
	BindUsername string
	BindPassword string
	BaseDN       string
	LDAP         *ldap.Conn
)

func InitLDAP(server, port string) {
	dial, err := ldap.Dial("tcp", fmt.Sprintf("%s:%s", server, port))
	if err != nil {
		log.Printf("Unable to connect to ldap: %v", err)
		os.Exit(1)
	}
	LDAP = dial
}

func CheckAuth(email, password string) error {

	err := LDAP.StartTLS(&tls.Config{InsecureSkipVerify: true})
	if err != nil {
		return err
	}

	err = LDAP.Bind(BindUsername, BindPassword)
	if err != nil {
		return nil
	}

	result, err := SearchLDAPbyMail(email)
	if err != nil {
		return err
	}

	if len(result.Entries) != 1 {
		return err
	}

	userDN := result.Entries[0].DN

	err = LDAP.Bind(userDN, password)
	if err != nil {
		return err
	}

	err = LDAP.Bind(BindUsername, BindPassword)
	if err != nil {
		return err
	}
	return nil
}

func SearchLDAPbyMail(email string) (*ldap.SearchResult, error) {
	return LDAP.Search(ldap.NewSearchRequest(
		BaseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		fmt.Sprintf("(&(objectClass=inetOrgPerson)(mail=%s))", email),
		[]string{"dn"},
		nil,
	))
}
