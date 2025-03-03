package flatbuffers

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"text/template"
	"xCelFlow/config"
	"xCelFlow/entities"
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
	render.Register("flatbuffers", newtsRender)
}

func newtsRender(render *render.Render) (render.IRender, error) {
	Schema := render.Schema.(*config.FlatbuffersSchema)

	if err := instance(Schema); err != nil {
		return nil, err
	}
	r := &FBSRender{Render: render, FlatbuffersSchema: Schema, tmpl: tmpl}

	return r, nil
}

func instance(schema *config.FlatbuffersSchema) error {
	once.Do(func() {
		tmpl = template.New("flatbuffers").Funcs(entities.FuncMap)
		for _, tmplStr := range []string{packageTemplate, dataSetTemplate, tailTemplate, fbTemplate} {
			if _, err := tmpl.Parse(tmplStr); err != nil {
				initErr = err
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
		return err
	}
	defer func() { _ = fileIO.Close() }()

	// 模板必备数据
	data := map[string]any{"Table": r}
	if err = r.tmpl.Execute(fileIO, data); err != nil {
		return err
	}

	return nil
}

// binExport 导出序列化数据
func (r *FBSRender) binExport() error {
	jsonParser, err := parser.CloneParser("json", r.Table)
	if err != nil {
		return err
	}

	jsonRender, err := render.NewRender("json", jsonParser.GetTable())
	if err != nil {
		return err
	}

	if err = jsonRender.Execute(); err != nil {
		return err
	}

	flatc := r.GetFlatc()
	binDir := r.GetBinDirectory()
	fbFilename := filepath.Join(r.GetFbsDirectory(), r.Filename())
	jr := jsonRender.(*json.JSONRender)
	jsonFilename := filepath.Join(jr.GetJsonDirectory(), jr.Filename())
	cmd := exec.Command(flatc, "--no-warnings", "--unknown-json", "-o", binDir, "-b", fbFilename, jsonFilename)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("error: %s", string(output))
	}

	return nil
}

func (r *FBSRender) Verify() error {
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
