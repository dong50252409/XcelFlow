package typescript

const headTemplate = `{{- define "IMPORT_TEMPLATE" -}}
import {flatbuffers, flexbuffers, BaseData_Fbs, cacheObjRes, FbsData, FbsDataList } from './index'
{{- end -}}`

const baseClassTemplate = `{{- define "BASE_CLASS_TEMPLATE" -}}
{{- $configName := .Table.ConfigName | toUpperCamelCase -}}
{{- $pkFields := .Table.GetPrimaryKeyFields -}}
{{- $pkLastIndex := len $pkFields | add -1 -}}
export class {{ $configName }} extends BaseData_Fbs<Cfg{{ $configName }}> {
    protected _keys: string[] = [{{- range $pkIndex, $pkField := $pkFields -}}'{{ $pkField.Name | toLowerCamelCase }}'{{ if lt $pkIndex $pkLastIndex }}, {{ end }}{{ end }}];

    public getFbsDataList(data: ArrayBuffer): FbsDataList<Cfg{{ $configName }}> {
        return new Cfg{{ $configName }}List(data);
    }

    public get({{- range $pkIndex, $pkField := $pkFields -}}{{ $pkField.Name | toLowerCamelCase }}: {{ $pkField.Type }}{{ if lt $pkIndex $pkLastIndex }}, {{ end }}{{ end }}): Cfg{{ .Table.ConfigName | toUpperCamelCase }} {
        return super.get({{- range $pkIndex, $pkField := $pkFields -}}{{ $pkField.Name | toLowerCamelCase }}{{ if lt $pkIndex $pkLastIndex }}, {{ end }}{{ end }});
    }
}
{{- end -}}`

const listClassTemplate = `{{- define "LIST_CLASS_TEMPLATE" -}}
{{- $configName := .Table.ConfigName | toUpperCamelCase -}}
export class Cfg{{ $configName }}List extends FbsDataList<Cfg{{ $configName }}> {
    private bb: flatbuffers.ByteBuffer | null = null;
    private bb_pos = 0;

    __init(data: ArrayBuffer) {
        const bb = new flatbuffers.ByteBuffer(new Uint8Array(data));
        this.bb_pos = bb.readInt32(bb.position()) + bb.position();
        this.bb = bb;
    }

    public get length(): number {
        const offset = this.bb!.__offset(this.bb_pos, 4);
        return offset ? this.bb!.__vector_len(this.bb_pos + offset) : 0;
    }

    public getFbsData(index: number, obj?: Cfg{{ $configName }}): Cfg{{ $configName }} | null {
        const offset = this.bb!.__offset(this.bb_pos, 4);
        obj = obj || new Cfg{{ $configName }}();
        return offset ? obj.__init(this.bb!.__indirect(this.bb!.__vector(this.bb_pos + offset) + index * 4), this.bb!) : obj;
    }
}
{{- end -}}`

const classTemplate = `{{- define "CLASS_TEMPLATE" -}}
{{- $configName := .Table.ConfigName | toUpperCamelCase -}}
export class Cfg{{ $configName }} extends FbsData {
    private bb: flatbuffers.ByteBuffer | null = null;
    private bb_pos = 0;

    __init(i: number, bb: flatbuffers.ByteBuffer): Cfg{{ $configName }} {
        this.bb_pos = i;
        this.bb = bb;
        return this;
    }

    __toObject(i: number): unknown {
        const offset = this.bb!.__offset(this.bb_pos, i);
        if (offset) {
            const array = new Uint8Array(this.bb!.bytes().buffer, this.bb!.bytes().byteOffset + this.bb!.__vector(this.bb_pos + offset), this.bb!.__vector_len(this.bb_pos + offset));
            return flexbuffers.toObject(array.slice().buffer);
        }
        return null
    }

    public clone(): Cfg{{ $configName }} {
        let obj = new Cfg{{ $configName }}();
        if (this.bb) {
            obj.__init(this.bb_pos, this.bb!);
        }
        return obj;
    }

    {{ range $index, $field := .Table.FieldRowIter }}
    {{- $decorator := $field.Type.DecoratorStr }}
    {{- if $decorator }}
    {{ $decorator }}
    {{- end }}
    public get {{ $field.Name | toLowerCamelCase }}(): {{ $field.Type }} {
        {{- if $field.Type.IsReferenceType }}
        return this.__toObject({{ $index | calOffset }}) as {{ $field.Type }}
        {{- else }}
        const offset = this.bb!.__offset(this.bb_pos, {{ $index | calOffset }});
        return offset ? this.bb!.{{ $field.Type.MethodStr }}(this.bb_pos + offset) : {{ $field.Type.DefaultValueStr }};
        {{- end }}
    }
    {{ end }}
}
{{- end -}}`

const enumTemplate = `{{- define "ENUM_TEMPLATE" -}}
{{ range $_, $macro := .Table.GetMacroDecorators }}
// {{ $macro.KeyField.Comment }}
export enum {{ $macro.MacroName | toUpperCamelCase }}Enum { 
	{{- range $_, $macroDetail := $macro.List }}
    /** {{ $macroDetail.Comment }} */
    {{ $macroDetail.Key | toUpperSnakeCase }} = {{ $macroDetail.Value | $macro.ValueField.Type.Convert }},
	{{- end }}
}
{{ end }}
{{- end -}}`

const tsTemplate = `{{- template "IMPORT_TEMPLATE" .}}

{{ template "BASE_CLASS_TEMPLATE" .}}

{{ template "LIST_CLASS_TEMPLATE" .}}

{{ template "CLASS_TEMPLATE" .}}

{{ template "ENUM_TEMPLATE" . -}}`
