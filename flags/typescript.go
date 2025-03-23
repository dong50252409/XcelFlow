package flags

import "flag"

var (
	TsDirectory         string
	TsMethodInCamelCase bool
)

func init() {
	flag.StringVar(&TsDirectory, "ts_dir", "", "指定TypeScript导出目录.")
	flag.BoolVar(&TsMethodInCamelCase, "ts_method_in_camel_case", false, "指定TypeScript方法名是否转换为驼峰命名.")
}

func mergeTypeScriptArgs(schemaArgs map[string]interface{}) {
	schemaArgs["ts_directory"] = TsDirectory
	schemaArgs["ts_method_in_camel_case"] = TsMethodInCamelCase
}
