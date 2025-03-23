package core

import (
	"fmt"
	"iter"
	"path/filepath"
	"xCelFlow/util"
)

type Table struct {
	// 路径
	Path string
	// 文件名
	Filename string
	// 表名
	Name string
	// 字段
	Fields []*Field
	// Field的真实长度
	FieldLen int
	// 装饰器
	Decorators []ITableDecorator
	// DataSet真实长度
	DataSetLen int
	// 主体数据
	DataSet [][]any
	// 原始数据
	Records [][]string
}

// NewTable 新建表
func NewTable(path string, records [][]string) *Table {
	filename := filepath.Base(path)
	if name := util.SubTableName(filename); name == "" {
		panic(fmt.Sprintf("文件名：%s 格式错误 配表描述(表名).ext", filename))
	} else {
		return &Table{
			Path:       filepath.Dir(path),
			Filename:   filename,
			Name:       name,
			Decorators: make([]ITableDecorator, 0),
			Records:    records,
		}
	}
}

// FieldRowIter 字段迭代器，仅返回有效字段
func (tbl *Table) FieldRowIter() iter.Seq2[int, *Field] {
	return func(yield func(int, *Field) bool) {
		index := 0
		for _, field := range tbl.Fields {
			if field.Name != "" {
				if !yield(index, field) {
					return
				}
				index++
			}
		}
	}
}

// GetFieldByName 获取字段根据字段名
func (tbl *Table) GetFieldByName(fieldName string) *Field {
	for _, field := range tbl.FieldRowIter() {
		if field.Name == fieldName {
			return field
		}
	}
	return nil
}

// GetFieldByColumn 获取字段根据表列数
func (tbl *Table) GetFieldByColumn(column int) *Field {
	if len(tbl.Fields) > column {
		return tbl.Fields[column]
	}
	return nil
}

// GetPrimaryKeyFields 获取主键字段列表
func (tbl *Table) GetPrimaryKeyFields() []*Field {
	for _, d := range tbl.Decorators {
		d1, ok := d.(IPrimaryKey)
		if ok {
			return d1.GetFields()
		}
	}
	return []*Field{}
}

// GetPrimaryKeyFieldNames 获取主键字段名列表
func (tbl *Table) GetPrimaryKeyFieldNames() []string {
	fields := tbl.GetPrimaryKeyFields()
	names := make([]string, len(fields))
	for i, field := range fields {
		names[i] = field.Name
	}
	return names
}

// GetPrimaryKeyValues 获取主键值列表
func (tbl *Table) GetPrimaryKeyValues() [][]any {
	fields := tbl.GetPrimaryKeyFields()
	keyLen := len(fields)
	var list = make([][]any, tbl.DataSetLen)
	for rowIndex, dataRow := range tbl.DataSetIter() {
		items := make([]any, keyLen)
		for keyIndex, field := range fields {
			items[keyIndex] = dataRow[field.Column]
		}
		list[rowIndex] = items
	}
	return list
}

// GetPrimaryKeyValuesByString 获取主键值列表,并将值转为字符串
func (tbl *Table) GetPrimaryKeyValuesByString() [][]string {
	fields := tbl.GetPrimaryKeyFields()
	keyLen := len(fields)
	var list = make([][]string, tbl.DataSetLen)
	for rowIndex, dataRow := range tbl.DataSetIter() {
		items := make([]string, keyLen)
		for keyIndex, field := range fields {
			v := dataRow[field.Column]
			items[keyIndex] = field.Type.Convert(v)
		}
		list[rowIndex] = items
	}
	return list
}

// GetMacros 获取宏装饰器集合列表
func (tbl *Table) GetMacros() []*Macro {
	var macroList []*Macro
	for _, d := range tbl.Decorators {
		if d1, ok := d.(IMacro); ok {
			macroList = append(macroList, d1.GetMacro())
		}
	}
	return macroList
}

// DataSetIter 数据集迭代器，仅返回有效行
func (tbl *Table) DataSetIter() iter.Seq2[int, []any] {
	return func(yield func(int, []any) bool) {
		index := 0
		for _, row := range tbl.DataSet {
			if row[0] != nil {
				// 通过 yield 返回索引和行数据
				if !yield(index, row) {
					return // 如果 yield 返回 false，终止迭代
				}
				index++
			}
		}
	}
}
