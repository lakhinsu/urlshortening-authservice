package utils

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
	ldap "gopkg.in/ldap.v2"
)

var (
	ldap_string     string
	ldap_tls        string
	userdn          string
	user_fqdn_aname string
	ssl_verify      string
	serverName      string
)

func init() {
	ldap_string = GetEnvVar("LDAP_CONNECTION_STRING")

	ldap_tls = GetEnvVar("LDAP_TLS")

	userdn = GetEnvVar("LDAP_USER_DN")

	user_fqdn_aname = GetEnvVar("LDAP_USER_FQDN_ATTRIBUTE_NAME")

	if ldap_tls == "true" {
		ssl_verify = GetEnvVar("LDAP_SSL_VERIFY")
		if ssl_verify == "true" {
			serverName = GetEnvVar("LDAP_SERVER_NAME")
			if serverName == "" {
				log.Error().Msg("LDAP_SERVER_NAME env variable is not set. " +
					"Either set LDAP_SSL_VERIFY=false" +
					" or set LDAP_SERVER_NAME for certificate validation")
			}
		}
	}
}

func GetLdapUser(username string, password string) (map[string]string, int, error) {
	// Dial up with LDAP server and bind read only user
	var ldap_user_dial *ldap.Conn

	// dial up with the server
	if ldap_tls != "true" {
		ldap_dial, err := ldap.Dial("tcp", ldap_string)
		if err != nil {
			log.Debug().Err(err).
				Msgf("Error occurred while dialing to the LDAP server, URL: %s", ldap_string)
			return nil, http.StatusInternalServerError, errors.New(err.Error())
		}
		ldap_user_dial = ldap_dial
	} else {
		var ldap_dial *ldap.Conn
		var err error
		if ssl_verify == "true" {
			ldap_dial, err = ldap.DialTLS("tcp", ldap_string, &tls.Config{ServerName: serverName})
		} else {
			ldap_dial, err = ldap.DialTLS("tcp", ldap_string, &tls.Config{InsecureSkipVerify: true})
		}
		if err != nil {
			log.Debug().Err(err).
				Msgf("Error occurred while dialing to the LDAP server, URL: %s", ldap_string)
			return nil, http.StatusInternalServerError, errors.New(err.Error())
		}
		ldap_user_dial = ldap_dial
	}

	// Load the USER DN fron ENV
	user_bind := fmt.Sprintf("uid=%s,%s", username, userdn)
	// Bind as the user to verify their password
	err := ldap_user_dial.Bind(user_bind, password)
	if err != nil {
		log.Debug().Err(err).
			Msgf("Error occurred while binding user with LDAP server, userdn:%s", userdn)
		defer ldap_user_dial.Close()
		return nil, http.StatusUnauthorized, errors.New(err.Error())
	}

	searchRequest := ldap.NewSearchRequest(
		userdn, // The base dn to search
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(objectClass=organizationalPerson)(uid=%s))", username), // The filter to apply
		[]string{user_fqdn_aname}, // A list attributes to retrieve
		nil,
	)

	sr, err := ldap_user_dial.Search(searchRequest)
	if err != nil {
		log.Debug().Err(err).Msg("Error occurred while searching user with LDAP server")
		defer ldap_user_dial.Close()
		return nil, http.StatusInternalServerError,
			errors.New("failed to authenticate user, internal server error")
	}
	user_fqdn_value := sr.Entries[0].GetAttributeValue("mail")

	results := map[string]string{
		"email": user_fqdn_value,
	}
	defer ldap_user_dial.Close()
	return results, http.StatusOK, nil
}
