package main

import (
	"fmt"
	"strings"

	"github.com/go-ldap/ldap/v3"
)

func auth(config *Config, conn *ldap.Conn, user, pass string) error {

	key := "sAMAccountName"
	if strings.IndexByte(user, '@') != -1 {
		key = "userPrincipalName"
	}

	fmt.Println("find user")

	search := ldap.NewSearchRequest(config.UsersDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectclass=user)"+ // object type is user
			fmt.Sprintf("(%s=%s)", key, ldap.EscapeFilter(user))+ // search by username
			"(!userAccountControl:1.2.840.113556.1.4.803:=2)"+ // account is active
			")",
		[]string{
			"sAMAccountName",
			"sn",
			"givenName",
			"displayName",
			"userPrincipalName",
			"memberOf",
			"objectClass",
		},
		nil)

	results, err := conn.Search(search)
	if err != nil {
		return err
	}

	if len(results.Entries) != 1 {
		return fmt.Errorf("Search for %s %#v matched %d results", key, user, len(results.Entries))
	}

	entry := results.Entries[0]

	for _, attr := range entry.Attributes {
		fmt.Println(attr.Name + ":")
		for _, value := range attr.Values {
			fmt.Println("\t" + value)
		}
	}

	fmt.Println(entry.DN)
	if err := conn.Bind(entry.DN, pass); err != nil {
		return err
	}
	return nil
}
