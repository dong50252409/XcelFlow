package entities

type Field struct {
	// 字段在表中所在列
	Column int
	// 字段名
	Name string
	// 字段类型
	Type ITypeSystem
	// 字段描述
	Comment string
	// 装饰器
	Decorators map[string]IFieldDecorator
	// 是否为主键
	IsPrimaryKey bool
	// 默认值
	DefaultValue string
}

// Convert 将值类型转为目标语言的值字符串
func (f *Field) Convert(v any) string {
	if v != nil {
		return f.Type.Convert(v)
	} else {
		return f.DefaultValue
	}
}
