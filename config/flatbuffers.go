package config

type FlatbuffersSchema struct {
	baseSchema
	flatc        string `toml:"flatc"`
	namespace    string `toml:"namespace"`
	fbsDirectory string
	binDirectory string
}

// 初始化Flatbuffers配置
func initFlatbuffers(schemaArgs map[string]interface{}, bs baseSchema) Schema {
	return &FlatbuffersSchema{
		baseSchema:   bs,
		flatc:        schemaArgs["flatc"].(string),
		namespace:    schemaArgs["namespace"].(string),
		fbsDirectory: schemaArgs["fbs_directory"].(string),
		binDirectory: schemaArgs["bin_directory"].(string),
	}
}

// GetFlatc 获取flatc路径
func (f FlatbuffersSchema) GetFlatc() string {
	return f.flatc
}

// GetNamespace 获取flatbuffers命名空间
func (f FlatbuffersSchema) GetNamespace() string {
	return f.namespace
}

// GetFbsDirectory 获取描述文件目录
func (f FlatbuffersSchema) GetFbsDirectory() string {
	return f.fbsDirectory
}

// GetBinDirectory 获取序列化文件目录
func (f FlatbuffersSchema) GetBinDirectory() string {
	return f.binDirectory
}
