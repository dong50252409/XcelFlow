package flags

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

const VERSION = "0.1.0"

var (
	Help              bool
	Source            string
	FieldTypeRow      int
	FieldDecoratorRow int
	FieldCommentRow   int
	BodyStartRow      int
	Verify            bool
	SchemaName        string
	TomlPath          string
)

func init() {
	flag.BoolVar(&Help, "h", false, "显示帮助信息.")
	flag.StringVar(&Source, "src", "", "指定配置表文件路径.")
	flag.IntVar(&FieldTypeRow, "type_row", 3, "指定字段类型所在行.")
	flag.IntVar(&FieldDecoratorRow, "dec_row", 4, "指定字段装饰器所在行.")
	flag.IntVar(&FieldCommentRow, "desc_row", 1, "指定字段注释所在行.")
	flag.IntVar(&BodyStartRow, "start_row", 6, "指定表体开始行.")
	flag.BoolVar(&Verify, "verify", false, "是否对生成的文件进行校验.")
	flag.StringVar(&TomlPath, "config_path", "config.toml", "指定config.toml文件路径.")
	flag.StringVar(&SchemaName, "schema_name", "", "指定config.toml中的区域进行导出.")

	flag.Usage = usage
}

func usage() {
	supportedSchema := []string{
		"erlang",
		"flatbuffers",
		"json",
		"typescript",
	}

	if _, err := fmt.Fprintf(os.Stdout, `cfg_exporter version: %s
Usage: cfg_exporter -s %s

Options:
`, VERSION, strings.Join(supportedSchema, "|")); err != nil {
		panic(err)
	}
	flag.PrintDefaults()
}
