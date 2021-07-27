package main

import (
	"fmt"
	"strings"

	"github.com/go-ldap/ldap/v3"
	"github.com/yobert/progress"
)

type User struct {
	DN     string   `json:"dn"`
	Parsed []Attr   `json:"parsed"`
	Name   string   `json:"name"`
	Key    string   `json:"key"`
	Groups []string `json:"groups"`
}

func (user *User) String() string {
	return fmt.Sprintf("user %#v (key %#v, DN %#v)", user.Name, user.Key, user.DN)
}

func users(config *Config, conn *ldap.Conn, groups map[string]*Group) (map[string]*User, error) {
	search := ldap.NewSearchRequest(config.UsersDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		config.UsersFilter,
		nil,
		nil)

	results, err := conn.SearchWithPaging(search, 500)
	if err != nil {
		return nil, err
	}

	r := map[string]*User{}

	bar := progress.NewBar(len(results.Entries), "users")

	for _, result := range results.Entries {
		key, parsed, err := dnToKey(result.DN)
		if err != nil {
			return nil, err
		}

		if dupval, dup := r[key]; dup {
			return nil, fmt.Errorf("Key %#v from DN %#s duplicates existing DN %#v", key, result.DN, dupval.DN)
		}

		user := &User{
			DN:     result.DN,
			Parsed: parsed,
			Key:    key,
		}

		dn, err := ldap.ParseDN(result.DN)
		if err != nil {
			return nil, err
		}

		if len(dn.RDNs) == 0 {
			return nil, fmt.Errorf("Empty RDNs from DN %#v", result.DN)
		}

		rdn := dn.RDNs[0]
		if len(rdn.Attributes) != 1 {
			return nil, fmt.Errorf("Not sure how to handle RDN with %d attributes from DN %#v", len(rdn.Attributes), result.DN)
		}
		attr := rdn.Attributes[0]

		user.Name = attr.Value

		r[key] = user

		// try to query all the groups of the user
		if !strings.Contains(result.DN, "#") {
			groupsearch := ldap.NewSearchRequest(config.GroupsDN,
				ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
				fmt.Sprintf(config.GroupsFilter, ldap.EscapeFilter(result.DN)),
				nil, nil)
			groupresults, err := conn.SearchWithPaging(groupsearch, 500)
			if err != nil {
				return nil, fmt.Errorf("%w handling DN: %s", err, result.DN)
			}
			for _, groupresult := range groupresults.Entries {
				groupkey, _, err := dnToKey(groupresult.DN)
				if err != nil {
					return nil, err
				}

				group, ok := groups[groupkey]
				if !ok {
					return nil, fmt.Errorf("User group DN %#v key %#v did not match any groups", groupresult.DN, groupkey)
				}

				user.Groups = append(user.Groups, groupkey)
				group.Users = append(group.Users, key)
			}
		}

		bar.Next()
	}
	bar.Done()

	fmt.Println(len(results.Entries), "users")
	return r, nil
}
