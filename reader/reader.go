package reader

import (
	"path/filepath"
	"strings"
	"xCelFlow/config"
)

type IReader interface {
	Read() ([][]string, error)
}

type Reader struct {
	IReader
	*config.TomlConfig
	Path string
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
		Path:       path,
		TomlConfig: &config.Config,
	})

	return r, nil
}
