package flags

import "flag"

// flatbuffers 专用参数
var (
	Flatc        string
	Namespace    string
	FbsDirectory string
	BinDirectory string
)

func init() {
	flag.StringVar(&Flatc, "flatc_path", "", "指定flatc路径.")
	flag.StringVar(&Namespace, "namespace", "", "指定flatbuffers的命名空间.")
	flag.StringVar(&FbsDirectory, "fbs_dir", "", "指定Flatbuffers描述文件导出目录.")
	flag.StringVar(&BinDirectory, "bin_dir", "", "指定Flatbuffers序列化数据导出目录.")
}

func mergeFlatbuffersArgs(schemaArgs map[string]interface{}) {
	schemaArgs["flatc"] = Flatc
	schemaArgs["namespace"] = Namespace
	schemaArgs["fbs_directory"] = FbsDirectory
	schemaArgs["bin_directory"] = BinDirectory
}
