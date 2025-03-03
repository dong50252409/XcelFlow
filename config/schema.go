package config

type baseSchema struct {
	Schema
	FieldNameRow    int    `toml:"field_name_row"`
	FilePrefix      string `toml:"file_prefix"`
	TableNamePrefix string `toml:"table_name_prefix"`
}

type Schema interface {
	GetFieldNameRow() int
	GetFilePrefix() string
	GetTableNamePrefix() string
}

func initSchemas(cfg map[string]interface{}) map[string]Schema {
	schema := make(map[string]Schema, len(cfg["schema"].(map[string]interface{})))
	for k, v := range cfg["schema"].(map[string]interface{}) {
		v1 := v.(map[string]interface{})
		initSchema(k, v1, schema)
	}
	return schema
}

func initSchema(schemaName string, schemaArgs map[string]interface{}, schema map[string]Schema) {
	bs := baseSchema{
		FieldNameRow:    int(schemaArgs["field_name_row"].(int64)),
		FilePrefix:      schemaArgs["file_prefix"].(string),
		TableNamePrefix: schemaArgs["table_name_prefix"].(string),
	}
	switch {
	case schemaName == "erlang":
		schema[schemaName] = initErlang(schemaArgs, bs)
	case schemaName == "flatbuffers":
		schema[schemaName] = initFlatbuffers(schemaArgs, bs)
	case schemaName == "json":
		schema[schemaName] = initJSON(schemaArgs, bs)
	case schemaName == "typescript":
		schema[schemaName] = initTypeScript(schemaArgs, bs)
	}
}

// GetFieldNameRow 获取字段名行
func (b *baseSchema) GetFieldNameRow() int {
	return b.FieldNameRow
}

// GetFilePrefix 获取文件前缀
func (b *baseSchema) GetFilePrefix() string {
	return b.FilePrefix
}

// GetTableNamePrefix 获取表名前缀
func (b *baseSchema) GetTableNamePrefix() string {
	return b.TableNamePrefix
}
