package erlang

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"text/template"
	"xCelFlow/config"
	"xCelFlow/entities"
	"xCelFlow/render"

	"github.com/stoewer/go-strcase"
)

type ERLRender struct {
	*render.Render
	*config.ErlangSchema
	hrlTmpl         *template.Template
	erlTmpl         *template.Template
	customTemplates map[string]string
}

var (
	once            sync.Once
	hrlTmpl         *template.Template
	erlTmpl         *template.Template
	customTemplates map[string]string
	initErr         error
)

func init() {
	render.Register("erlang", newErlangRender)
}

func newErlangRender(render *render.Render) (render.IRender, error) {
	schema := render.Schema.(*config.ErlangSchema)

	if err := instance(schema); err != nil {
		return nil, err
	}

	r := &ERLRender{
		Render:          render,
		ErlangSchema:    schema,
		hrlTmpl:         hrlTmpl,
		erlTmpl:         erlTmpl,
		customTemplates: customTemplates,
	}

	return r, nil
}

// instance 初始化模板
func instance(schema *config.ErlangSchema) error {
	once.Do(func() {
		// 初始化hrl模板
		hrlTmpl = template.New("hrl").Funcs(entities.FuncMap)
		for _, tmplStr := range []string{hrlHeadTemplate, hrlRecordTemplate, hrlMacroTemplate, hrlTailTemplate, hrlTemplate} {
			if _, err := hrlTmpl.Parse(tmplStr); err != nil {
				initErr = err
				return
			}
		}

		// 初始化erl模板
		erlTmpl = template.New("erl").Funcs(entities.FuncMap)
		for _, tmplStr := range []string{erlHeadTemplate, erlGetTemplate, erlListTemplate} {
			if _, err := erlTmpl.Parse(tmplStr); err != nil {
				initErr = err
				return
			}
		}

		// 载入自定义erl模板
		customTemplates = make(map[string]string)
		_ = filepath.Walk(schema.GetErlTemplates(), func(path string, info os.FileInfo, _ error) error {
			if !info.IsDir() {
				name := strings.TrimSuffix(info.Name(), filepath.Ext(info.Name()))
				customTemplates[name] = path
			}
			return nil
		})

		// 创建hrl目录
		if err := os.MkdirAll(schema.GetHrlDirectory(), os.ModePerm); err != nil {
			initErr = fmt.Errorf("导出路径创建失败 %s", err)
			return
		}

		// 创建erl目录
		if err := os.MkdirAll(schema.GetErlDirectory(), os.ModePerm); err != nil {
			initErr = fmt.Errorf("导出路径创建失败 %s", err)
			return
		}
	})

	return initErr
}

// Execute 执行导出
func (r *ERLRender) Execute() error {
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
	fp := filepath.Join(r.GetHrlDirectory(), r.hrlFilename())
	fileIO, err := os.Create(fp)
	if err != nil {
		return err
	}
	defer func() { _ = fileIO.Close() }()

	// 模板必备数据
	data := map[string]any{"Table": r}

	if err = r.hrlTmpl.Execute(fileIO, data); err != nil {
		return err
	}

	fmt.Printf("导出配置：%s\n", fp)
	return nil
}

func (r *ERLRender) erlExecute() error {
	fp := filepath.Join(r.GetErlDirectory(), r.erlFilename())
	fileIO, err := os.Create(fp)
	if err != nil {
		return err
	}
	defer func() { _ = fileIO.Close() }()

	erlTmpl, err := r.erlTmpl.Clone()
	if err != nil {
		return err
	}

	err = r.loadCustomTemplate(erlTmpl)
	if err != nil {
		return err
	}

	data := map[string]any{"Table": r}

	if err = erlTmpl.Execute(fileIO, data); err != nil {
		return err
	}

	fmt.Printf("导出配置：%s\n", fp)

	return err
}

func (r *ERLRender) loadCustomTemplate(template *template.Template) error {
	var err error
	if path, ok := r.customTemplates[r.Name]; ok && path != "" {
		if content, err := os.ReadFile(path); err == nil {
			if _, err := template.Parse(string(content)); err == nil {
				return nil
			}
		}
		return err
	}
	for regx, path := range r.customTemplates {
		if strings.HasSuffix(r.Name, regx) {
			if content, err := os.ReadFile(path); err == nil {
				if _, err := template.Parse(string(content)); err == nil {
					return nil
				}
			}
			return err
		}
	}
	_, err = template.Parse(erlTemplate)
	return err
}

// Verify 验证导出结果
func (r *ERLRender) Verify() error {
	if r.FieldLen == 0 {
		fmt.Printf("没有生成数据文件，跳过验证：%s\n", r.erlFilename())
		return nil
	}

	var out string
	switch runtime.GOOS {
	case "windows":
		out = os.Getenv("TEMP")
	default:
		out = "/dev/null"
	}
	hrlDir := r.GetHrlDirectory()
	erl := filepath.Join(r.GetErlDirectory(), r.erlFilename())

	fmt.Printf("开始验证生成结果：%s\n", erl)
	cmd := exec.Command("erlc", "-Werror", "-Wall", "-o", out, "-I", hrlDir, erl)
	if result, err := cmd.Output(); err != nil {
		fmt.Printf("%s", result)
		return err
	}
	fmt.Printf("验证完成：%s\n", erl)
	return nil
}

func (r *ERLRender) hrlFilename() string {
	return strcase.SnakeCase(r.GetTableNamePrefix()+r.Name) + ".hrl"
}

func (r *ERLRender) erlFilename() string {
	return strcase.SnakeCase(r.GetTableNamePrefix()+r.Name) + ".erl"
}

// ConfigName 配置名
func (r *ERLRender) ConfigName() string {
	return strcase.SnakeCase(r.GetTableNamePrefix() + r.Name)
}
