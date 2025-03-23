package config

type TypeScriptSchema struct {
	baseSchema
	TsDirectory         string
	TsMethodInCamelCase bool
}

// 初始化TypeScript配置
func initTypeScript(schemaArgs map[string]interface{}, bs baseSchema) Schema {
	return &TypeScriptSchema{
		baseSchema:          bs,
		TsDirectory:         schemaArgs["ts_directory"].(string),
		TsMethodInCamelCase: schemaArgs["ts_method_in_camel_case"].(bool),
	}
}

func (s *TypeScriptSchema) GetTsDirectory() string {
	return s.TsDirectory
}

func (s *TypeScriptSchema) GetTsMethodInCamelCase() bool {
	return s.TsMethodInCamelCase
}
