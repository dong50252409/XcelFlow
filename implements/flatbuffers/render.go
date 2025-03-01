package flatbuffers

import (
	"fmt"
	"github.com/stoewer/go-strcase"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
	"xCelFlow/config"
	"xCelFlow/entities"
	"xCelFlow/implements/json"
	"xCelFlow/parser"
	"xCelFlow/render"
)

type FBSRender struct {
	*render.Render
	Schema *config.FlatbuffersSchema
}

func init() {
	render.Register("flatbuffers", newtsRender)
}

func newtsRender(render *render.Render) render.IRender {
	return &FBSRender{render, render.Schema.(*config.FlatbuffersSchema)}
}

func (r *FBSRender) Execute() error {
	if err := r.Render.ExecuteBefore(); err != nil {
		return err
	}

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
	fbsDir := r.FbsDir()
	if err := os.MkdirAll(fbsDir, os.ModePerm); err != nil {
		return err
	}

	fp := filepath.Join(fbsDir, r.Filename())
	fileIO, err := os.Create(fp)
	if err != nil {
		return err
	}
	defer func() { _ = fileIO.Close() }()

	// 必备数据
	data := map[string]any{"Table": r}

	// 解析模板字符串
	tmpl := template.New("flatbuffers").Funcs(entities.FuncMap)

	for _, tmplStr := range []string{dataSetTemplate, tailTemplate, fbTemplate} {
		if tmpl, err = tmpl.Parse(tmplStr); err != nil {
			return err
		}
	}

	// 执行模板渲染并输出到文件
	if err = tmpl.Execute(fileIO, data); err != nil {
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

	flatc := r.Schema.GetFlatc()
	binDir := r.BinDir()
	fbFilename := filepath.Join(r.FbsDir(), r.Filename())
	jr := jsonRender.(*json.JSONRender)
	jsonFilename := filepath.Join(jr.ExportDir(), jr.Filename())
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
	return strcase.UpperCamelCase(r.Schema.GetTableNamePrefix() + r.Name)
}

func (r *FBSRender) Filename() string {
	return strcase.SnakeCase(r.Schema.GetFilePrefix()+r.Name) + ".fbs"
}

func (r *FBSRender) Namespace() string {
	return r.Schema.GetNamespace()
}

// FbsDir flatbuffers描述文件目录
func (r *FBSRender) FbsDir() string {
	return filepath.Join(r.ExportDir(), r.Schema.GetFbsDirectory())
}

// BinDir flatbuffers序列化数据目录
func (r *FBSRender) BinDir() string {
	return filepath.Join(r.ExportDir(), r.Schema.GetBinDirectory())
}
