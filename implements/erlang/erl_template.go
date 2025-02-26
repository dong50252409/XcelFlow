package erlang

const erlHeadTemplate = `
{{- define "head" -}}
%%Auto Create, Don't Edit%%
-module({{ .Table.ConfigName | toSnakeCase }}).
-include("{{ .Table.ConfigName | toSnakeCase }}.hrl").
-compile(export_all).
-compile(nowarn_export_all).
-compile({no_auto_import, [get/1]}).
{{ end }}
`

const erlGetTemplate = `
{{- define "get" -}}
{{/* 声明模板渲染所需的变量 */}}
{{- $configName := .Table.ConfigName | toLower -}}
{{- $dataSetIter := .Table.DataSetIter -}}
{{- $fieldRowIter := .Table.FieldRowIter -}}
{{- $dsLastIndex := .Table.FieldLen | add -1 -}}
{{- $pkValuesList := .Table.GetPrimaryKeyValuesByString -}}
{{- $pkLen := index $pkValuesList 0 | len -}}
{{- $pkLastIndex := $pkLen | add -1 -}}
{{- $pkSeq := seq $pkLen -}}
{{- range $rowIndex, $dataRow := $dataSetIter -}}
get({{ index $pkValuesList $rowIndex | joinByComma }}) ->
    #{{ $configName | toSnakeCase }}{
        {{- range $fieldIndex, $field := $fieldRowIter }}
        {{ $field.Name | toSnakeCase }} = {{ index $dataRow $field.Column | $field.Convert }}{{ if lt $fieldIndex $dsLastIndex }},{{ end }}
        {{- end }}
    };
{{ end -}}

get({{ range $pkIndex, $_ := $pkSeq }}_ID{{ $pkIndex }}{{ if lt $pkIndex $pkLastIndex }}, {{ end }}{{ end }}) ->
    throw({error_config, ?MODULE, {{ range $pkIndex, $_ := $pkSeq }}_ID{{ $pkIndex }}{{ if lt $pkIndex $pkLastIndex }}, {{ end }}{{ end }}}).
{{ end }}
`
const erlListTemplate = `
{{- define "list" -}}
{{/* 声明模板渲染所需的变量 */}}
{{- $pkValuesList := .Table.GetPrimaryKeyValuesByString -}}
{{- $pkLastIndex := len $pkValuesList | add -1 -}}
{{- $pkLen := .Table.GetPrimaryKeyFields | len -}}
{{- if eq $pkLen 1 -}}
list() ->
    [{{- range $pkIndex, $pkValues := $pkValuesList -}}{{ $pkValues | joinByComma }}{{ if lt $pkIndex $pkLastIndex }}, {{ end }}{{- end -}}].
{{- else -}}
list() ->
    [{{- range $pkIndex, $pkValues := $pkValuesList -}}{{ "{" }}{{ $pkValues | joinByComma }}{{ "}" }}{{ if lt $pkIndex $pkLastIndex }}, {{ end }}{{- end -}}].
{{ end -}}
{{- end -}}
`

const erlTemplate = `
{{- template "head" .}}
{{ template "get" .}}
{{ template "list" .}}
`
