package reader

import (
	"cfg_exporter/config"
	"path/filepath"
	"strings"
)

type IReader interface {
	Read() ([][]string, error)
}

type Reader struct {
	IReader
	Path                string
	FieldNameIndexList  []int
	FieldTypeIndex      int
	FieldDecoratorIndex int
	FieldCommentIndex   int
	BodyStartIndex      int
}

var registry = make(map[string]func(reader *Reader) IReader)

// Register 注册文件读取器
func Register(key string, cls func(reader *Reader) IReader) {
	registry[key] = cls
}

func NewReader(path string) (IReader, error) {
	if strings.HasPrefix(filepath.Base(path), "~$") {
		return nil, errorTableTempFile(path)
	}

	ext := filepath.Ext(path)[1:]
	cls, ok := registry[ext]
	if !ok {
		return nil, errorTableNotSupported(path)
	}

	r := cls(&Reader{
		Path:                path,
		FieldNameIndexList:  config.GetFieldNameIndexList(),
		FieldTypeIndex:      config.GetFieldTypeIndex(),
		FieldDecoratorIndex: config.GetFieldDecoratorIndex(),
		FieldCommentIndex:   config.GetFieldCommentIndex(),
		BodyStartIndex:      config.GetBodyStartIndex()})

	return r, nil
}
