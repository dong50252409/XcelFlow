package core

// IRender 渲染接口
type IRender interface {
	// Execute 执行导出
	Execute() error

	// Verify 验证导出结果
	Verify() error
}
