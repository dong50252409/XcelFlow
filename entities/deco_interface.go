package entities

type IDecorator interface {
	Name() string
}

type ITableDecorator interface {
	IDecorator
	RunTableDecorator(tbl *Table) error
}

type IFieldDecorator interface {
	IDecorator
	RunFieldDecorator(tbl *Table, field *Field) error
}
