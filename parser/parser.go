package parser

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"xCelFlow/config"
	"xCelFlow/entities"
	"xCelFlow/util"
)

type IParser interface {
	ParseFieldName() error
	ParseFieldType() error
	ParseFieldDecorators() error
	ParseFieldComment()
	ParseFieldDefaultValue()
	ParseRow() error
	RunDecorators() error
	GetTable() *entities.Table
}

type Parser struct {
	*entities.Table
	*config.TomlConfig
	schemaName  string
	NewTypeFunc func(string, *entities.Field) (entities.ITypeSystem, error)
}

var (
	parserRegistry       = make(map[string]func(p *Parser) IParser)
	splitMultiConsRegexp = regexp.MustCompile(`\r\n|\r|\n`)
)

// RegisterParser 注册解析器
func RegisterParser(name string, cls func(p *Parser) IParser) {
	parserRegistry[name] = cls
}

// NewParser 新建解析器
func NewParser(path string, schemaName string, records [][]string) (IParser, error) {
	if parser, ok := parserRegistry[schemaName]; ok {
		filename := filepath.Base(path)
		name, err := util.SubTableName(filename)
		if err != nil {
			return nil, err
		}
		p := &Parser{
			Table: &entities.Table{
				Path:       filepath.Dir(path),
				Filename:   filename,
				Name:       name,
				Decorators: make([]entities.ITableDecorator, 0),
				Records:    records,
			},
			TomlConfig:  &config.Config,
			schemaName:  schemaName,
			NewTypeFunc: entities.NewType,
		}

		return parser(p), nil
	}
	return nil, fmt.Errorf("未找到解析器：%s", schemaName)
}

func CloneParser(schemaName string, tbl *entities.Table) (IParser, error) {
	if parser, ok := parserRegistry[schemaName]; ok {
		p := &Parser{
			Table: &entities.Table{
				Path:       tbl.Path,
				Filename:   tbl.Filename,
				Name:       tbl.Name,
				Decorators: tbl.Decorators,
				DataSetLen: tbl.DataSetLen,
				DataSet:    tbl.DataSet,
				Records:    tbl.Records,
			},
			TomlConfig:  &config.Config,
			schemaName:  schemaName,
			NewTypeFunc: entities.NewType,
		}

		iParser := parser(p)
		if err := iParser.ParseFieldName(); err != nil {
			return nil, err
		}
		if err := iParser.ParseFieldType(); err != nil {
			return nil, err
		}
		return iParser, nil
	}
	return nil, fmt.Errorf("未找到解析器：%s", schemaName)
}

// Parse 解析配表
func (p *Parser) Parse() error {
	err := p.ParseFieldName()
	if err != nil {
		return err
	}

	err = p.ParseFieldType()
	if err != nil {
		return err
	}

	err = p.ParseFieldDecorators()
	if err != nil {
		return err
	}

	p.ParseFieldComment()

	err = p.ParseRow()
	if err != nil {
		return err
	}

	p.ParseFieldDefaultValue()

	return nil
}

// ParseFieldName 解析字段名
func (p *Parser) ParseFieldName() error {
	fieldNameIndex := p.GetSchemaFieldNameIndex(p.schemaName)
	fieldNameRow := p.Records[fieldNameIndex]
	p.Fields = make([]*entities.Field, len(fieldNameRow))

	// 检查是否有重复字段名
	fieldSet := make(map[string]struct{})
	for column, fieldName := range fieldNameRow {
		fieldName = strings.TrimSpace(fieldName)
		if _, ok := fieldSet[fieldName]; !ok {
			if fieldName != "" {
				fieldSet[fieldName] = struct{}{}
			}
			field := &entities.Field{Column: column, Name: fieldName, Decorators: make(map[string]entities.IFieldDecorator)}
			p.Fields[column] = field
		} else {
			return fmt.Errorf("单元格：%s、%s\n错误：%s 重复定义", util.ToCell(fieldNameIndex, column), util.ToCell(fieldNameIndex, column), fieldName)
		}
	}
	p.FieldLen = len(fieldSet)
	return nil
}

// ParseFieldType 解析字段类型
func (p *Parser) ParseFieldType() error {
	fieldTypeRow := p.Records[p.GetFieldTypeIndex()]

	for _, field := range p.Fields {
		if fieldType := strings.ReplaceAll(fieldTypeRow[field.Column], " ", ""); fieldType != "" {
			t, err := p.NewTypeFunc(fieldType, field)
			if err != nil {
				return fmt.Errorf("单元格：%s\n错误：%s", util.ToCell(p.GetFieldTypeIndex(), field.Column), err)
			}
			field.Type = t
		} else {
			// 检查是否有字段存在但未填写类型的情况
			for _, rowIndex := range p.GetFieldNameIndexList() {
				if p.Records[rowIndex][field.Column] != "" {
					return fmt.Errorf("单元格：%s\n错误：类型不能为空", util.ToCell(p.GetFieldTypeIndex(), field.Column))
				}
			}
		}
	}
	return nil
}

// ParseFieldDecorators 解析装饰器信息
func (p *Parser) ParseFieldDecorators() error {
	fieldDecoratorRow := p.Records[p.GetFieldDecoratorIndex()]
	if fieldDecoratorRow == nil {
		return nil
	}

	// 第一列默认为主键列
	if !strings.Contains(fieldDecoratorRow[0], "p_key") {
		fieldDecoratorRow[0] = "p_key\n" + fieldDecoratorRow[0]
	}

	for column, fieldDec := range fieldDecoratorRow {
		if fieldDec = strings.TrimSpace(fieldDec); fieldDec != "" {
			var field *entities.Field
			if field = p.GetFieldByColumn(column); field == nil {
				field = &entities.Field{Column: column, Decorators: make(map[string]entities.IFieldDecorator)}
			}

			parts := splitMultiConsRegexp.Split(fieldDec, -1)
			for _, part := range parts {
				if err := entities.NewDecorator(p.Table, field, part); err != nil {
					return fmt.Errorf("单元格：%s\n错误：%s 装饰器 %s", field.Name, part, err)
				}
			}
		}
	}
	return nil
}

// ParseFieldComment 解析字段注释
func (p *Parser) ParseFieldComment() {
	fileCommentRow := p.Records[p.GetFieldCommentIndex()]
	for _, field := range p.Fields {
		if len(fileCommentRow) > field.Column {
			field.Comment = util.Quoted(strings.TrimSpace(fileCommentRow[field.Column]))
		}
	}
}

// ParseFieldDefaultValue 解析字段默认值
func (p *Parser) ParseFieldDefaultValue() {
	for _, field := range p.Fields {
		if field.Type != nil {
			field.DefaultValue = field.Type.DefaultValueStr()
		}
	}
}

// ParseRow 解析行
func (p *Parser) ParseRow() error {
	records := p.Records[p.GetBodyStartIndex():]

	// 初始化数据
	totalRows, rowLength := len(records), len(p.Fields)
	flat := make([]interface{}, totalRows*rowLength)
	p.DataSet = make([][]interface{}, totalRows)
	// 通过切片重组二维结构
	for i := range p.DataSet {
		p.DataSet[i] = flat[i*rowLength : (i+1)*rowLength : (i+1)*rowLength]
	}

	for rowIndex, row := range records {
		for _, field := range p.Fields {
			if len(row) > field.Column {
				if cell := strings.TrimSpace(row[field.Column]); cell != "" && field.Type != nil {
					value, err := field.Type.ParseString(cell)
					if err != nil {
						return fmt.Errorf("单元格：%s\n错误：%s", util.ToCell(rowIndex+p.GetBodyStartIndex(), field.Column), err)
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
func (p *Parser) RunDecorators() error {
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

// GetTable 获取配表
func (p *Parser) GetTable() *entities.Table {
	return p.Table
}
