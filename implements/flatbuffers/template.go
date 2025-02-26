package flatbuffers

const dataSetTemplate = `
{{- define "tableSchema" -}}
table {{ .Table.Name | toUpperCamelCase }}{
    {{- range $_, $field := .Table.FieldRowIter }}
    {{ $field.Name | toLowerCamelCase }}: {{ $field.Type }};
	{{- end }}
}
{{- end -}}
`

const tailTemplate = `
{{- define "rootType" -}}
table {{ .Table.ConfigName }}List{
    {{ .Table.Name | toLowerCamelCase }}List: [{{ .Table.Name | toUpperCamelCase }}];
}

root_type {{ .Table.ConfigName }}List;
{{- end -}}
`

const fbTemplate = `
{{- if ne .Table.Schema.GetNamespace "" -}}
{{"namespace" }} {{ .Table.GetNamespace }};

{{ template "tableSchema" . }}
{{ else }}
{{- template "tableSchema" . }}
{{ end }}

{{ template "rootType" . }}
`
