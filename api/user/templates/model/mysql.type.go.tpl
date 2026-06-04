{{- $short := (shortname .Name "err" "res" "sqlStr" "db" "XOLog") -}}
{{- $table := (schema .Schema .Table.TableName) -}}
{{- if .Comment -}}
// {{ .Comment }}
{{- else -}}
// {{ .Name }} represents a row from '{{ $table }}'.
{{- end }}
type {{ .Name }} struct {
{{- range .Fields }}
        {{ .Name }} {{ retype .Type }} `json:"{{ .Col.ColumnName }}"`
{{- end }}
}

// {{ .Name }}Columns defines and stores column names for '{{ $table }}'.
type {{ .Name }}Columns struct {
{{- range .Fields }}
    {{ .Name }} string
{{- end }}
}

// Get table name
func Get{{  .Name  }}TableName() string {
    return "{{  .Table.TableName  }}"
}

// Get field string
func Get{{ .Name }}FieldString () string {
    return `{{ colnames .Fields }}`
}

// Get field string slice
func Get{{ .Name }}FieldStringSlice () []string {
    return strings.Split(`{{colnames .Fields .PrimaryKey.Name }}`,",")
}

// Get add field
func Get{{ .Name }}AddField ( {{ $short }}  *{{ .Name }}) []interface{} {
    return [] interface{} {
            		{{- range .Fields }}
            		    {{- if not .Col.IsPrimaryKey }}
                            {{ $short }}.{{ .Name }},
                        {{- end }}
                    {{- end }}
            }
}

// Get scan field
func Get{{ .Name }}ScanField ( {{ $short }}  *{{ .Name }}) []interface{} {
    return [] interface{} {
    		{{- range .Fields }}
                    &{{ $short }}.{{ .Name }},
            {{- end }}
    	}
}

// Get{{ .Name }}Columns get {{ .Name }}Columns
func Get{{ .Name }}Columns () {{ .Name }}Columns {
    return {{ .Name }}Columns {
    		{{- range .Fields }}
    		    {{ .Name }}: "{{ .Col.ColumnName }}",
    		{{- end }}
    	}
}
