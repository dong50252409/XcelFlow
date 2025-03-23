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
	"xCelFlow/core"
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

func newErlangRender(render *render.Render) (core.IRender, error) {
	schema := render.Schema.(*config.ErlangSchema)

	if err := instance(schema); err != nil {
		return nil, fmt.Errorf("初始化失败 %s", err)
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
		hrlTmpl = template.New("hrl").Funcs(core.FuncMap)
		for _, tmplStr := range []string{hrlHeadTemplate, hrlRecordTemplate, hrlMacroTemplate, hrlTailTemplate, hrlTemplate} {
			if _, err := hrlTmpl.Parse(tmplStr); err != nil {
				initErr = fmt.Errorf("初始化hrl模板失败 %s", err)
				return
			}
		}

		// 初始化erl模板
		erlTmpl = template.New("erl").Funcs(core.FuncMap)
		for _, tmplStr := range []string{erlHeadTemplate, erlGetTemplate, erlListTemplate, erlTemplate} {
			if _, err := erlTmpl.Parse(tmplStr); err != nil {
				initErr = fmt.Errorf("初始化erl模板失败 %s", err)
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
	} else {
		if err := r.erlExecute(); err != nil {
			return err
		}
	}
	return nil
}

func (r *ERLRender) hrlExecute() error {
	fp := filepath.Join(r.GetHrlDirectory(), r.hrlFilename())
	fileIO, err := os.Create(fp)
	if err != nil {
		return fmt.Errorf("创建文件失败：%s 错误：%s", fp, err)
	}
	defer func() { _ = fileIO.Close() }()

	// 模板必备数据
	data := map[string]any{"Table": r}

	if err = r.hrlTmpl.Execute(fileIO, data); err != nil {
		return fmt.Errorf("执行模板失败 %s", err)
	}

	fmt.Printf("导出配置：%s\n", fp)
	return nil
}

func (r *ERLRender) erlExecute() error {
	fp := filepath.Join(r.GetErlDirectory(), r.erlFilename())
	fileIO, err := os.Create(fp)
	if err != nil {
		return fmt.Errorf("创建文件失败：%s 错误：%s", fp, err)
	}
	defer func() { _ = fileIO.Close() }()

	erlTmpl, err := r.erlTmpl.Clone()
	if err != nil {
		return fmt.Errorf("克隆模板失败 %s", err)
	}

	if err = r.loadCustomTemplate(erlTmpl); err != nil {
		return fmt.Errorf("导入自定义模板失败 %s", err)
	}

	data := map[string]any{"Table": r}

	if err = erlTmpl.Execute(fileIO, data); err != nil {
		return fmt.Errorf("执行模板失败 %s", err)
	}

	fmt.Printf("导出配置：%s\n", fp)
	return nil
}

func (r *ERLRender) loadCustomTemplate(template *template.Template) error {
	// 收集所有匹配的自定义模板内容
	contentList := make([]string, 0)

	// 处理精确匹配的模板
	if path, ok := r.customTemplates[r.Name]; ok && path != "" {
		if content, err := os.ReadFile(path); err != nil {
			return fmt.Errorf("读取自定义模板失败 %s", err)
		} else {
			contentList = append(contentList, string(content))
		}
	}

	// 处理通配匹配的模板
	for regx, path := range r.customTemplates {
		if strings.HasSuffix(r.Name, regx) {
			if content, err := os.ReadFile(path); err != nil {
				return fmt.Errorf("读取自定义模板失败 %s", err)
			} else {
				contentList = append(contentList, string(content))
			}
		}
	}

	// 将所有自定义内容包装到 CUSTOM_TEMPLATE 定义中
	if len(contentList) > 0 {
		wrappedContent := fmt.Sprintf(`{{- define "CUSTOM_TEMPLATE" -}}%s{{- end -}}`, strings.Join(contentList, "\n"))
		if _, err := template.Parse(wrappedContent); err != nil {
			return fmt.Errorf("解析自定义模板失败 %s", err)
		}
		customTemplate := fmt.Sprintf("%s\n\n{{ template \"CUSTOM_TEMPLATE\" . }}", erlTemplate)
		if _, err := template.Parse(customTemplate); err != nil {
			return fmt.Errorf("解析自定义模板失败 %s", err)
		}
	}

	return nil
}

// Verify 验证导出结果
func (r *ERLRender) Verify() error {
	if r.FieldLen == 0 {
		fmt.Printf("没有生成数据文件，跳过验证：%s\n", r.erlFilename())
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
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("验证失败 %s", string(output))
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
