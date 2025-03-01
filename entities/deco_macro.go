package entities

import (
	"fmt"
	"strconv"
	"strings"
	"xCelFlow/util"
)

type MacroDetail struct {
	Key     any
	Value   any
	Comment string
}

type Macro struct {
	// 宏名，一个表可能包含多个不同的宏数据
	MacroName    string
	KeyField     *Field
	ValueField   *Field
	CommentField *Field
	List         []MacroDetail
}

func init() {
	decoratorRegister("macro", newMacro)
}

func newMacro(tbl *Table, field *Field, str string) error {
	if param := util.SubParam(str); param != "" {

		if l := strings.Split(param, ","); len(l) == 1 {
			valueFieldName := l[0]
			valueField := tbl.GetFieldByName(valueFieldName)
			if valueField == nil {
				return fmt.Errorf("%s 宏 %s 值字段不存在", field.Name, valueFieldName)
			}

			tbl.Decorators = append(tbl.Decorators, &Macro{MacroName: field.Name, KeyField: field, ValueField: valueField})
			return nil
		} else if len(l) == 2 {
			valueFieldName, commentFieldName := l[0], l[1]
			valueField := getField(tbl, valueFieldName)
			if valueField == nil {
				return fmt.Errorf("%s 宏 %s 值字段不存在", field.Name, valueFieldName)
			}
			commentField := getField(tbl, commentFieldName)
			if commentField == nil {
				return fmt.Errorf("%s 宏 %s 描述字段不存在", field.Name, commentFieldName)
			}

			tbl.Decorators = append(tbl.Decorators, &Macro{MacroName: field.Name, KeyField: field, ValueField: valueField, CommentField: commentField})
			return nil
		}
	}
	return fmt.Errorf("参数格式错误 macro(值字段名[,描述字段名])")
}

func getField(tbl *Table, fieldName string) *Field {
	if colRow, err := strconv.Atoi(fieldName); err == nil {
		return tbl.GetFieldByColumn(colRow - 1)
	} else {
		return tbl.GetFieldByName(fieldName)
	}
}

func (m *Macro) Name() string {
	return "macro"
}

func (m *Macro) RunTableDecorator(tbl *Table) error {
	for _, row := range tbl.DataSetIter() {
		if key := row[m.KeyField.Column]; key != "" {
			var comment string
			if m.CommentField != nil {
				comment = util.Quoted(row[m.CommentField.Column].(string))
			}
			m.List = append(m.List, MacroDetail{
				Key:     key,
				Value:   row[m.ValueField.Column],
				Comment: comment,
			})
		}
	}
	return nil
}
