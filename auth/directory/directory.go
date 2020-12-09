package directory

import (
	"crypto/tls"
	"fmt"
	"github.com/go-ldap/ldap/v3"
)

var (
	LDAPUsername string
	LDAPPassword string
	LDAPServer   string
	BaseDN       string
)

func LDAPDial() (error, *ldap.Conn) {
	dial, err := ldap.Dial("tcp", LDAPServer)
	if err != nil {
		return err, nil
	}
	err = dial.StartTLS(&tls.Config{InsecureSkipVerify: true})
	if err != nil {
		return err, nil
	}
	err = dial.Bind(LDAPUsername, LDAPPassword)
	if err != nil {
		return err, nil
	}
	return nil, dial
}

func CheckAuth(email, password string) error {
	err, dial := LDAPDial()
	if err != nil {
		return err
	}
	defer dial.Close()

	result, err := SearchLDAPbyMail(email, dial)
	if err != nil {
		return err
	}

	if len(result.Entries) != 1 {
		return ldap.NewError(ldap.LDAPResultInvalidCredentials, nil)
	}

	err = dial.Bind(result.Entries[0].DN, password)
	if err != nil {
		return err
	}

	return nil
}

func SearchLDAPbyMail(email string, dial *ldap.Conn) (*ldap.SearchResult, error) {
	return dial.Search(ldap.NewSearchRequest(
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
