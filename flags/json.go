package flags

import "flag"

var (
	JsonDirectory string
)

func init() {
	flag.StringVar(&JsonDirectory, "json_dir", "", "指定JSON导出目录.")
}

func mergeJSONArgs(schemaArgs map[string]interface{}) {
	schemaArgs["json_directory"] = JsonDirectory
}
