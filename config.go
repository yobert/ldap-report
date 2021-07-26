package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	Server   string
	Bind     string
	Password string

	UsersDN     string `json:"users_dn"`
	UsersFilter string `json:"users_filter"`

	GroupsDN     string `json:"groups_dn"`
	GroupsFilter string `json:"groups_filter"`
	GroupsAll    string `json:"groups_all"`
}

func parseConfigFile(path string) (*Config, error) {
	fh, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fh.Close()
	dec := json.NewDecoder(fh)
	var config Config
	if err := dec.Decode(&config); err != nil {
		return nil, err
	}
	return &config, nil
}
