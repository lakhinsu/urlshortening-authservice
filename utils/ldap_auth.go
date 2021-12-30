package utils

import (
	"crypto/tls"
	"errors"
	"fmt"

	ldap "gopkg.in/ldap.v2"
)

var (
	ldap_string     string
	ldap_tls        string
	userdn          string
	user_fqdn_aname string
)

func init() {
	ldap_string = ReadEnvVar("LDAP_SERVER")

	ldap_tls = ReadEnvVar("LDAP_TLS")

	userdn = ReadEnvVar("LDAP_USER_DN")

	user_fqdn_aname = ReadEnvVar("LDAP_USER_FQDN_ATTRIBUTE_NAME")
}

func GetLdapUser(username string, password string) (map[string]string, error) {
	// Dial up with LDAP server and bind read only user

	// ldap_string := ReadEnvVar("LDAP_SERVER")

	// dial up with the server
	ldap_user_dial, err := ldap.Dial("tcp", ldap_string)
	if err != nil {
		fmt.Println(err.Error())
		defer ldap_user_dial.Close()
		return nil, errors.New("failed to authenticate user, internal server error")
	}

	// ldap_tls := ReadEnvVar("LDAP_TLS")

	// If LDAP_TLS is set to True then reconnect with TLS client.
	if ldap_tls == "True" {
		// Reconnect with TLS
		err = ldap_user_dial.StartTLS(&tls.Config{})
		if err != nil {
			defer ldap_user_dial.Close()
			fmt.Println(err.Error())
			return nil, errors.New("failed to authenticate user, internal server error")
		}
	}

	// Load the USER DN fron ENV
	// userdn := ReadEnvVar("LDAP_USER_DN")
	user_bind := fmt.Sprintf("uid=%s,%s", username, userdn)
	// Bind as the user to verify their password
	err = ldap_user_dial.Bind(user_bind, password)
	if err != nil {
		fmt.Println(err.Error())
		defer ldap_user_dial.Close()
		return nil, errors.New("failed to authenticate user, incorrect user name or password")
	}

	// user_fqdn_aname := ReadEnvVar("LDAP_USER_FQDN_ATTRIBUTE_NAME")

	searchRequest := ldap.NewSearchRequest(
		userdn, // The base dn to search
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(objectClass=organizationalPerson)(uid=%s))", username), // The filter to apply
		[]string{user_fqdn_aname}, // A list attributes to retrieve
		nil,
	)

	sr, err := ldap_user_dial.Search(searchRequest)
	if err != nil {
		fmt.Println(err.Error())
		defer ldap_user_dial.Close()
		return nil, errors.New("failed to authenticate user, internal server error")
	}
	user_fqdn_value := sr.Entries[0].GetAttributeValue("mail")

	results := map[string]string{
		"email": user_fqdn_value,
	}
	defer ldap_user_dial.Close()
	return results, nil
}
