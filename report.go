package main

import (
	"fmt"
	"os"
	"sort"
)

func report(data *Data) error {
	users := data.Users
	groups := data.Groups

	fh, err := os.Create("report.html")
	if err != nil {
		return err
	}
	defer fh.Close()

	var groupkeys []string
	for key, group := range groups {
		//sort.Sort(group.Users)
		sort.Strings(group.Users)
		groupkeys = append(groupkeys, key)
	}
	sort.Strings(groupkeys)

	var userkeys []string
	for key, user := range users {
		//sort.Sort(user.Groups)
		sort.Strings(user.Groups)
		userkeys = append(userkeys, key)
	}

	for _, key := range groupkeys {
		group := groups[key]
		fmt.Println(key, group)
	}

	fmt.Println("------")

	for _, key := range userkeys {
		user := users[key]
		fmt.Println(key, user)
	}

	//printtree(groups)

	//		dn, err := ldap.ParseDN(result.DN)

	return fh.Close()
}
