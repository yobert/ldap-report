package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/go-ldap/ldap/v3"
)

type Data struct {
	Users  map[string]*User  `json:"users"`
	Groups map[string]*Group `json:"groups"`
}

type Attr struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

func dnToKey(in string) (string, []Attr, error) {
	dn, err := ldap.ParseDN(in)
	if err != nil {
		return "", nil, err
	}

	key := ""
	var attrs []Attr

	for i := len(dn.RDNs) - 1; i >= 0; i-- {
		rdn := dn.RDNs[i]
		if len(rdn.Attributes) != 1 {
			return "", nil, fmt.Errorf("Not sure how to handle RDN with %d attributes from DN %#v", len(rdn.Attributes), in)
		}
		attr := rdn.Attributes[0]

		// We'll reverse the DN chunks just to make things easier in the JS
		attrs = append(attrs, Attr{
			Type:  attr.Type,
			Value: attr.Value,
		})

		if strings.IndexByte(attr.Value, '/') != -1 {
			return "", nil, fmt.Errorf("Bad character in DN %#v", in)
		}

		if len(key) > 0 {
			key += "/" + attr.Value
		} else {
			key = attr.Value
		}
	}

	return key, attrs, nil
}

func loaddata(path string) (*Data, error) {
	fh, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fh.Close()
	dec := json.NewDecoder(fh)
	var r Data
	if err := dec.Decode(&r); err != nil {
		return nil, err
	}
	return &r, nil
}
func savedata(path string, data *Data) error {
	fh, err := os.Create(path)
	if err != nil {
		return err
	}
	defer fh.Close()
	enc := json.NewEncoder(fh)
	if err := enc.Encode(data); err != nil {
		return err
	}
	return fh.Close()
}

func printtree(tree map[string]interface{}) {
	printtreei(" ▄", tree, "", "")
}
func printtreei(title string, tree map[string]interface{}, pre1, pre2 string) {
	//title = strings.TrimPrefix(title, "CN=")
	fmt.Printf("%s%s\n", pre1, title)

	var keys []string
	for k := range tree {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for i, k := range keys {
		v := tree[k]
		vv, ok := v.(map[string]interface{})
		if !ok {
			panic("wtf")
		}

		if i+1 == len(keys) {
			printtreei(k, vv, pre2+leaf, pre2+blank)
		} else {
			printtreei(k, vv, pre2+branch, pre2+carry)
		}
	}
}
