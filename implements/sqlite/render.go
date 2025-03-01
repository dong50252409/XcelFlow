package sqlite

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"xCelFlow/render"

	"github.com/stoewer/go-strcase"
)

type SQLRender struct {
	*render.Render
	db *sql.DB    // 数据库连接
	mu sync.Mutex // 用于并发控制
}

func init() {
	render.Register("sqlite", newSQLRender)
}

func newSQLRender(render *render.Render) render.IRender {
	dbDir := render.ExportDir()
	if err := os.MkdirAll(dbDir, os.ModePerm); err != nil {
		// 由于 newSQLRender 不返回 error，这里记录错误日志
		fmt.Printf("failed to create directory: %v\n", err)
		return nil
	}

	dbPath := filepath.Join(dbDir, strcase.SnakeCase(render.Name)+".db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		fmt.Printf("failed to open database: %v\n", err)
		return nil
	}

	// 设置数据库连接参数
	db.SetMaxOpenConns(1) // SQLite 建议只使用一个连接

	return &SQLRender{
		Render: render,
		db:     db,
	}
}

func (r *SQLRender) Execute() error {
	if err := r.Render.ExecuteBefore(); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// 创建表（如果不存在）
	if err := r.createTableIfNotExists(); err != nil {
		return err
	}

	// 写入数据
	if err := r.writeData(); err != nil {
		return err
	}

	return nil
}

func (r *SQLRender) createTableIfNotExists() error {
	// 使用 CREATE TABLE IF NOT EXISTS 语句
	schema := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s)",
		r.ConfigName(),
		r.generateTableSchema())

	_, err := r.db.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to create table: %v", err)
	}
	return nil
}

func (r *SQLRender) writeData() error {
	// 开启事务
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	// 构建批量插入的 SQL
	var valueStrings []string
	var valueArgs []interface{}

	baseSQL := fmt.Sprintf("INSERT INTO %s (%s) VALUES ",
		r.ConfigName(),
		r.generateColumnNames())

	// 收集所有的值
	for _, row := range r.DataSetIter() {
		values := r.extractValues(row)
		valueStrings = append(valueStrings, r.generateValuePlaceholders())
		valueArgs = append(valueArgs, values...)
	}

	// 如果没有数据要插入，直接返回
	if len(valueStrings) == 0 {
		return nil
	}

	// 构建完整的插入语句
	stmt := baseSQL + strings.Join(valueStrings, ",")

	// 执行批量插入
	if _, err := tx.Exec(stmt, valueArgs...); err != nil {
		return fmt.Errorf("failed to insert rows: %v", err)
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

// 生成列名字符串
func (r *SQLRender) generateColumnNames() string {
	// 实现此方法返回所有列名，例如: "id, name, age"
	return ""
}

// 生成值占位符
func (r *SQLRender) generateValuePlaceholders() string {
	// 根据列数生成占位符，例如: "(?, ?, ?)"
	columnCount := len(r.Columns())
	placeholders := make([]string, columnCount)
	for i := range placeholders {
		placeholders[i] = "?"
	}
	return "(" + strings.Join(placeholders, ",") + ")"
}

func (r *SQLRender) generateTableSchema() string {
	return ""
}

func (r *SQLRender) Verify() error {
	return nil
}

func (r *SQLRender) ConfigName() string {
	return strcase.SnakeCase(r.Name)
}
