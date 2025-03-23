package main

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"xCelFlow/config"
	"xCelFlow/core"
	"xCelFlow/parser"
	"xCelFlow/reader"
	"xCelFlow/render"
	"xCelFlow/util"

	_ "xCelFlow/implements/erlang"

	_ "xCelFlow/implements/flatbuffers"

	_ "xCelFlow/implements/json"

	_ "xCelFlow/implements/typescript"
)

func Run(path string, schemaName string) error {
	filename := filepath.Base(path)
	if tableName := util.SubTableName(filename); tableName == "" {
		panic(fmt.Sprintf("文件名：%s 文件名格式错误 格式：配表描述(表名)", filename))
	}

	// 跳过临时文件
	if strings.HasPrefix(filepath.Base(path), "~$") {
		return nil
	}

	records, err := ReadFile(path)
	if err != nil {
		return err
	}
	tbl := core.NewTable(path, records)
	p, err := Parser(schemaName, tbl)
	if err != nil {
		return err
	}

	if err = RenderTable(schemaName, p.GetTable()); err != nil {
		return err
	}
	return nil
}

func ReadFile(path string) ([][]string, error) {
	fmt.Printf("读取配置文件：%s\n", path)
	r, err := reader.NewReader(path)
	if err != nil {
		return nil, err
	}

	records, err := r.Read()
	if err != nil {
		return nil, err
	}

	return records, nil
}

func Parser(schemaName string, tbl *core.Table) (core.IParser, error) {
	p, err := parser.NewParser(schemaName, tbl)
	if err != nil {
		return nil, err
	}

	if err = p.Parse(); err != nil {
		return nil, err
	}

	return p, nil
}

func RenderTable(schemaName string, tbl *core.Table) error {
	r, err := render.NewRender(schemaName, tbl)
	if err != nil {
		return err
	}

	if r == nil {
		return nil
	}

	if err = r.Execute(); err != nil {
		return err
	}

	if config.Config.GetVerify() {
		if err = r.Verify(); err != nil {
			return err
		}
	}
	return nil
}

func printStack() {
	// 获取完整的堆栈信息
	buf := make([]byte, 2048)
	n := runtime.Stack(buf, false)
	stack := string(buf[:n])

	// 按行分割堆栈信息
	lines := strings.Split(stack, "\n")

	// 跳过前 skip*2 行（每个调用占两行：函数名和调用位置）
	if len(lines) > 4 {
		lines = lines[4:]
	}

	// 重新组合并打印堆栈信息
	fmt.Printf("堆栈信息:\n%s\n", strings.Join(lines, "\n"))
}
