package render

import (
	"fmt"
	"xCelFlow/config"
	"xCelFlow/entities"
)

type IRender interface {
	Execute() error
	Verify() error
}

type Render struct {
	*entities.Table
	SchemaName string
	Schema     config.Schema
}

var renderRegistry = make(map[string]func(render *Render) (IRender, error))

func Register(key string, cls func(render *Render) (IRender, error)) {
	renderRegistry[key] = cls
}

func NewRender(schemaName string, table *entities.Table) (IRender, error) {
	if checkSkip(table) {
		fmt.Printf("没有定义字段名跳过生成配置文件：%s\n", table.Filename)
		return nil, nil
	}

	Schema := config.Config.GetSchema(schemaName)

	if cls, ok := renderRegistry[schemaName]; ok {
		cr, err := cls(&Render{Table: table, SchemaName: schemaName, Schema: Schema})
		fmt.Printf("开始导出%s配置：%s\n", schemaName, table.Filename)
		return cr, err
	}
	return nil, fmt.Errorf("配置表：%s 渲染模板：%s 还没有被支持", table.Filename, schemaName)
}

func checkSkip(table *entities.Table) bool {
	if table.FieldLen == 0 {
		for _, macroDec := range table.GetMacroDecorators() {
			if len(macroDec.List) > 0 {
				return false
			}
		}
		return true
	}
	return false
}
