{{- $table := (schema .Schema .Table.TableName) -}}
{{- $lowerName := lowerFirstLetter .Name -}}

// Generated from '{{ $table }}'.
type {{ $lowerName }}Repository struct{}

func (r *{{ $lowerName }}Repository) SelectAll(db database.DB) ([]*entity.{{ .Name }}, error) {
	query := `
		SELECT {{ colnames .Fields }}
		FROM {{ $table }}
		ORDER BY {{ colnames .PrimaryKeyFields }}
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	s := []*entity.{{ .Name }}{}
	for rows.Next() {
		e := entity.{{ .Name }}{}
		err = rows.Scan({{ fieldnames .Fields "&e" }})
		if err != nil {
			return nil, err
		}
		s = append(s, &e)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return s, nil
}

{{ if .PrimaryKey }}
func (r *{{ $lowerName }}Repository) Insert(db database.DB, e *entity.{{ .Name }}) error {
{{ if .Table.ManualPk -}}
	query := `
		INSERT INTO {{ $table }} (
			{{ colnames .Fields }}
		) VALUES (
			{{ colvals .Fields }}
		)
	`

	_, err := db.Exec(query, {{ fieldnames .Fields "e" }})
	if err != nil {
		return err
	}
{{ else -}}
	query := `
		INSERT INTO {{ $table }} (
			{{ colnames .Fields .PrimaryKey.Name }}
		) VALUES (
			{{ colvals .Fields .PrimaryKey.Name }}
		)
	`

	res, err := db.Exec(query, {{ fieldnames .Fields "e" .PrimaryKey.Name }})
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	e.{{ .PrimaryKey.Name }} = {{ if or (eq .PrimaryKey.Type "uint") (eq .PrimaryKey.Type "uint64") }}int{{ else }}{{ .PrimaryKey.Type }}{{ end }}(id)
{{ end -}}
	return nil
}

{{ if ne (fieldnamesmulti .Fields "e" .PrimaryKeyFields) "" }}
func (r *{{ $lowerName }}Repository) Update(db database.DB, e *entity.{{ .Name }}) error {
{{ if gt ( len .PrimaryKeyFields ) 1 -}}
	query := `
		UPDATE {{ $table }} SET
			{{ colnamesquerymulti .Fields ", " 0 .PrimaryKeyFields }}
		WHERE {{ colnamesquery .PrimaryKeyFields " AND " }}
	`

	_, err := db.Exec(query, {{ fieldnamesmulti .Fields "e" .PrimaryKeyFields }}, {{ fieldnames .PrimaryKeyFields "e"}})
	return err
{{- else -}}
	query := `
		UPDATE {{ $table }} SET
			{{ colnamesquery .Fields ", " .PrimaryKey.Name }}
		WHERE {{ colname .PrimaryKey.Col }} = ?
	`

	_, err := db.Exec(query, {{ fieldnames .Fields "e" .PrimaryKey.Name }}, e.{{ .PrimaryKey.Name }})
	return err
{{- end }}
}

{{ else }}
	// Update statements omitted due to lack of fields other than primary key
{{ end }}

func (r *{{ $lowerName }}Repository) Delete(db database.DB, e *entity.{{ .Name }}) error {
{{ if gt ( len .PrimaryKeyFields ) 1 -}}
	query := `
		DELETE FROM {{ $table }}
		WHERE {{ colnamesquery .PrimaryKeyFields " AND " }}
	`

	_, err := db.Exec(query, {{ fieldnames .PrimaryKeyFields "e" }})
{{- else -}}
	query := `
		DELETE FROM {{ $table }}
		WHERE {{ colname .PrimaryKey.Col }} = ?
	`

	_, err := db.Exec(query, e.{{ .PrimaryKey.Name }})
{{ end -}}
	return err
}
{{- end }}
