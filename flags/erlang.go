package flags

import "flag"

// erlang 专用参数
var (
	HrlDirectory         string
	ErlDirectory         string
	ErlTemplateDirectory string
)

func init() {
	flag.StringVar(&HrlDirectory, "hrl_dir", "", "指定Erlang头文件导出目录.")
	flag.StringVar(&ErlDirectory, "erl_dir", "", "指定Erlang配置文件导出目录.")
	flag.StringVar(&ErlTemplateDirectory, "erl_template_dir", "", "指定Erlang模板文件路径.")
}

func mergeErlangArgs(schemaArgs map[string]interface{}) {
	// erlang
	schemaArgs["hrl_directory"] = HrlDirectory
	schemaArgs["erl_directory"] = ErlDirectory
	schemaArgs["erl_template_directory"] = ErlTemplateDirectory
}
