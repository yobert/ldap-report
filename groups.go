package main

import (
	"fmt"

	"github.com/go-ldap/ldap/v3"
)

type Group struct {
	DN    string
	Name  string
	Key   string
	Users []string
}

func (group *Group) String() string {
	return fmt.Sprintf("group %#v (key %#v, DN %#v)", group.Name, group.Key, group.DN)
}

func groups(config *Config, conn *ldap.Conn) (map[string]*Group, error) {
	// list groups?
	search := ldap.NewSearchRequest(config.GroupsDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		config.GroupsAll,
		nil, //[]string{"CN"},
		nil)

	results, err := conn.SearchWithPaging(search, 500)
	if err != nil {
		return nil, err
	}

	r := map[string]*Group{}

	for _, result := range results.Entries {
		key, err := dnToKey(result.DN)
		if err != nil {
			return nil, err
		}

		if dupval, dup := r[key]; dup {
			return nil, fmt.Errorf("Key %#v from DN %#s duplicates existing DN %#v", key, result.DN, dupval.DN)
		}

		group := &Group{
			DN:  result.DN,
			Key: key,
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

		group.Name = attr.Value

		r[key] = group
	}

	fmt.Println(len(results.Entries), "groups")
	return r, nil
}
