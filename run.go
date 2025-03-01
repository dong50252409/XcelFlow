package main

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
	"xCelFlow/config"
	"xCelFlow/entities"
	_ "xCelFlow/implements/erlang"
	_ "xCelFlow/implements/flatbuffers"
	_ "xCelFlow/implements/json"
	_ "xCelFlow/implements/typescript"
	"xCelFlow/parser"
	"xCelFlow/reader"
	"xCelFlow/render"
)

func run(path string, schemaName string) error {
	records, err := readFile(path)
	if err != nil {
		if errors.Is(err, reader.ErrorTableTempFile) {
			return nil
		}
		return err
	}

	t, err := parserTable(path, schemaName, records)
	if err != nil {
		return err
	}

	err = renderTable(schemaName, t)
	if err != nil {
		return err
	}
	return nil
}

func readFile(path string) ([][]string, error) {
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

func parserTable(path string, schemaName string, records [][]string) (*entities.Table, error) {
	p, err := parser.NewParser(path, schemaName, records)
	if err != nil {
		return nil, err
	}

	if err = p.ParseFieldName(); err != nil {
		return nil, err
	}

	if err = p.ParseFieldType(); err != nil {
		return nil, err
	}

	if err = p.ParseFieldDecorators(); err != nil {
		return nil, err
	}

	p.ParseFieldComment()

	if err = p.ParseRow(); err != nil {
		return nil, err
	}

	p.ParseFieldDefaultValue()

	if err = p.RunDecorators(); err != nil {
		return nil, err
	}

	return p.GetTable(), nil
}

func renderTable(schemaName string, t *entities.Table) error {
	if r, err := render.NewRender(schemaName, t); err != nil {
		return err
	} else if r != nil {
		if err = r.Execute(); err != nil {
			return err
		}
		if config.Config.GetVerify() {
			if err = r.Verify(); err != nil {
				return err
			}
		}
	}
	fmt.Printf("导出配置完成：%s\n", t.Filename)
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
