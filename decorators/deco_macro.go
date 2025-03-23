package decorators

import (
	"fmt"
	"strconv"
	"strings"
	"xCelFlow/core"
	"xCelFlow/util"
)

type DecoMacro struct {
	*core.Macro
}

func init() {
	decoratorRegister("macro", newMacro)
}

func newMacro(tbl *core.Table, field *core.Field, str string) error {
	if param := util.SubParam(str); param != "" {

		if l := strings.Split(param, ","); len(l) == 1 {
			valueFieldName := l[0]
			valueField := tbl.GetFieldByName(valueFieldName)
			if valueField == nil {
				return fmt.Errorf("%s 宏 %s 值字段不存在", field.Name, valueFieldName)
			}
			baseMacro := &core.Macro{MacroName: field.Name, KeyField: field, ValueField: valueField}
			tbl.Decorators = append(tbl.Decorators, &DecoMacro{Macro: baseMacro})
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
			baseMacro := &core.Macro{MacroName: field.Name, KeyField: field, ValueField: valueField, CommentField: commentField}
			tbl.Decorators = append(tbl.Decorators, &DecoMacro{Macro: baseMacro})
			return nil
		}
	}
	return fmt.Errorf("参数格式错误 macro(值字段名[,描述字段名])")
}

func getField(tbl *core.Table, fieldName string) *core.Field {
	if colRow, err := strconv.Atoi(fieldName); err == nil {
		return tbl.GetFieldByColumn(colRow - 1)
	} else {
		return tbl.GetFieldByName(fieldName)
	}
}

func (m *DecoMacro) Name() string {
	return "macro"
}

func (m *DecoMacro) RunTableDecorator(tbl *core.Table) error {
	for _, row := range tbl.DataSetIter() {
		if key := row[m.KeyField.Column]; key != nil && key != "" {
			var comment string
			if m.CommentField != nil {
				comment = util.Quoted(row[m.CommentField.Column].(string))
			}
			m.DetailList = append(m.DetailList, &core.MacroDetail{
				Key:     fmt.Sprintf("%v", key),
				Value:   row[m.ValueField.Column],
				Comment: comment,
			})
		}
	}
	return nil
}

func (m *DecoMacro) GetMacro() *core.Macro {
	return m.Macro
}
