package json

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"sync"
	"xCelFlow/config"
	"xCelFlow/core"
	"xCelFlow/render"

	"github.com/stoewer/go-strcase"
)

type JSONRender struct {
	*render.Render
	*config.JSONSchema
}

var (
	once    sync.Once
	initErr error
)

func init() {
	render.Register("json", newJSONRender)
}

func newJSONRender(render *render.Render) (core.IRender, error) {
	Schema := render.Schema.(*config.JSONSchema)

	if err := instance(Schema); err != nil {
		return nil, fmt.Errorf("初始化失败 %s", err)
	}
	r := &JSONRender{Render: render, JSONSchema: Schema}

	return r, nil
}

func instance(schema *config.JSONSchema) error {
	once.Do(func() {
		if err := os.MkdirAll(schema.GetJsonDirectory(), os.ModePerm); err != nil {
			initErr = fmt.Errorf("导出路径创建失败 %s", err)
			return
		}
	})
	return initErr
}

func (r *JSONRender) Execute() error {
	fp := filepath.Join(r.GetJsonDirectory(), r.Filename())
	fileIO, err := os.Create(fp)
	if err != nil {
		return fmt.Errorf("创建文件失败:%s %s", fp, err)
	}
	defer func() { _ = fileIO.Close() }()

	rootMap := make(map[string]any)
	if dataList := r.dataListExecute(); len(dataList) > 0 {
		rootMap[r.ConfigName()+"DetailList"] = dataList
	}
	if macroMap := r.macroListExecute(); len(macroMap) > 0 {
		rootMap[r.ConfigName()+"MacroMap"] = macroMap
	}
	if len(rootMap) > 0 {
		if jsonData, err := json.MarshalIndent(rootMap, "", "    "); err == nil {
			if _, err = fileIO.Write(jsonData); err != nil {
				return fmt.Errorf("写入文件失败:%s %s", fp, err)
			}
		}
	}

	fmt.Printf("导出配置：%s\n", fp)
	return nil
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
				rowMap[strcase.LowerCamelCase(field.Name)] = convert(field.DefaultValue)
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
				switch v1, v2 := dataList[i][strcase.LowerCamelCase(field.Name)].(core.AnyT), dataList[j][strcase.LowerCamelCase(field.Name)].(core.AnyT); {
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
	for _, macro := range r.GetMacros() {
		var childMacroMap = make(map[string]any)
		for _, macroDetail := range macro.DetailList {
			childMacroMap[strcase.UpperSnakeCase(macro.KeyField.Type.Convert(macroDetail.Key))] = macroDetail.Value
		}
		macroMap[strcase.LowerCamelCase(macro.MacroName)] = childMacroMap
	}
	return macroMap
}

func (r *JSONRender) Verify() error {
	return nil
}

func (r *JSONRender) ConfigName() string {
	return strcase.LowerCamelCase(r.GetTableNamePrefix() + r.Name)
}

func (r *JSONRender) Filename() string {
	return strcase.KebabCase(r.GetFilePrefix()+r.Name) + ".json"
}
