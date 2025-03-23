package flatbuffers

const packageTemplate = `
{{- define "PACKAGE_TEMPLATE" -}}
namespace {{ .Table.GetNamespace }};
{{- end -}}
`

const dataSetTemplate = `
{{- define "TABLE_SCHEMA_TEMPLATE" -}}
table {{ .Table.Name | toUpperCamelCase }}{
    {{- range $_, $field := .Table.FieldRowIter }}
    {{ $field.Name | toLowerCamelCase }}: {{ $field.Type }};
	{{- end }}
}
{{- end -}}
`

const tailTemplate = `
{{- define "ROOT_TYPE_TEMPLATE" -}}
table {{ .Table.ConfigName }}DetailList{
    {{ .Table.Name | toLowerCamelCase }}DetailList: [{{ .Table.Name | toUpperCamelCase }}];
}

root_type {{ .Table.ConfigName }}DetailList;
{{- end -}}
`

const fbTemplate = `
{{- template "PACKAGE_TEMPLATE" . }}

{{ template "TABLE_SCHEMA_TEMPLATE" . }}

{{ template "ROOT_TYPE_TEMPLATE" . }}`
