package config

import (
	"fmt"
	"os"
	"xCelFlow/flags"

	"github.com/pelletier/go-toml/v2"
)

var Config TomlConfig

type TomlConfig struct {
	source             string
	fieldCommentRow    int
	fieldTypeRow       int
	fieldDecoratorRow  int
	bodyStartRow       int
	verify             bool
	sqliteDirectory    string
	schemas            map[string]Schema
	fieldNameIndexList []int
}

// NewTomlConfig 载入toml配置
func NewTomlConfig(path string) {
	content, _ := os.ReadFile(path)
	var cfg map[string]interface{}
	if err := toml.Unmarshal(content, &cfg); err != nil {
		panic(fmt.Errorf("解析配置文件失败 %s", err))
	}

	Config = TomlConfig{
		source:            cfg["source"].(string),
		fieldTypeRow:      int(cfg["field_type_row"].(int64)),
		fieldDecoratorRow: int(cfg["field_decorator_row"].(int64)),
		fieldCommentRow:   int(cfg["field_comment_row"].(int64)),
		bodyStartRow:      int(cfg["body_start_row"].(int64)),
		verify:            cfg["verify"].(bool),
		sqliteDirectory:   cfg["sqlite_directory"].(string),
		schemas:           initSchemas(cfg),
	}
	Config.fieldNameIndexList = getFieldNameIndexList()
}

// NewTomlConfigByFlags 根据命令行参数生成配置
func NewTomlConfigByFlags() {
	Config = TomlConfig{
		source:            flags.Source,
		fieldTypeRow:      flags.FieldTypeRow,
		fieldDecoratorRow: flags.FieldDecoratorRow,
		fieldCommentRow:   flags.FieldCommentRow,
		bodyStartRow:      flags.BodyStartRow,
		verify:            flags.Verify,
		schemas:           make(map[string]Schema),
	}
	Config.fieldNameIndexList = getFieldNameIndexList()
	schemaArgs := flags.GetSchemaArgs()
	initSchema(flags.SchemaName, schemaArgs, Config.schemas)
}

func getFieldNameIndexList() []int {
	rowSet := make(map[int]struct{})
	for _, schema := range Config.schemas {
		rowSet[schema.GetFieldNameRow()] = struct{}{}
	}
	rowList := make([]int, 0, len(rowSet))
	for k := range rowSet {
		rowList = append(rowList, k-1)
	}
	return rowList
}

// GetSource 获取表文件路径
func (t TomlConfig) GetSource() string {
	return t.source
}

// GetFieldCommentRow 获取字段注释行
func (t TomlConfig) GetFieldCommentRow() int {
	return t.fieldCommentRow
}

// GetFieldTypeRow 获取字段类型行
func (t TomlConfig) GetFieldTypeRow() int {
	return t.fieldTypeRow
}

// GetFieldDecoratorRow 获取字段装饰器行
func (t TomlConfig) GetFieldDecoratorRow() int {
	return t.fieldDecoratorRow
}

// GetBodyStartRow 获取主体数据开始行
func (t TomlConfig) GetBodyStartRow() int {
	return t.bodyStartRow
}

// GetVerify 获取是否校验生成后的配置
func (t TomlConfig) GetVerify() bool {
	return t.verify
}

func (t TomlConfig) SetVerify(b bool) {
	t.verify = b
}

// GetSqliteDirectory 获取SQLite导出目录
func (t TomlConfig) GetSqliteDirectory() string {
	return t.sqliteDirectory
}

// GetSchema 获取模式字典
func (t TomlConfig) GetSchema(schema string) Schema {
	return t.schemas[schema]
}

// GetSchemas 获取模式字典
func (t TomlConfig) GetSchemas() map[string]Schema {
	return t.schemas
}

// GetFieldTypeIndex 获取字段类型索引
func (t TomlConfig) GetFieldTypeIndex() int {
	return t.fieldTypeRow - 1
}

// GetFieldCommentIndex 获取字段注释索引
func (t TomlConfig) GetFieldCommentIndex() int {
	return t.fieldCommentRow - 1
}

// GetFieldDecoratorIndex 获取字段装饰器索引
func (t TomlConfig) GetFieldDecoratorIndex() int {
	return t.fieldDecoratorRow - 1
}

// GetSchemaFieldNameIndex 获取模式字段名索引
func (t TomlConfig) GetSchemaFieldNameIndex(schemaName string) int {
	return t.GetSchema(schemaName).GetFieldNameRow() - 1
}

// GetFieldNameIndexList 获取字段名索引列集合
func (t TomlConfig) GetFieldNameIndexList() []int {
	return t.fieldNameIndexList
}

// GetBodyStartIndex 获取主体数据开始索引
func (t TomlConfig) GetBodyStartIndex() int {
	return t.bodyStartRow - 1
}
