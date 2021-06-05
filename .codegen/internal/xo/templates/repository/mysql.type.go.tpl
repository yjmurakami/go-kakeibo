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

	ctx, cancel := context.WithTimeout(context.Background(), database.QueryTimeout)
	defer cancel()

	rows, err := db.QueryContext(ctx, query)
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
			{{ colnames .Fields "Version" }}
		) VALUES (
			{{ colvals .Fields "Version" }}
		)
	`

	ctx, cancel := context.WithTimeout(context.Background(), database.QueryTimeout)
	defer cancel()

	_, err := db.ExecContext(ctx, query, {{ fieldnames .Fields "e" "Version" }})
	if err != nil {
		return err
	}
{{ else -}}
	query := `
		INSERT INTO {{ $table }} (
			{{ colnames .Fields .PrimaryKey.Name "Version" }}
		) VALUES (
			{{ colvals .Fields .PrimaryKey.Name "Version" }}
		)
	`

	ctx, cancel := context.WithTimeout(context.Background(), database.QueryTimeout)
	defer cancel()

	res, err := db.ExecContext(ctx, query, {{ fieldnames .Fields "e" .PrimaryKey.Name "Version" }})
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

	ctx, cancel := context.WithTimeout(context.Background(), database.QueryTimeout)
	defer cancel()

	result, err := db.ExecContext(ctx, query, {{ fieldnamesmulti .Fields "e" .PrimaryKeyFields }}, {{ fieldnames .PrimaryKeyFields "e"}})
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
{{- else -}}
	query := `
		UPDATE {{ $table }} SET
			{{ colnamesquery .Fields ", " .PrimaryKey.Name "Version" }}, version = version + 1
		WHERE {{ colname .PrimaryKey.Col }} = ? AND version = ?
	`

	ctx, cancel := context.WithTimeout(context.Background(), database.QueryTimeout)
	defer cancel()

	result, err := db.ExecContext(ctx, query, {{ fieldnames .Fields "e" .PrimaryKey.Name "Version" }}, e.{{ .PrimaryKey.Name }}, e.Version)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	e.Version += 1
	return nil
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

	ctx, cancel := context.WithTimeout(context.Background(), database.QueryTimeout)
	defer cancel()

	result, err := db.ExecContext(ctx, query, {{ fieldnames .PrimaryKeyFields "e" }})
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
{{- else -}}
	query := `
		DELETE FROM {{ $table }}
		WHERE {{ colname .PrimaryKey.Col }} = ?
	`

	ctx, cancel := context.WithTimeout(context.Background(), database.QueryTimeout)
	defer cancel()

	result, err := db.ExecContext(ctx, query, e.{{ .PrimaryKey.Name }})
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
{{ end -}}
	return nil
}
{{- end }}
