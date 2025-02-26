package flags

import "flag"

// typescript 专用参数
var (
	TSClsDirectory string
	TSFbsDirectory string
)

func init() {
	flag.StringVar(&TSClsDirectory, "ts_cls_dir", "", "指定TypeScript的配置类导出目录.")
	flag.StringVar(&TSFbsDirectory, "ts_fbs_dir", "", "指定Typescript的Flatbuffers类导出目录.")
}

func mergeTypeScriptArgs(schemaArgs map[string]interface{}) {
	schemaArgs["ts_cls_directory"] = TSClsDirectory
	schemaArgs["ts_fbs_directory"] = TSFbsDirectory
}
