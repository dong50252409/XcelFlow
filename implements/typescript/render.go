package typescript

import (
	"fmt"
	"github.com/stoewer/go-strcase"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
	"xCelFlow/config"
	"xCelFlow/entities"
	"xCelFlow/implements/flatbuffers"
	"xCelFlow/parser"
	"xCelFlow/render"
)

type TSRender struct {
	*render.Render
	Schema *config.TypeScriptSchema
}

func init() {
	render.Register("typescript", newtsRender)
}

func newtsRender(render *render.Render) render.IRender {
	return &TSRender{render, render.Schema.(*config.TypeScriptSchema)}
}

func (r *TSRender) Execute() error {
	if err := r.Render.ExecuteBefore(); err != nil {
		return err
	}

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
	clsDir := r.clsDir()
	if err := os.MkdirAll(clsDir, os.ModePerm); err != nil {
		return err
	}

	fp := filepath.Join(clsDir, r.Filename())
	fileIO, err := os.Create(fp)
	if err != nil {
		return err
	}
	defer func() { _ = fileIO.Close() }()

	// 必备数据
	data := map[string]any{"Table": r}

	// 解析模板字符串
	tmpl := template.New("typescript").Funcs(entities.FuncMap)

	for _, tmplStr := range []string{headTemplate, baseClassTemplate, innerClass1Template, innerClass2Template, enumTemplate, tsTemplate} {
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

	return nil
}

// fbsExecute 导出flatbuffers接口文件
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

	fr := fbsRender.(*flatbuffers.FBSRender)
	flatc := fr.Schema.GetFlatc()
	fbsDir := r.fbsDir()
	fbFilename := filepath.Join(fr.FbsDir(), fr.Filename())
	cmd := exec.Command(flatc, "--ts", "--no-warnings", "--ts-omit-entrypoint", "-o", fbsDir, fbFilename)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("error:%s", string(output))
	}

	return nil
}

func (r *TSRender) Verify() error {
	return nil
}

func (r *TSRender) Filename() string {
	return strcase.KebabCase(r.Schema.GetFilePrefix()+r.Name) + ".ts"
}

func (r *TSRender) ConfigName() string {
	return r.Schema.GetTableNamePrefix() + r.Name
}

func (r *TSRender) InnerConfigName() string {
	return strcase.UpperCamelCase("cfg_" + r.Name)
}

func (r *TSRender) clsDir() string {
	return filepath.Join(r.ExportDir(), r.Schema.GetTSClsDirectory())
}

func (r *TSRender) fbsDir() string {
	return filepath.Join(r.ExportDir(), r.Schema.GetTSFbsDirectory())
}
