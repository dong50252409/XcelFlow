package config

type TypeScriptSchema struct {
	baseSchema
	tsClsDirectory string
	tsFbsDirectory string
}

// 初始化TypeScript配置
func initTypeScript(schemaArgs map[string]interface{}, bs baseSchema) Schema {
	return &TypeScriptSchema{
		baseSchema:     bs,
		tsClsDirectory: schemaArgs["ts_cls_directory"].(string),
		tsFbsDirectory: schemaArgs["ts_fbs_directory"].(string),
	}
}

// GetTSClsDirectory 获取配置类文件路径
func (ts TypeScriptSchema) GetTSClsDirectory() string {
	return ts.tsClsDirectory
}

// GetTSFbsDirectory 获取fbs类文件路径
func (ts TypeScriptSchema) GetTSFbsDirectory() string {
	return ts.tsFbsDirectory
}
