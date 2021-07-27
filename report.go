package main

import (
	"fmt"
	"html/template"
	"os"
	"sort"
	"strings"

	"github.com/go-ldap/ldap/v3"
)

var tmpl = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<title>LDAP report</title>
		<style>
html, body {
	margin: 0;
	padding: 0;
	font-family: sans-serif;
}
.children {
	margin-left: 4rem;
}
.children_closed {
	display: none;
}
.name_toggle {
	display: block;
	color: black;
	text-decoration: none;
}
.plus {
	font-family: monospace;
}
a > .name {
	color: #05a;
}
.count {
	opacity: 0.5;
}
.count_zero {
	display: none;
}
		</style>
		<script>
function toggle(id) {
	var plus = document.getElementById("plus_"+id)
	if(!plus)
		return

	var div = document.getElementById("children_"+id)
	if(!div)
		return

	if(div.className == "children children_closed") {
		div.className = "children"
		plus.innerHTML = "-"
	} else {
		div.className = "children children_closed"
		plus.innerHTML = "+"
	}
}
		</script>
	</head>
	<body>

{{ define "node" }}
{{ if .DN }}
<div class="node">

{{ if .Children }}
{{ if gt (len .Children) 5 }}
	<a class="name_toggle" title="{{ .DN }}" href="#" onclick="toggle({{ .ID }}); return false"><span class="plus" id="plus_{{ .ID }}">+</span> <span class="name">{{ .Name }}</span> <span class="count">{{ len .Children }}</span></a>
{{ else }}
	<a class="name_toggle" title="{{ .DN }}" href="#" onclick="toggle({{ .ID }}); return false"><span class="plus" id="plus_{{ .ID }}">-</span> <span class="name">{{ .Name }}</span> <span class="count">{{ len .Children }}</span></a>
{{ end }}
{{ else }}
	<span class="name_toggle" title="{{ .DN }}" href="#" onclick="toggle({{ .ID }}); return false"><span class="plus" id="plus_{{ .ID }}">&nbsp;</span> <span class="name">{{ .Name }}</span> <span class="count count_zero">{{ len .Children }}</span></span>
{{ end }}

{{ if .Children }}

{{ if gt (len .Children) 5 }}
<div class="children children_closed" id="children_{{ .ID }}">
{{ else }}
<div class="children" id="children_{{ .ID }}">
{{ end }}

{{ range .Children }}
{{ template "node" . }}
{{ end }}

</div>
{{ end }}

</div>
{{ else }}
{{ range .Children }}
{{ template "node" . }}
{{ end }}
{{ end }}
{{ end }}

{{ template "node" . }}

	</body>
</html>
`))

type tmplData struct {
	ID          string
	DN          string
	SubDN       string
	Name        string
	Children    []*tmplData
	ChildrenMap map[string]*tmplData
}

func sortdata(data *tmplData) {
	keys := make([]string, 0, len(data.ChildrenMap))
	children := make([]*tmplData, 0, len(data.ChildrenMap))
	for key := range data.ChildrenMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		child := data.ChildrenMap[key]
		children = append(children, child)
		sortdata(child)
	}
	data.Children = children
}

func report(data *Data) error {
	fh, err := os.Create("report.html")
	if err != nil {
		return err
	}
	defer fh.Close()

	id := 0

	root := &tmplData{
		ChildrenMap: map[string]*tmplData{},
	}

	for _, group := range data.Groups {
		dn, err := ldap.ParseDN(group.DN)
		if err != nil {
			return err
		}

		tree := root
		treedn := ""

		for i := len(dn.RDNs) - 1; i >= 0; i-- {
			rdn := dn.RDNs[i]
			attr := rdn.Attributes[0]
			subdn := attr.Type + "=" + ldap.EscapeFilter(attr.Value)
			key := strings.ToLower(subdn)

			if treedn == "" {
				treedn = subdn
			} else {
				treedn = subdn + "," + treedn
			}

			if tree.ChildrenMap[key] == nil {
				node := &tmplData{
					ID:          fmt.Sprintf("node%d", id),
					DN:          treedn,
					SubDN:       subdn,
					Name:        attr.Value,
					ChildrenMap: map[string]*tmplData{},
				}
				id++
				tree.ChildrenMap[key] = node
			}

			tree = tree.ChildrenMap[key]
		}

		for _, userkey := range group.Users {
			user := data.Users[userkey]
			key := strings.ToLower(user.Name)

			if tree.ChildrenMap[key] == nil {
				node := &tmplData{
					ID:          fmt.Sprintf("node%d", id),
					DN:          user.DN,
					SubDN:       user.DN,
					Name:        user.Name,
					ChildrenMap: map[string]*tmplData{},
				}
				id++
				tree.ChildrenMap[key] = node
			}
		}
	}

	/*	var groupkeys []string
		for key, group := range groups {
			sort.Strings(group.Users)
			groupkeys = append(groupkeys, key)
		}
		sort.Strings(groupkeys)

		var userkeys []string
		for key, user := range users {
			sort.Strings(user.Groups)
			userkeys = append(userkeys, key)
		}
		sort.Strings(userkeys)

		tmplGroupIndex := map[string]*tmplDataGroup{}

		for _, key := range groupkeys {
			group := groups[key]
			tgroup := &tmplDataGroup{
				DN:   group.DN,
				Name: group.Name,
			}
			tmplGroupIndex[key] = tgroup
			tdata.Groups = append(tdata.Groups, tgroup)
		}

		for _, key := range userkeys {
			user := users[key]
			tuser := &tmplDataUser{
				DN:   user.DN,
				Name: user.Name,
			}
			for _, gk := range user.Groups {
				tg := tmplGroupIndex[gk]
				tg.Users = append(tg.Users, tuser)
				tuser.Groups = append(tuser.Groups, tg)
			}
			tdata.Users = append(tdata.Users, tuser)
		}
	*/

	sortdata(root)

	if err := tmpl.Execute(fh, root); err != nil {
		return err
	}

	if err := fh.Close(); err != nil {
		return err
	}
	return nil
	/*
		tdata.Title = "Users"

		fh2, err := os.Create("users.html")
		if err != nil {
			return err
		}
		defer fh2.Close()

		if err := tmpl2.Execute(fh2, tdata); err != nil {
			return err
		}

		return fh2.Close()*/
}
