//go:build console

package main

import (
	"cfg_exporter/config"
	"cfg_exporter/flags"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
			printStack()
		}
	}()

	flag.Parse()
	if flag.NFlag() == 0 || flags.Help {
		flag.Usage()
	} else {
		if flags.TomlPath != "" && flags.SchemaName != "" {
			if _, err := os.Stat(flags.TomlPath); !os.IsNotExist(err) {
				config.NewTomlConfig(flags.TomlPath)
			} else {
				panic(fmt.Errorf("配置文件不存在 %s", flags.TomlPath))
			}
		} else {
			config.NewTomlConfigByFlags()
		}

		var filepathList []string
		_ = filepath.WalkDir(config.Config.GetSource(),
			func(path string, d os.DirEntry, _ error) error {
				if !d.IsDir() {
					filepathList = append(filepathList, path)
				}
				return nil
			})
		runTask(filepathList)
	}
}
