package config

type TypeScriptSchema struct {
	baseSchema
	TsDirectory string
}

// 初始化TypeScript配置
func initTypeScript(schemaArgs map[string]interface{}, bs baseSchema) Schema {
	return &TypeScriptSchema{baseSchema: bs, TsDirectory: schemaArgs["ts_directory"].(string)}
}

func (s *TypeScriptSchema) GetTsDirectory() string {
	return s.TsDirectory
}
