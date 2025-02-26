package erlang

const hrlHeadTemplate = `
{{- define "head" -}}
%% Auto Create, Don't Edit
-ifndef({{ .Table.ConfigName | toUpperSnakeCase }}_HRL).
-define({{ .Table.ConfigName | toUpperSnakeCase }}_HRL, true).
{{- end -}}
`

const hrlRecordTemplate = `
{{- define "record" -}}
{{- $fieldLen := .Table.FieldLen -}}
{{- $lastIndex := add $fieldLen -1 }}
{{- if ne $fieldLen 0 }}
-record({{ .Table.ConfigName }}, {
    {{- range $index, $field := .Table.FieldRowIter  }}
    {{ $field.Name | toSnakeCase }} = {{ $field.DefaultValue }} :: {{ $field.Type }}{{ if lt $index $lastIndex }},{{ end }}	% {{ $field.Comment }}
    {{- end }}
}).
{{ end }}
{{- end -}}
`

const hrlMacroTemplate = `
{{- define "macro" -}}
{{- range $_, $macro := .Table.GetMacroDecorators -}}
%% {{ $macro.KeyField.Comment }}
{{ range $_, $macroDetail := $macro.List -}}
-define({{ $macroDetail.Key | $macro.KeyField.Type.Convert | toUpperSnakeCase }}, {{ $macroDetail.Value | $macro.ValueField.Type.Convert }}).	  % {{ $macroDetail.Comment }}
{{ end }}
{{- end -}}
{{- end -}}
`

const hrlTailTemplate = `
{{- define "tail" -}}
-endif.
{{- end -}}
`

const hrlTemplate = `
{{- template "head" . }}
{{ template "record" . }}
{{ template "macro" . }}
{{ template "tail" . }}
`
