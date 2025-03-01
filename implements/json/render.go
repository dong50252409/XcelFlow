package json

import (
	"encoding/json"
	"github.com/stoewer/go-strcase"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"xCelFlow/entities"
	"xCelFlow/render"
)

type JSONRender struct {
	*render.Render
}

func init() {
	render.Register("json", newJSONRender)
}

func newJSONRender(render *render.Render) render.IRender {
	return &JSONRender{render}
}

func (r *JSONRender) Execute() error {
	if err := r.Render.ExecuteBefore(); err != nil {
		return err
	}

	jsonDir := r.ExportDir()
	if err := os.MkdirAll(jsonDir, os.ModePerm); err != nil {
		return err
	}

	fp := filepath.Join(jsonDir, r.Filename())
	fileIO, err := os.Create(fp)
	if err != nil {
		return err
	}
	defer func() { _ = fileIO.Close() }()

	rootMap := make(map[string]any)
	if dataList := r.dataListExecute(); len(dataList) > 0 {
		rootMap[r.ConfigName()+"List"] = dataList
	}
	if macroMap := r.macroListExecute(); len(macroMap) > 0 {
		rootMap[r.ConfigName()+"MacroMap"] = macroMap
	}
	if len(rootMap) > 0 {
		if jsonData, err := json.MarshalIndent(rootMap, "", "    "); err == nil {
			_, err = fileIO.Write(jsonData)
		}
	}

	return err
}

func (r *JSONRender) dataListExecute() (dataList []map[string]any) {
	if r.FieldLen == 0 {
		return nil
	}

	dataList = make([]map[string]any, 0, r.DataSetLen)
	for _, rowData := range r.DataSetIter() {
		rowMap := make(map[string]any, r.FieldLen)
		for _, field := range r.FieldRowIter() {
			switch v := rowData[field.Column]; v {
			case nil, "":
				continue
			default:
				v1 := convert(v)
				rowMap[strcase.LowerCamelCase(field.Name)] = v1
			}
		}
		dataList = append(dataList, rowMap)
	}

	sort.Slice(dataList, func(i, j int) bool {
		for _, field := range r.GetPrimaryKeyFields() {
			switch field.Type.Kind() {
			case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				switch v1, v2 := dataList[i][strcase.LowerCamelCase(field.Name)].(int64), dataList[j][strcase.LowerCamelCase(field.Name)].(int64); {
				case v1 < v2:
					return true
				case v1 > v2:
					return false
				}
			case reflect.Float32, reflect.Float64:
				switch v1, v2 := dataList[i][strcase.LowerCamelCase(field.Name)].(float64), dataList[j][strcase.LowerCamelCase(field.Name)].(float64); {
				case v1 < v2:
					return true
				case v1 > v2:
					return false
				}
			case reflect.String:
				switch v1, v2 := dataList[i][strcase.LowerCamelCase(field.Name)].(string), dataList[j][strcase.LowerCamelCase(field.Name)].(string); {
				case v1 < v2:
					return true
				case v1 > v2:
					return false
				}
			case reflect.Interface:
				switch v1, v2 := dataList[i][strcase.LowerCamelCase(field.Name)].(entities.AnyT), dataList[j][strcase.LowerCamelCase(field.Name)].(entities.AnyT); {
				case v1 < v2:
					return true
				case v1 > v2:
					return false
				}
			default:
				return true
			}
		}
		return true
	})

	return dataList
}

func (r *JSONRender) macroListExecute() (macroMap map[string]any) {
	macroMap = make(map[string]any)
	for _, macro := range r.GetMacroDecorators() {
		var childMacroMap = make(map[string]any)
		for _, macroDetail := range macro.List {
			childMacroMap[strcase.UpperSnakeCase(macro.KeyField.Convert(macroDetail.Key))] = macroDetail.Value
		}
		macroMap[strcase.LowerCamelCase(macro.MacroName)] = childMacroMap
	}
	return macroMap
}

func (r *JSONRender) Verify() error {
	return nil
}

func (r *JSONRender) ConfigName() string {
	return strcase.LowerCamelCase(r.Schema.GetTableNamePrefix() + r.Name)
}

func (r *JSONRender) Filename() string {
	return strcase.KebabCase(r.Schema.GetFilePrefix()+r.Name) + ".json"
}
