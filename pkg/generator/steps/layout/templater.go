package layout

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/generator/steps/layout/tpl"
)

type PathTransformerFunc func(context interface{}, path string) (string, error)

const (
	templateExtension = ".tpl"
)

type RenderParams struct {
	// Path to directory with templates tree
	TemplatesPath string

	// Path to save result
	TargetPath string

	// Allows to overwrite the path of file or directory before moving result to target directory
	PathTransformer PathTransformerFunc

	FuncMap template.FuncMap
}

func renderTemplate(context interface{}, fs embed.FS, tplPath string, targetPath string, funcMap template.FuncMap) error {
	tmpl, err := template.
		New(filepath.Base(tplPath)).
		Funcs(funcMap).
		ParseFS(fs, tplPath)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(targetPath), 0770); err != nil {
		return err
	}

	file, err := os.Create(targetPath)
	if err != nil {
		return err
	}

	err = tmpl.Execute(file, context)
	if err != nil {
		return err
	}

	return nil
}

func copyFile(fs embed.FS, path string, targetPath string) error {
	data, err := fs.ReadFile(path)
	if err != nil {
		return err
	}

	err = os.WriteFile(targetPath, data, 0644)

	if err != nil {
		return err
	}

	return nil
}

func RenderTemplateTree(ctx *gencontext.GenContext, model interface{}, params RenderParams) error {
	ctx.Logger.Infof("Template render: starting... TemplatesPath: %s. TargetPath: %s", params.TemplatesPath, params.TargetPath)

	assetsFs := tpl.GetTplFs()
	return fs.WalkDir(assetsFs, params.TemplatesPath, func(path string, d fs.DirEntry, err error) error {
		ctx.Logger.Infof("Template render: visiting %s", path)
		if err != nil {
			return err
		}

		destPath := strings.ReplaceAll(path, params.TemplatesPath, "")
		if params.PathTransformer != nil {
			destPath, err = params.PathTransformer(model, destPath)
			if err != nil {
				return err
			}
		}
		destPath = filepath.Join(params.TargetPath, destPath)

		if d.IsDir() {
			ctx.Logger.Infof("Template render: found dir %s. Creating: %s", path, destPath)
			return os.MkdirAll(destPath, 0755)
		}

		if filepath.Ext(path) == templateExtension {
			filePath := strings.ReplaceAll(destPath, templateExtension, "")
			ctx.Logger.Infof("Template render: found tpl %s. Creating: %s", path, filePath)
			return renderTemplate(model, assetsFs, path, filePath, params.FuncMap)
		}

		ctx.Logger.Infof("Template render: found file %s. Creating: %s", path, destPath)
		if err := copyFile(assetsFs, path, destPath); err != nil {
			return err
		}

		return nil
	})
}
