package erlang

import (
	"fmt"
	"github.com/stoewer/go-strcase"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"
	"xCelFlow/config"
	"xCelFlow/entities"
	"xCelFlow/render"
)

type ERLRender struct {
	*render.Render
	Schema    *config.ErlangSchema
	templates map[string]string
}

func init() {
	render.Register("erlang", newErlangRender)
}

func newErlangRender(render *render.Render) render.IRender {
	schema := render.Schema.(*config.ErlangSchema)
	templateDir := filepath.Join(render.ExportDir(), schema.GetErlTemplates())
	templates := getTemplates(templateDir)
	r := &ERLRender{render, schema, templates}
	return r
}

func getTemplates(templateDir string) map[string]string {
	templates := make(map[string]string)
	_ = filepath.Walk(templateDir, func(path string, info os.FileInfo, _ error) error {
		if info.IsDir() {
			return nil
		}
		if matched, _ := filepath.Match("*.tmpl", info.Name()); matched {
			name := strings.TrimSuffix(info.Name(), filepath.Ext(info.Name()))
			templates[name] = path
		}
		return nil
	})
	return templates
}

// Execute 执行导出
func (r *ERLRender) Execute() error {
	if err := r.Render.ExecuteBefore(); err != nil {
		return err
	}
	if err := r.hrlExecute(); err != nil {
		return err
	}

	if r.FieldLen == 0 {
		fmt.Printf("没有定义字段名跳过生成数据文件：%s\n", r.erlFilename())
	} else if err := r.erlExecute(); err != nil {
		return err
	}

	return nil
}

func (r *ERLRender) hrlExecute() error {
	hrlDir := r.hrlDir()
	if err := os.MkdirAll(hrlDir, os.ModePerm); err != nil {
		return err
	}

	fp := filepath.Join(hrlDir, r.hrlFilename())
	fileIO, err := os.Create(fp)
	if err != nil {
		return err
	}
	defer func() { _ = fileIO.Close() }()

	// 必备数据
	data := map[string]any{"Table": r}

	// 解析模板字符串
	tmpl := template.New("hrl").Funcs(entities.FuncMap)

	for _, tmplStr := range []string{hrlHeadTemplate, hrlRecordTemplate, hrlMacroTemplate, hrlTailTemplate, hrlTemplate} {
		tmpl, err = tmpl.Parse(tmplStr)
		if err != nil {
			return err
		}
	}

	// 执行模板渲染并输出到文件
	err = tmpl.Execute(fileIO, data)
	if err != nil {
		return err
	}

	fmt.Printf("导出配置：%s\n", fp)

	return nil
}

func (r *ERLRender) erlExecute() error {
	erlDir := r.erlDir()
	if err := os.MkdirAll(erlDir, os.ModePerm); err != nil {
		return err
	}

	fp := filepath.Join(erlDir, r.erlFilename())
	fileIO, err := os.Create(fp)
	if err != nil {
		return err
	}
	defer func() { _ = fileIO.Close() }()

	// 必备数据
	data := map[string]any{"Table": r}

	// 解析模板字符串
	tmpl := template.New("erl").Funcs(entities.FuncMap)

	for _, tmplStr := range r.GetTemplateList() {
		if tmpl, err = tmpl.Parse(tmplStr); err != nil {
			return err
		}
	}

	// 执行模板渲染并输出到文件
	if err = tmpl.Execute(fileIO, data); err != nil {
		return err
	}

	fmt.Printf("导出配置：%s\n", fp)

	return nil
}

func (r *ERLRender) GetTemplateList() []string {
	for regx, path := range r.templates {
		if strings.HasSuffix(r.Name, regx) {
			if content, err := os.ReadFile(path); err != nil {
				continue
			} else {
				return []string{erlHeadTemplate, erlGetTemplate, erlListTemplate, string(content)}
			}
		}
	}
	return []string{erlHeadTemplate, erlGetTemplate, erlListTemplate, erlTemplate}
}

// Verify 验证导出结果
func (r *ERLRender) Verify() error {
	if r.FieldLen == 0 {
		fmt.Printf("没有生成数据文件，跳过验证：%s\n", r.erlFilename())
		return nil
	}

	hrlDir := r.hrlDir()
	erl := filepath.Join(r.erlDir(), r.erlFilename())
	fmt.Printf("开始验证生成结果：%s\n", erl)
	var out string
	switch runtime.GOOS {
	case "windows":
		out = os.Getenv("TEMP")
	default:
		out = "/dev/null"
	}
	cmd := exec.Command("erlc", "-Werror", "-Wall", "-o", out, "-I", hrlDir, erl)
	result, err := cmd.Output()
	if err != nil {
		fmt.Printf("%s", result)
		return err
	}
	return nil
}

func (r *ERLRender) hrlDir() string {
	return filepath.Join(r.ExportDir(), r.Schema.GetHrlDirectory())
}

func (r *ERLRender) erlDir() string {
	return filepath.Join(r.ExportDir(), r.Schema.GetErlDirectory())
}

func (r *ERLRender) hrlFilename() string {
	return strcase.SnakeCase(r.Schema.GetTableNamePrefix()+r.Name) + ".hrl"
}

func (r *ERLRender) erlFilename() string {
	return strcase.SnakeCase(r.Schema.GetTableNamePrefix()+r.Name) + ".erl"
}

// ConfigName 配置名
func (r *ERLRender) ConfigName() string {
	return strcase.SnakeCase(r.Schema.GetTableNamePrefix() + r.Name)
}
