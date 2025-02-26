package render

import (
	"cfg_exporter/config"
	"cfg_exporter/entities"
	"fmt"
	"os"
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

var renderRegistry = make(map[string]func(render *Render) IRender)

func Register(key string, cls func(render *Render) IRender) {
	renderRegistry[key] = cls
}

func NewRender(schemaName string, table *entities.Table) (IRender, error) {
	if checkSkip(table) {
		fmt.Printf("没有定义字段名跳过生成配置文件：%s\n", table.Filename)
		return nil, nil
	}

	if cls, ok := renderRegistry[schemaName]; ok {
		r := &Render{table, schemaName, config.Config.GetSchema(schemaName)}
		return cls(r), nil
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

func (r Render) ExportDir() string {
	return r.Schema.GetDestination()
}

func (r Render) ExecuteBefore() error {
	fmt.Printf("开始导出%s配置：%s\n", r.SchemaName, r.Table.Filename)
	dir := r.ExportDir()
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("导出路径创建失败 %s", err)
	}
	return nil
}
