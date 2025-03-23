package core

// IReader 文件读取接口
type IReader interface {
	// Read 读取文件
	Read() ([][]string, error)

	// CanRead 检查是否能处理该类型的数据
	CanRead(filename string) bool
}
