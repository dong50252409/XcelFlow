package erlang

const hrlHeadTemplate = `{{- define "HRL_HEAD_TEMPLATE" -}}
%% Auto Create, Don't Edit
-ifndef({{ .Table.ConfigName | toUpperSnakeCase }}_HRL).
-define({{ .Table.ConfigName | toUpperSnakeCase }}_HRL, true).
{{- end -}}`

const hrlRecordTemplate = `{{- define "HRL_RECORD_TEMPLATE" -}}
{{- $fieldLen := .Table.FieldLen -}}
{{- $lastIndex := add $fieldLen -1 -}}
{{- if ne $fieldLen 0 -}}
-record({{ .Table.ConfigName }}, {
    {{- range $index, $field := .Table.FieldRowIter  }}
    {{ $field.Name | toSnakeCase }} = {{ $field.Type.DefaultValueStr }} :: {{ $field.Type }}{{ if lt $index $lastIndex }},{{ end }}	% {{ $field.Comment }}
    {{- end }}
}).
{{- end -}}
{{- end -}}`

const hrlMacroTemplate = `{{- define "HRL_MACRO_TEMPLATE" -}}
{{- range $_, $macro := .Table.GetMacroDecorators -}}
%% {{ $macro.KeyField.Comment }}
{{ range $_, $macroDetail := $macro.List -}}
-define({{ $macroDetail.Key | toUpperSnakeCase }}, {{ $macroDetail.Value | $macro.ValueField.Type.Convert }}).	  % {{ $macroDetail.Comment }}
{{ end }}
{{- end -}}
{{- end -}}`

const hrlTailTemplate = `{{- define "HRL_TAIL_TEMPLATE" -}}
-endif.
{{- end -}}`

// 默认hrl模板
const hrlTemplate = `{{ template "HRL_HEAD_TEMPLATE" . }}

{{ template "HRL_RECORD_TEMPLATE" . }}

{{ template "HRL_MACRO_TEMPLATE" . }}

{{ template "HRL_TAIL_TEMPLATE" . }}`

const erlHeadTemplate = `{{- define "ERL_HEAD_TEMPLATE" -}}
%%Auto Create, Don't Edit
-module({{ .Table.ConfigName | toSnakeCase }}).
-include("{{ .Table.ConfigName | toSnakeCase }}.hrl").
-compile(export_all).
-compile(nowarn_export_all).
-compile({no_auto_import, [get/1]}).
{{- end -}}`

const erlGetTemplate = `{{- define "ERL_GET_TEMPLATE" -}}
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
{{- end -}}`

const erlListTemplate = `{{- define "ERL_LIST_TEMPLATE" -}}
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
{{- end -}}`

// 默认erl模板
const erlTemplate = `{{ template "ERL_HEAD_TEMPLATE" . }}

{{ template "ERL_GET_TEMPLATE" . }}

{{ template "ERL_LIST_TEMPLATE" . }}`
