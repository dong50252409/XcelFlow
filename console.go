//go:build console

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"xCelFlow/config"
	"xCelFlow/flags"
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
		return
	}
	if flags.TomlPath != "" && flags.SchemaName != "" {
		if _, err := os.Stat(flags.TomlPath); !os.IsNotExist(err) {
			config.NewTomlConfig(flags.TomlPath)
		} else {
			fmt.Printf("配置文件不存在 %s\n", flags.TomlPath)
			return
		}
	} else {
		config.NewTomlConfigByFlags()
	}

	if _, err := os.Stat(config.Config.GetSource()); !os.IsNotExist(err) {
		var filepathList []string
		_ = filepath.WalkDir(config.Config.GetSource(),
			func(path string, d os.DirEntry, _ error) error {
				if !d.IsDir() {
					filepathList = append(filepathList, path)
				}
				return nil
			})
		runTask(filepathList)
	} else {
		fmt.Printf("配置表目录文件不存在 %s\n", config.Config.GetSource())
	}
}
