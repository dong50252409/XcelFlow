package parser

import (
	"fmt"
	"regexp"
	"strings"
	"xCelFlow/config"
	"xCelFlow/core"
	"xCelFlow/decorators"
	"xCelFlow/types"
	"xCelFlow/util"
)

type FormParser struct {
	*config.TomlConfig
	*core.Table
	// 模式名
	schemaName string
	// 新建类型函数
	NewTypeFunc func(string) (core.IType, error)
	// 检查是否有重复字段名
	FieldSet map[string]struct{}
}

var (
	parserRegistry       = make(map[string]func(p *FormParser) core.IParser)
	splitMultiConsRegexp = regexp.MustCompile(`\r\n|\r|\n`)
)

// RegisterParser 注册解析器
func RegisterParser(name string, cls func(p *FormParser) core.IParser) {
	parserRegistry[name] = cls
}

// NewParser 新建解析器
func NewParser(schemaName string, tbl *core.Table) (core.IParser, error) {
	if parser, ok := parserRegistry[schemaName]; ok {
		p := &FormParser{
			Table:       tbl,
			TomlConfig:  &config.Config,
			schemaName:  schemaName,
			NewTypeFunc: types.NewType,
			FieldSet:    make(map[string]struct{}),
		}
		return parser(p), nil
	}
	return nil, fmt.Errorf("未找到解析器：%s", schemaName)
}

func CloneParser(schemaName string, tbl *core.Table) (core.IParser, error) {
	if parser, ok := parserRegistry[schemaName]; ok {
		p := &FormParser{
			Table:       tbl,
			TomlConfig:  &config.Config,
			schemaName:  schemaName,
			NewTypeFunc: types.NewType,
			FieldSet:    make(map[string]struct{}),
		}

		cloneParser := parser(p)
		fields := p.Fields
		cloneFields := make([]*core.Field, len(tbl.Fields))
		p.Fields = cloneFields

		for column, field := range fields {
			cloneField := &core.Field{
				Column:       column,
				Comment:      field.Comment,
				DefaultValue: field.DefaultValue,
				Decorators:   field.Decorators,
			}
			p.Fields[column] = cloneField

			fieldName, err := cloneParser.ParseFiledNameByColumn(column)
			if err != nil {
				return nil, err
			}
			cloneField.Name = fieldName

			fieldType, err := cloneParser.ParseFieldTypeByColumn(column)
			if err != nil {
				return nil, err
			}
			cloneField.Type = fieldType

			if field.Type != nil {
				DefaultValueStr := field.Type.Convert(field.DefaultValue)
				cloneField.DefaultValueStr = DefaultValueStr
			}
		}
		p.FieldLen = len(p.FieldSet)

		return cloneParser, nil
	}
	return nil, fmt.Errorf("未找到解析器：%s", schemaName)
}

// Parse 解析配表
func (p *FormParser) Parse() error {
	if err := p.ParseHead(); err != nil {
		return err
	}

	if err := p.ParseRow(); err != nil {
		return err
	}

	return nil
}

// ParseHead 解析表头
func (p *FormParser) ParseHead() error {
	records := p.Records[0]
	fields := make([]*core.Field, len(records))
	p.Fields = fields

	for column := range records {
		field := &core.Field{
			Column:     column,
			Decorators: make(map[string]core.IFieldDecorator),
		}
		p.Fields[column] = field

		fieldName, err := p.ParseFiledNameByColumn(column)
		if err != nil {
			return err
		}
		field.Name = fieldName

		fieldType, err := p.ParseFieldTypeByColumn(column)
		if err != nil {
			return err
		}
		field.Type = fieldType

		fieldComment, err := p.ParseFieldCommentByColumn(column)
		if err != nil {
			return err
		}
		field.Comment = fieldComment

		if field.Type != nil {
			field.DefaultValue = field.Type.DefaultValue()
			field.DefaultValueStr = field.Type.Convert(field.DefaultValue)
		}
	}

	p.FieldLen = len(p.FieldSet)

	if err := p.ParseFieldDecorators(); err != nil {
		return err
	}
	return nil
}

// ParseFiledNameByColumn 解析字段名
func (p *FormParser) ParseFiledNameByColumn(column int) (string, error) {
	fieldNameIndex := p.GetSchemaFieldNameIndex(p.schemaName)
	fieldNameRow := p.Records[fieldNameIndex]

	if fieldName := strings.TrimSpace(fieldNameRow[column]); fieldName != "" {
		if _, ok := p.FieldSet[fieldName]; !ok {
			p.FieldSet[fieldName] = struct{}{}
			return fieldName, nil
		}
		return "", fmt.Errorf("单元格：%s、%s\n错误：%s 重复定义", util.ToCell(fieldNameIndex, column), util.ToCell(fieldNameIndex, column), fieldName)
	}
	return "", nil
}

// ParseFieldTypeByColumn 解析字段类型
func (p *FormParser) ParseFieldTypeByColumn(column int) (core.IType, error) {
	fieldTypeIndex := p.GetFieldTypeIndex()
	fieldTypeRow := p.Records[fieldTypeIndex]

	if fieldType := strings.ReplaceAll(fieldTypeRow[column], " ", ""); fieldType != "" {
		t, err := p.NewTypeFunc(fieldType)
		if err != nil {
			return nil, fmt.Errorf("单元格：%s\n错误：%s", util.ToCell(p.GetFieldTypeIndex(), column), err)
		}
		return t, nil
	} else if p.Fields[column].Name != "" {
		return nil, fmt.Errorf("单元格：%s\n错误：类型不能为空", util.ToCell(p.GetFieldTypeIndex(), column))
	}
	return nil, nil
}

func (p *FormParser) ParseFieldCommentByColumn(column int) (string, error) {
	fieldCommentIndex := p.GetFieldCommentIndex()
	fieldCommentRow := p.Records[fieldCommentIndex]
	if len(fieldCommentRow) > column {
		return util.Quoted(strings.TrimSpace(fieldCommentRow[column])), nil
	}
	return "", fmt.Errorf("单元格：%s\n错误：注释不能为空", util.ToCell(p.GetFieldCommentIndex(), column))
}

// ParseFieldDecorators 解析字段装饰器
func (p *FormParser) ParseFieldDecorators() error {
	// 第一列默认为主键列
	fieldDecoratorRow := p.Records[p.GetFieldDecoratorIndex()]
	if !strings.Contains(fieldDecoratorRow[0], "p_key") {
		fieldDecoratorRow[0] = "p_key\n" + fieldDecoratorRow[0]
	}

	for column, field := range p.Fields {
		if err := p.ParseFieldDecoratorsByColumn(column, field); err != nil {
			return err
		}
	}
	return nil
}

// ParseFieldDecoratorsByColumn 解析字段装饰器
func (p *FormParser) ParseFieldDecoratorsByColumn(column int, field *core.Field) error {
	fieldDecoratorIndex := p.GetFieldDecoratorIndex()
	fieldDecoratorRow := p.Records[fieldDecoratorIndex]

	if len(fieldDecoratorRow) > column {
		if fieldDec := strings.TrimSpace(fieldDecoratorRow[column]); fieldDec != "" {
			parts := splitMultiConsRegexp.Split(fieldDec, -1)
			for _, part := range parts {
				if err := decorators.NewDecorator(p.Table, field, part); err != nil {
					return fmt.Errorf("单元格：%s\n装饰器：%s\n错误：%s", field.Name, part, err)
				}
			}
		}
	}
	return nil
}

// ParseRow 解析行
func (p *FormParser) ParseRow() error {
	records := p.Records[p.GetBodyStartIndex():]

	// 初始化数据，通过切片重组二维结构
	totalRows, rowLength := len(records), len(p.Fields)
	flat := make([]interface{}, totalRows*rowLength)
	p.DataSet = make([][]interface{}, totalRows)
	for i := range p.DataSet {
		p.DataSet[i] = flat[i*rowLength : (i+1)*rowLength : (i+1)*rowLength]
	}

	for rowIndex, row := range records {
		for _, field := range p.Fields {
			if len(row) > field.Column {
				if cell := strings.TrimSpace(row[field.Column]); cell != "" && field.Type != nil {
					value, err := field.Type.ParseString(cell)
					if err != nil {
						panic(fmt.Sprintf("单元格：%s\n错误：%s", util.ToCell(rowIndex+p.GetBodyStartIndex(), field.Column), err))
					}
					p.DataSet[rowIndex][field.Column] = value
				}
			}
		}
		if p.DataSet[rowIndex][0] != nil {
			p.DataSetLen += 1
		}
	}
	return nil
}

// RunDecorators 运行装饰器
func (p *FormParser) RunDecorators() error {
	for _, d := range p.Decorators {
		if err := d.RunTableDecorator(p.Table); err != nil {
			return fmt.Errorf("装饰器：%s\n%s", d.Name(), err)
		}
	}

	for _, field := range p.Fields {
		for k, d := range field.Decorators {
			if err := d.RunFieldDecorator(p.Table, field); err != nil {
				return fmt.Errorf("装饰器：%s\n%s", k, err)
			}
		}
	}
	return nil
}

func (p *FormParser) GetTable() *core.Table {
	return p.Table
}
