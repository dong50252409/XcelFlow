package config

type ErlangSchema struct {
	baseSchema
	hrlDirectory         string
	erlDirectory         string
	erlTemplateDirectory string
}

// 初始化Erlang配置
func initErlang(schemaArgs map[string]interface{}, bs baseSchema) Schema {
	return &ErlangSchema{
		baseSchema:           bs,
		hrlDirectory:         schemaArgs["hrl_directory"].(string),
		erlDirectory:         schemaArgs["erl_directory"].(string),
		erlTemplateDirectory: schemaArgs["erl_template_directory"].(string),
	}
}

// GetHrlDirectory 获取hrl文件路径
func (e ErlangSchema) GetHrlDirectory() string {
	return e.hrlDirectory
}

// GetErlDirectory 获取erl文件路径
func (e ErlangSchema) GetErlDirectory() string {
	return e.erlDirectory
}

// GetErlTemplates 获取erl模板文件路径
func (e ErlangSchema) GetErlTemplates() string {
	return e.erlTemplateDirectory
}
