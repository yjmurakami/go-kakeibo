{{- $table := (schema .Schema .Table.TableName) -}}
// Generated from '{{ $table }}'.
type {{ .Name }} struct {
{{- range .Fields }}
	{{ .Name }} {{ if or (eq .Type "uint") (eq .Type "uint64") }}int{{ else }}{{ retype .Type }}{{ end }} // {{ .Col.ColumnName }}
{{- end }}
}
