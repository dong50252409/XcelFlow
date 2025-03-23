package flatbuffers

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"text/template"
	"xCelFlow/config"
	"xCelFlow/core"
	"xCelFlow/implements/json"
	"xCelFlow/parser"
	"xCelFlow/render"

	"github.com/stoewer/go-strcase"
)

type FBSRender struct {
	*render.Render
	*config.FlatbuffersSchema
	tmpl *template.Template
}

var (
	once    sync.Once
	tmpl    *template.Template
	initErr error
)

func init() {
	render.Register("flatbuffers", newFBSRender)
}

func newFBSRender(render *render.Render) (core.IRender, error) {
	Schema := render.Schema.(*config.FlatbuffersSchema)

	if err := instance(Schema); err != nil {
		return nil, fmt.Errorf("初始化失败 %s", err)
	}
	r := &FBSRender{Render: render, FlatbuffersSchema: Schema, tmpl: tmpl}

	return r, nil
}

func instance(schema *config.FlatbuffersSchema) error {
	once.Do(func() {
		tmpl = template.New("flatbuffers").Funcs(core.FuncMap)
		for _, tmplStr := range []string{packageTemplate, dataSetTemplate, tailTemplate, fbTemplate} {
			if _, err := tmpl.Parse(tmplStr); err != nil {
				initErr = fmt.Errorf("初始化失败 %s", err)
				return
			}
		}

		// 创建目录
		if err := os.MkdirAll(schema.GetFbsDirectory(), os.ModePerm); err != nil {
			initErr = fmt.Errorf("导出路径创建失败 %s", err)
			return
		}

		if err := os.MkdirAll(schema.GetBinDirectory(), os.ModePerm); err != nil {
			initErr = fmt.Errorf("导出路径创建失败 %s", err)
			return
		}
	})

	return initErr
}

func (r *FBSRender) Execute() error {
	if err := r.fbsExecute(); err != nil {
		return err
	}
	if err := r.binExport(); err != nil {
		return err
	}
	return nil
}

// fbsExecute 导出描述文件
func (r *FBSRender) fbsExecute() error {
	fp := filepath.Join(r.GetFbsDirectory(), r.Filename())
	fileIO, err := os.Create(fp)
	if err != nil {
		return fmt.Errorf("创建文件失败：%s 错误：%s", fp, err)
	}
	defer func() { _ = fileIO.Close() }()

	// 模板必备数据
	data := map[string]any{"Table": r}
	if err = r.tmpl.Execute(fileIO, data); err != nil {
		return fmt.Errorf("执行模板失败 %s", err)
	}

	fmt.Printf("导出配置：%s\n", fp)
	return nil
}

// binExport 导出序列化数据
func (r *FBSRender) binExport() error {
	jsonParser, err := parser.CloneParser("json", r.Table)
	if err != nil {
		return fmt.Errorf("克隆解析器失败 %s", err)
	}

	jsonRender, err := render.NewRender("json", jsonParser.GetTable())
	if err != nil {
		return fmt.Errorf("克隆渲染器失败 %s", err)
	}

	if err := jsonRender.Execute(); err != nil {
		return fmt.Errorf("执行渲染器失败 %s", err)
	}

	flatc := r.GetFlatc()
	binDir := r.GetBinDirectory()
	fbFilename := filepath.Join(r.GetFbsDirectory(), r.Filename())
	jr := jsonRender.(*json.JSONRender)
	jsonFilename := filepath.Join(jr.GetJsonDirectory(), jr.Filename())
	cmd := exec.Command(flatc, "--no-warnings", "--unknown-json", "-o", binDir, "-b", fbFilename, jsonFilename)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("执行命令失败 %s", string(output))
	}
	fmt.Printf("导出配置：%s\n", binDir)
	return nil
}

func (r *FBSRender) ConfigName() string {
	return strcase.UpperCamelCase(r.GetTableNamePrefix() + r.Name)
}

func (r *FBSRender) Filename() string {
	return strcase.SnakeCase(r.GetFilePrefix()+r.Name) + ".fbs"
}

func (r *FBSRender) Namespace() string {
	return r.GetNamespace()
}
