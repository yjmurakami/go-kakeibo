{{- if .Comment -}}
// {{ .Comment }}
{{- end }}
func {{ .Name }} (db database.DB{{ range .QueryUniqueParams }}, {{ .Name }} {{ .Type }}{{ end }}) ({{ if not .OnlyOne }}[]{{ end }}*{{ .Type.Name }}, error) {
	query := `
		{{- range $i, $l := .Query }}
		{{ $l }}{{ end }}
	`

	ctx, cancel := context.WithTimeout(context.Background(), database.QueryTimeout)
	defer cancel()

{{ if .OnlyOne }}
	d := {{ .Type.Name }}{}
	err := db.QueryRowContext(ctx, query{{ range .QueryParams }}, {{ .Name }}{{ end }}).Scan({{ fieldnames .Type.Fields "&d" }})
	if err != nil {
		return nil, err
	}
	return &d, nil
{{- else }}
	rows, err := db.QueryContext(ctx, query{{ range .QueryParams }}, {{ .Name }}{{ end }})
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	s := []*{{ .Type.Name }}{}
	for rows.Next() {
		d := {{ .Type.Name }}{}
		err = rows.Scan({{ fieldnames .Type.Fields "&d" }})
		if err != nil {
			return nil, err
		}
		s = append(s, &d)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return s, nil
{{- end }}
}
