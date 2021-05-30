{{- $table := (schema .Schema .Table.TableName) -}}
{{- if .Comment -}}
// {{ .Comment }}
{{- end }}
// Generated from {{ .Name }}.sql

type {{ .Name }} struct {
{{- range .Fields }}
	{{ .Name }} {{ if or (eq .Type "uint") (eq .Type "uint64") }}int{{ else }}{{ retype .Type }}{{ end }} // {{ .Col.ColumnName }}
{{- end }}
}
