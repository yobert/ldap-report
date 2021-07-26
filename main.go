package main

import (
	"fmt"
	"os"

	"github.com/go-ldap/ldap/v3"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {

	config, err := parseConfigFile("config.json")
	if err != nil {
		return err
	}

	data, err := loaddata("data.json")
	if err != nil {
		fmt.Println(err)

		// load new data
		fmt.Print("connecting to " + config.Server + ": ")
		conn, err := ldap.Dial("tcp", config.Server)
		if err != nil {
			return err
		}
		fmt.Println("OK")

		fmt.Print("unauthenticated bind: ")
		if err := conn.UnauthenticatedBind(""); err != nil {
			return err
		}
		fmt.Println("OK")

		fmt.Print(" administrative bind: ")
		if err := conn.Bind(config.Bind, config.Password); err != nil {
			return err
		}
		fmt.Println("OK")

		g, err := groups(config, conn)
		if err != nil {
			return err
		}

		u, err := users(config, conn, g)
		if err != nil {
			return err
		}

		fmt.Print("unauthenticated bind: ")
		if err := conn.UnauthenticatedBind(""); err != nil {
			return err
		}
		fmt.Println("OK")

		data = &Data{
			Users:  u,
			Groups: g,
		}

		if err := savedata("data.json", data); err != nil {
			return err
		}
	}

	if err := report(data); err != nil {
		return err
	}

	return nil
}

var (
	blank  = "   "
	branch = " ├─ "
	carry  = " │  "
	leaf   = " └─ "
)
