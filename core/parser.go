package core

type IParser interface {
	// Parse 解析
	Parse() error
	// ParseHead 解析表头
	ParseHead() error
	// ParseFiledNameByColumn 解析字段名
	ParseFiledNameByColumn(column int) (string, error)
	// ParseFieldTypeByColumn 解析字段类型
	ParseFieldTypeByColumn(column int) (IType, error)
	// ParseFieldCommentByColumn 解析字段描述
	ParseFieldCommentByColumn(column int) (string, error)
	// ParseRow 解析一行
	ParseRow() error
	// RunDecorators 运行字段属性
	RunDecorators() error
	// GetTable 获取表
	GetTable() *Table
}
