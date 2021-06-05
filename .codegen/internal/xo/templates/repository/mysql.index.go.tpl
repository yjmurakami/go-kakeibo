{{- $table := (schema .Schema .Type.Table.TableName) -}}
{{- $lowerName := lowerFirstLetter .Type.Name -}}

// Generated from index '{{ .Index.IndexName }}'.
func (r *{{ $lowerName }}Repository) Select{{ .FuncName }}(db database.DB{{ goparamlist .Fields true true }}) ({{ if not .Index.IsUnique }}[]{{ end }}*entity.{{ .Type.Name }}, error) {
	query := `
		SELECT {{ colnames .Type.Fields }}
		FROM {{ $table }}
		WHERE {{ colnamesquery .Fields " AND " }}
		{{- if not .Index.IsUnique }}
		ORDER BY {{ colnames .Type.PrimaryKeyFields }}{{ end }}
	`

	ctx, cancel := context.WithTimeout(context.Background(), database.QueryTimeout)
	defer cancel()

{{- if .Index.IsUnique }}
	e := entity.{{ .Type.Name }}{}
	err := db.QueryRowContext(ctx, query{{ goparamlist .Fields true false }}).Scan({{ fieldnames .Type.Fields "&e" }})
	if err != nil {
		return nil, err
	}
	return &e, nil
{{- else }}
	rows, err := db.QueryContext(ctx, query{{ goparamlist .Fields true false }})
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	s := []*entity.{{ .Type.Name }}{}
	for rows.Next() {
		e := entity.{{ .Type.Name }}{}
		err = rows.Scan({{ fieldnames .Type.Fields "&e" }})
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
{{- end }}
}
