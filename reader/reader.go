package reader

import (
	"path/filepath"
	"strings"
	"xCelFlow/config"
	"xCelFlow/core"
)

type Reader struct {
	*config.TomlConfig
	Path string
}

var registry = make(map[string]func(reader *Reader) core.IReader)

// Register 注册文件读取器
func Register(key string, cls func(reader *Reader) core.IReader) {
	registry[key] = cls
}

func NewReader(path string) (core.IReader, error) {
	ext := filepath.Ext(path)[1:]
	cls, ok := registry[ext]
	if !ok {
		return nil, errorTableNotSupported(path)
	}

	r := cls(&Reader{
		TomlConfig: &config.Config,
		Path:       path,
	})

	return r, nil
}

func (r *Reader) CanRead(filename string) bool {
	ext := filepath.Ext(filename)
	return strings.EqualFold(ext, ".csv") || strings.EqualFold(ext, ".xlsx")
}
