package flags

import "flag"

// schema 通用参数
var (
	FieldNameRow    int
	FilePrefix      string
	TableNamePrefix string
	schemaArgs      map[string]interface{}
)

func init() {
	flag.IntVar(&FieldNameRow, "name_row", 0, "指定字段名所在行.")
	flag.StringVar(&FilePrefix, "file_prefix", "", "指定导出文件的前缀.")
	flag.StringVar(&TableNamePrefix, "table_name_prefix", "", "指定导出表的前缀.")
}

func GetSchemaArgs() map[string]interface{} {
	if schemaArgs != nil {
		schemaArgs = map[string]interface{}{
			"field_name_row":    FieldNameRow,
			"file_prefix":       FilePrefix,
			"table_name_prefix": TableNamePrefix,
		}
		mergeErlangArgs(schemaArgs)
		mergeFlatbuffersArgs(schemaArgs)
		mergeJSONArgs(schemaArgs)
		mergeTypeScriptArgs(schemaArgs)
	}

	return schemaArgs
}
