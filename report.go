package main

import (
	"html/template"
	"os"
)

var tmpl = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<title>LDAP report</title>
		<link rel="stylesheet" href="main.css"></link>
		<script src="main.js"></script>
		<script>
let main_data = {{ . }};
		</script>
	</head>
	<body onload="main();">
	</body>
</html>
`))

func report(data *Data) error {
	fh, err := os.Create("report.html")
	if err != nil {
		return err
	}
	defer fh.Close()

	if err := tmpl.Execute(fh, data); err != nil {
		return err
	}

	if err := fh.Close(); err != nil {
		return err
	}
	return nil
}
