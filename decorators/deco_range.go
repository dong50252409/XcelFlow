package decorators

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"xCelFlow/config"
	"xCelFlow/core"
	"xCelFlow/util"
)

// DecoRange 范围
type DecoRange struct {
	minValue any
	maxValue any
}

func init() {
	decoratorRegister("range", newRange)
}

func newRange(_ *core.Table, field *core.Field, str string) error {
	if field.Type != nil {
		if param := util.SubParam(str); param != "" {
			if l := strings.Split(param, ","); len(l) == 2 {
				switch field.Type.Kind() {
				case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					v1, err1 := strconv.ParseInt(l[0], 10, 64)
					v2, err2 := strconv.ParseInt(l[1], 10, 64)
					if err1 == nil && err2 == nil && v1 <= v2 {
						field.Decorators["range"] = &DecoRange{minValue: v1, maxValue: v2}
						return nil
					}
				case reflect.Float32, reflect.Float64:
					v1, err1 := strconv.ParseFloat(l[0], 64)
					v2, err2 := strconv.ParseFloat(l[1], 64)
					if err1 == nil && err2 == nil && v1 <= v2 {
						field.Decorators["range"] = &DecoRange{minValue: v1, maxValue: v2}
						return nil
					}
				default:
					return fmt.Errorf("类型无法使用此装饰器")
				}
			}
		}
		return fmt.Errorf("参数格式错误 range(最小值,最大值)")
	}
	return nil
}

func (r *DecoRange) RunFieldDecorator(tbl *core.Table, field *core.Field) error {
	for corIndex, row := range tbl.DataSetIter() {
		if v := row[field.Column]; v != nil {
			if err := r.Equal(corIndex, row[field.Column], field); err != nil {
				return err
			}
		}
	}
	return nil
}

func (*DecoRange) Name() string {
	return "range"
}

func (r *DecoRange) Equal(rowIndex int, v any, field *core.Field) error {
	switch field.Type.Kind() {
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if !(r.minValue.(int64) <= v.(int64) && v.(int64) <= r.maxValue.(int64)) {
			return fmt.Errorf("单元格：%s 数值必须在%d到%d之间", util.ToCell(rowIndex+config.Config.GetBodyStartRow(), field.Column), r.minValue, r.maxValue)
		}
		return nil
	case reflect.Float32, reflect.Float64:
		if !(r.minValue.(float64) <= v.(float64) && v.(float64) <= r.maxValue.(float64)) {
			return fmt.Errorf("单元格：%s 数值必须在%d到%d之间", util.ToCell(rowIndex+config.Config.GetBodyStartRow(), field.Column), r.minValue, r.maxValue)
		}
		return nil
	default:
		return fmt.Errorf("类型无法使用此装饰器")
	}
}
