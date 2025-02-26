package typescript

const headTemplate = `
{{- define "head" -}}
import * as flatbuffers from '../../../../Plugins/FlatBuffers';
import { BaseData_Fbs } from '../../BaseData';
import { cacheObjRes, cacheStrRes, FbsData, FbsDataList } from '../../FbsDataListView';
import { {{ .Table.ConfigName | toUpperCamelCase }} } from '../FbsCls/{{ .Table.Name | toKebabCase }}';
import { {{ .Table.ConfigName | toUpperCamelCase }}List } from '../FbsCls/{{ .Table.Name | toKebabCase }}-list';
{{- end -}}
`

const baseClassTemplate = `
{{- define "base_class" -}}
export class {{ .Table.Name | toUpperCamelCase }} extends BaseData_Fbs<Cfg{{ .Table.Name | toUpperCamelCase }}> {
    public getFbsDataList(data: ArrayBuffer): FbsDataList<{{ .Table.InnerConfigName }}> {
        return new {{ .Table.InnerConfigName }}List(data);
    }

	{{- $pkFields := .Table.GetPrimaryKeyFields }}
	{{ $pkLen := $pkFields | len | add -1 }}
    get({{ range $pkIndex, $pkField := $pkFields }}{{ $pkField.Name | toLowerCamelCase }}: {{ $pkField.Type }}{{ if lt $pkIndex $pkLen }}, {{ end }}{{ end }}): {{ .Table.InnerConfigName }} {
	    return super.get({{ range $pkIndex, $pkField := $pkFields }}{{ $pkField.Name | toLowerCamelCase }}{{ if lt $pkIndex $pkLen }}, {{ end }}{{ end }});
    }
}
{{- end -}}
`

const innerClass1Template = `
{{- define "inner_class1" -}}
export class {{ .Table.InnerConfigName }} extends FbsData {
    private _fbs: {{ .Table.ConfigName | toUpperCamelCase }};

    __init(fbs: {{ .Table.ConfigName | toUpperCamelCase }}) {
		this._fbs = fbs;
    }

    public clone(): {{ .Table.InnerConfigName }} {
	    let newFD: {{ .Table.InnerConfigName }} = new {{ .Table.InnerConfigName }}();
	    newFD._fbs = new {{ .Table.ConfigName | toUpperCamelCase }}();
        newFD._fbs.bb = this._fbs.bb;
        newFD._fbs.bb_pos = this._fbs.bb_pos;
	    return newFD;
    }

	{{- $pkFields := .Table.GetPrimaryKeyFields }}
	{{ $pkLen := $pkFields | len | add -1 }}
    public checkKeyPrefix(): string{
        return {{ "\x60" }}{{ .Table.InnerConfigName }}.{{ range $pkIndex, $pkField := $pkFields }}${this.{{ $pkField.Name | toLowerCamelCase }}}{{ if lt $pkIndex $pkLen }}.{{ end }}{{ end }}{{ "\x60" }};
    }

    {{ range $index, $field := .Table.FieldRowIter }}
	{{ $field.Type.Decorator }}
	public get {{ $field.Name | toLowerCamelCase }}(): {{ $field.Type }} {
		return this._fbs.{{ $field.Name | toLowerCamelCase }}();
	}
    {{ end }}
}
{{- end -}}
`
const innerClass2Template = `
{{- define "inner_class2" -}}
export class {{ .Table.InnerConfigName }}List extends FbsDataList<{{ .Table.InnerConfigName }}> {
	private _fbsList: {{ .Table.ConfigName | toUpperCamelCase }}List;

	__init(data: ArrayBuffer) {
		this._fbsList = {{ .Table.ConfigName | toUpperCamelCase }}List.getRootAs{{ .Table.ConfigName | toUpperCamelCase }}List(new flatbuffers.ByteBuffer(new Uint8Array(data)));
	}

    public get length(): number {
        return this._fbsList.{{ .Table.Name | toLowerCamelCase }}Length();
    }

    public getFbsData(index: number, obj?: {{ .Table.InnerConfigName | toUpperCamelCase }}): {{ .Table.InnerConfigName | toUpperCamelCase }} {
        obj = obj ? obj : new {{ .Table.InnerConfigName | toUpperCamelCase }}();
        obj.__init(this._fbsList.{{ .Table.Name | toLowerCamelCase }}(index) as {{ .Table.ConfigName | toUpperCamelCase}});
        return obj;
    }
}
{{- end -}}
`

const enumTemplate = `
{{- define "enum" -}}
{{ range $_, $macro := .Table.GetMacroDecorators }}
// {{ $macro.KeyField.Comment }}
export enum {{ $macro.MacroName | toUpperCamelCase }}Enum { 
	{{- range $_, $macroDetail := $macro.List }}
    /** {{ $macroDetail.Comment }} */
    {{ $macroDetail.Key | $macro.KeyField.Type.Convert | toUpperSnakeCase }} = {{ $macroDetail.Value | $macro.ValueField.Type.Convert }},
	{{- end }}
}
{{ end }}
{{- end -}}
`

const tsTemplate = `
{{- template "head" .}}

{{ template "base_class" .}}

{{ template "inner_class1" .}}

{{ template "inner_class2" .}}

{{ template "enum" . -}}
`
