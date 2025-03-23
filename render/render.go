package render

import (
	"fmt"
	"xCelFlow/config"
	"xCelFlow/core"
)

type Render struct {
	*core.Table
	SchemaName string
	Schema     config.Schema
}

var renderRegistry = make(map[string]func(render *Render) (core.IRender, error))

func Register(key string, cls func(render *Render) (core.IRender, error)) {
	renderRegistry[key] = cls
}

func NewRender(schemaName string, table *core.Table) (core.IRender, error) {
	if checkSkip(table) {
		fmt.Printf("没有定义字段名跳过生成配置文件：%s\n", table.Filename)
		return nil, nil
	}

	Schema := config.Config.GetSchema(schemaName)

	if cls, ok := renderRegistry[schemaName]; ok {
		cr, err := cls(&Render{Table: table, SchemaName: schemaName, Schema: Schema})
		if err != nil {
			return nil, err
		}
		fmt.Printf("开始导出%s配置：%s\n", schemaName, table.Filename)
		return cr, nil
	}
	panic(fmt.Errorf("配置表：%s 渲染模板：%s 还没有被支持", table.Filename, schemaName))
}

func checkSkip(table *core.Table) bool {
	if table.FieldLen == 0 {
		for _, macroDec := range table.GetMacros() {
			if len(macroDec.DetailList) > 0 {
				return false
			}
		}
		return true
	}
	return false
}

func (r *Render) Verify() error {
	return nil
}
