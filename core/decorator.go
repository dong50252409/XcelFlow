package core

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

type IPrimaryKey interface {
	GetFields() []*Field
}

type Macro struct {
	MacroName    string
	KeyField     *Field
	ValueField   *Field
	CommentField *Field
	DetailList   []*MacroDetail
}

type MacroDetail struct {
	Key     string
	Value   any
	Comment string
}

// IMacro 宏
type IMacro interface {
	GetMacro() *Macro
}
