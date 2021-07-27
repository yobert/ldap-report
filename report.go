package main

import (
	"html/template"
	"os"
	"sort"
)

var tmpl = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<title>{{ .Title }}</title>
		<style>
html, body {
	margin: 0;
	padding: 0;
}
.flex {
	display: flex;
	flex-wrap: wrap;
}
.group {
	flex-grow: 1;
	border: 1px solid black;
	border-radius: 1rem;
	margin: 1rem 1rem 0 0;
}
.groupname {
	font-weight: bold;
	padding: 0.25rem 0.5rem;
}
.groupusers {
	border-top: 1px solid black;
	max-height: 20rem;
	overflow-y: auto;
	padding: 0.25rem 0.5rem;
}
		</style>
		<script>
		</script>
	</head>
	<body>

<div class="flex">
{{ range .Groups }}
	<div class="group">
		<div class="groupname" title="{{ .DN }}">{{ .Name }}</div>
{{ if .Users }}
		<div class="groupusers">
{{ range .Users }}
			<div class="groupuser" title="{{ .DN }}">{{ .Name }}</div>
{{ end }}
		</div>
{{ end }}
	</div>
{{ end }}
</div>

	</body>
</html>
`))

var tmpl2 = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<title>{{ .Title }}</title>
		<style>
html, body {
	margin: 0;
	padding: 0;
}
.flex {
	display: flex;
	flex-wrap: wrap;
}
.group {
	flex-grow: 1;
	border: 1px solid black;
	border-radius: 1rem;
	margin: 1rem 1rem 0 0;
}
.groupname {
	font-weight: bold;
	padding: 0.25rem 0.5rem;
}
.groupusers {
	border-top: 1px solid black;
	max-height: 20rem;
	overflow-y: auto;
	padding: 0.25rem 0.5rem;
}
		</style>
		<script>
		</script>
	</head>
	<body>

<div class="flex">
{{ range .Users }}
	<div class="group">
		<div class="groupname" title="{{ .DN }}">{{ .Name }}</div>
{{ if .Groups }}
		<div class="groupusers">
{{ range .Groups }}
			<div class="groupuser" title="{{ .DN }}">{{ .Name }}</div>
{{ end }}
		</div>
{{ end }}
	</div>
{{ end }}
</div>
	</body>
</html>
`))

type tmplData struct {
	Title  string
	Users  []*tmplDataUser
	Groups []*tmplDataGroup
}
type tmplDataUser struct {
	DN     string
	Name   string
	Groups []*tmplDataGroup
}
type tmplDataGroup struct {
	DN    string
	Name  string
	Users []*tmplDataUser
}

func report(data *Data) error {
	users := data.Users
	groups := data.Groups

	tdata := tmplData{}

	tdata.Title = "Groups"

	fh, err := os.Create("groups.html")
	if err != nil {
		return err
	}
	defer fh.Close()

	var groupkeys []string
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

	if err := tmpl.Execute(fh, tdata); err != nil {
		return err
	}

	if err := fh.Close(); err != nil {
		return err
	}

	tdata.Title = "Users"

	fh2, err := os.Create("users.html")
	if err != nil {
		return err
	}
	defer fh2.Close()

	if err := tmpl2.Execute(fh2, tdata); err != nil {
		return err
	}

	return fh2.Close()
}
