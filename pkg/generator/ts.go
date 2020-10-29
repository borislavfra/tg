package generator

import (
	"os"
	"path"

	"github.com/seniorGolang/tg/pkg/utils"
)

const (
	indexF             = "/index.ts"
	makeRequestConfigF = "/make-request-config.ts"
	typesF             = "/types.ts"
	requestF           = "/request.ts"
	apiCreatorF        = "/api-creator.ts"
	getSchemasF        = "/get-schemas.ts"
)

var TSBasicTypes = map[string]bool{
	"string":  true,
	"number":  true,
	"boolean": true,
}

type commonDirRenderer struct {
	name       string
	fileName   string
	renderFunc func(svc *service, path string) (err error)
}

type ts struct {
	*Transport

	schemas    tsSchemas
	knownTypes map[string]int
}

func (tsDoc *ts) render(outDir string) (err error) {

	for _, svcName := range tsDoc.serviceKeys() {

		svc := tsDoc.services[svcName]
		svc.Name = utils.ToLowerCamel(svc.Name)
		if _, err = os.Stat(path.Join(outDir, svc.Name)); os.IsNotExist(err) {
			if err = os.Mkdir(path.Join(outDir, svc.Name), os.ModePerm); err != nil {
				return
			}
		}

		for _, method := range svc.methods {
			tsDoc.registerStruct(method.requestStructName(), svc.pkgPath, method.tags, method.argumentsWithUploads())
			tsDoc.registerStruct(method.responseStructName(), svc.pkgPath, method.tags, method.results())
		}

		if _, err = os.Stat(outDir + indexF); os.IsNotExist(err) {
			if err = tsDoc.genIndex(svc.Methods, outDir+indexF, svc); err != nil {
				return
			}
		} else {
			if err = tsDoc.updateIndex(svc.Methods, outDir+indexF, svc); err != nil {
				return
			}
		}
		if _, err = os.Stat(outDir + apiCreatorF); os.IsNotExist(err) {
			if err = tsDoc.genApiCreator(svc.Methods, outDir+apiCreatorF, svc); err != nil {
				return
			}
		} else {
			if err = tsDoc.updateApiCreator(svc.Methods, outDir+apiCreatorF, svc); err != nil {
				return
			}
		}

		var commonDirs = []commonDirRenderer{
			{
				name:       "_methods",
				fileName:   indexF,
				renderFunc: tsDoc.genMethods,
			},
			{
				name:       "_schemas",
				fileName:   indexF,
				renderFunc: tsDoc.genSchemas,
			},
			{
				name:       "_types",
				fileName:   indexF,
				renderFunc: tsDoc.genTypes,
			},
			{
				name:       "batched",
				fileName:   "",
				renderFunc: tsDoc.genBatched,
			},
		}
		for _, dir := range commonDirs {
			filePath := path.Join(outDir, svc.Name, dir.name)
			if _, err = os.Stat(filePath); os.IsNotExist(err) {
				if err = os.Mkdir(filePath, os.ModePerm); err != nil {
					return
				}
			}
			if err = dir.renderFunc(svc, filePath+dir.fileName); err != nil {
				return
			}
		}

		for _, method := range svc.Methods {
			filePath := path.Join(outDir, svc.Name, utils.ToLowerCamel(method.Name))
			if _, err = os.Stat(filePath); os.IsNotExist(err) {
				if err = os.Mkdir(filePath, os.ModePerm); err != nil {
					return
				}
			}

			if err = tsDoc.genMakeRequestConfig(method, filePath+makeRequestConfigF, svc); err != nil {
				return
			}
			if err = tsDoc.genMethodTypes(filePath+typesF, svc.Name, method.Name); err != nil {
				return
			}
			if err = tsDoc.genMethodIndex(filePath + indexF); err != nil {
				return
			}
			if err = tsDoc.genRequest(filePath + requestF); err != nil {
				return
			}
		}
	}
	return
}

func newTS(tr *Transport) *ts {
	return &ts{
		Transport:  tr,
		schemas:    make(tsSchemas),
		knownTypes: make(map[string]int),
	}
}
