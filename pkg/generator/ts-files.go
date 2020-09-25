package generator

import (
	"github.com/seniorGolang/tg/pkg/utils"
	"os"
	"path"
	"text/template"
)

func (svc *service) renderTSFiles(outDir string) (err error) {
	indexT, err := readTemplate(&indexTemplate)
	if err != nil {
		return
	}
	responseSchemaT, err := readTemplate(&responseSchemaTemplate)
	if err != nil {
		return
	}
	makeRequestConfigT, err := readTemplate(&makeRequestConfigTemplate)
	if err != nil {
		return
	}

	if _, err = os.Stat(path.Join(outDir, svc.Name)); os.IsNotExist(err) {
		if err = os.Mkdir(path.Join(outDir, svc.Name), os.ModePerm); err != nil {
			return
		}
	}
	for _, v := range svc.Methods {
		convertedFunc := tsFunctionConverter(v, svc.Interface.Base)
		filePath := path.Join(outDir, svc.Name, convertedFunc.Name)
		if _, err = os.Stat(filePath); os.IsNotExist(err) {
			if err = os.Mkdir(filePath, os.ModePerm); err != nil {
				return
			}
		}
		if err = svc.genResponseSchema(responseSchemaT, convertedFunc, filePath); err != nil {
			return
		}
		if err = svc.genIndex(indexT, convertedFunc, filePath); err != nil {
			return
		}
		if err = svc.genMakeRequestConfig(makeRequestConfigT, convertedFunc, filePath); err != nil {
			return
		}
	}
	return
}

func (svc *service) genResponseSchema(t *template.Template, method *templateMethod, path string) (err error) {
	file, err := os.Create(path + "/response_schema.ts")
	if err != nil {
		return err
	}

	err = t.Execute(file, method.Results)
	if err != nil {
		return err
	}

	return
}

func (svc *service) genIndex(t *template.Template, method *templateMethod, path string) (err error) {
	file, err := os.Create(path + "/index.ts")
	if err != nil {
		return err
	}

	err = t.Execute(file, method)
	if err != nil {
		return err
	}

	return
}

func (svc *service) genMakeRequestConfig(t *template.Template, method *templateMethod, path string) (err error) {
	file, err := os.Create(path + "/make_request_config.ts")
	if err != nil {
		return err
	}

	err = t.Execute(file, method)
	if err != nil {
		return err
	}

	return
}

func readTemplate(templateData *string) (t *template.Template, err error) {
	var fm = template.FuncMap{
		"toLowCamel": utils.ToLowerCamel,
	}

	t, err = template.New("").Funcs(fm).Parse(*templateData)
	if err != nil {
		return nil, err
	}
	return
}
