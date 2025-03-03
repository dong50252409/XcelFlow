package typescript

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"text/template"
	"xCelFlow/config"
	"xCelFlow/entities"
	"xCelFlow/parser"
	"xCelFlow/render"

	"github.com/stoewer/go-strcase"
)

type TSRender struct {
	*render.Render
	*config.TypeScriptSchema
	tmpl *template.Template
}

var (
	once    sync.Once
	initErr error
	tmpl    *template.Template
)

func init() {
	render.Register("typescript", newtsRender)
}

func newtsRender(render *render.Render) (render.IRender, error) {
	Schema := render.Schema.(*config.TypeScriptSchema)

	if err := instance(Schema); err != nil {
		return nil, err
	}

	r := &TSRender{Render: render, TypeScriptSchema: Schema, tmpl: tmpl}

	return r, nil
}

func instance(schema *config.TypeScriptSchema) error {
	once.Do(func() {
		// 准备模板中用的函数
		funcMap := entities.CloneFuncMap()
		funcMap["calOffset"] = func(index int) int {
			return 2*index + 4
		}

		tmpl = template.New("typescript").Funcs(funcMap)
		for _, tmplStr := range []string{headTemplate, baseClassTemplate, listClassTemplate, classTemplate, enumTemplate, tsTemplate} {
			if _, err := tmpl.Parse(tmplStr); err != nil {
				initErr = err
				return
			}
		}

		if err := os.MkdirAll(schema.GetTsDirectory(), os.ModePerm); err != nil {
			initErr = fmt.Errorf("导出路径创建失败 %s", err)
			return
		}
	})
	return initErr
}

func (r *TSRender) Execute() error {
	if err := r.clsExecute(); err != nil {
		return err
	}

	if err := r.fbsExecute(); err != nil {
		return err
	}

	return nil
}

// clsExecute 导出封装类文件
func (r *TSRender) clsExecute() error {
	fp := filepath.Join(r.GetTsDirectory(), r.Filename())
	fileIO, err := os.Create(fp)
	if err != nil {
		return err
	}
	defer func() { _ = fileIO.Close() }()

	data := map[string]any{"Table": r}
	if err = r.tmpl.Execute(fileIO, data); err != nil {
		return err
	}

	return nil
}

// fbsExecute 导出flatbuffers数据
func (r *TSRender) fbsExecute() error {
	fbsParser, err := parser.CloneParser("flatbuffers", r.Table)
	if err != nil {
		return err
	}

	fbsRender, err := render.NewRender("flatbuffers", fbsParser.GetTable())
	if err != nil {
		return err
	}

	if err = fbsRender.Execute(); err != nil {
		return err
	}
	return nil
}

func (r *TSRender) Verify() error {
	return nil
}

func (r *TSRender) Filename() string {
	return strcase.KebabCase(r.GetFilePrefix()+r.Name) + ".ts"
}

func (r *TSRender) ConfigName() string {
	return r.GetTableNamePrefix() + r.Name
}
