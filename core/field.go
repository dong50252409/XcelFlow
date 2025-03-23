package core

type Field struct {
	// 字段在表中所在列
	Column int
	// 字段名
	Name string
	// 字段类型
	Type IType
	// 字段描述
	Comment string
	// 默认值
	DefaultValue any
	// 默认值字符串
	DefaultValueStr string
	// 装饰器
	Decorators map[string]IFieldDecorator
}

func (f *Field) Convert(v any) string {
	if v != nil {
		return f.Type.Convert(v)
	}
	return f.DefaultValueStr
}
