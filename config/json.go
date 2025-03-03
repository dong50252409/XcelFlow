package config

type JSONSchema struct {
	baseSchema
	JsonDirectory string
}

// 初始化JSON配置
func initJSON(schemaArgs map[string]interface{}, bs baseSchema) Schema {
	return &JSONSchema{baseSchema: bs, JsonDirectory: schemaArgs["json_directory"].(string)}
}

func (s *JSONSchema) GetJsonDirectory() string {
	return s.JsonDirectory
}
