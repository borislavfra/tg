package generator

import (
	"os"
	"path"
	"sort"
	"strings"

	. "github.com/seniorGolang/tg/pkg/typescript"
	"github.com/seniorGolang/tg/pkg/utils"
)

const (
	indexF             = "/index.ts"
	responseSchemaF    = "/response-schema.ts"
	makeRequestConfigF = "/make-request-config.ts"
	typesF             = "/types.ts"
	requestF           = "/request.ts"
	apiCreatorF        = "/api-creator.ts"
)

var TSBasicTypes = map[string]bool{
	"string":  true,
	"number":  true,
	"boolean": true,
}

func (svc *service) renderTSFiles(outDir string, sw *swagger) (err error) {

	svc.Name = utils.ToLowerCamel(svc.Name)
	if _, err = os.Stat(path.Join(outDir, svc.Name)); os.IsNotExist(err) {
		if err = os.Mkdir(path.Join(outDir, svc.Name), os.ModePerm); err != nil {
			return
		}
	}

	for _, method := range svc.methods {
		sw.registerStruct(method.requestStructName(), svc.pkgPath, method.tags, method.argumentsWithUploads())
		sw.registerStruct(method.responseStructName(), svc.pkgPath, method.tags, method.results())
	}

	convertedMethods := []*templateMethod{}
	for _, originalMethod := range svc.Methods {
		convertedMethods = append(convertedMethods, tsFunctionConverter(originalMethod, svc.Interface.Base, sw))
	}
	sort.Slice(convertedMethods, func(i, j int) bool {
		return convertedMethods[i].InterfaceBase.Name < convertedMethods[j].InterfaceBase.Name
	})

	if _, err = os.Stat(outDir + indexF); os.IsNotExist(err) {
		if err = svc.genIndex(convertedMethods, outDir+indexF); err != nil {
			return
		}
	} else {
		if err = svc.updateIndex(convertedMethods, outDir+indexF); err != nil {
			return
		}
	}
	if _, err = os.Stat(outDir + apiCreatorF); os.IsNotExist(err) {
		if err = svc.genApiCreator(convertedMethods, outDir+apiCreatorF); err != nil {
			return
		}
	} else {
		if err = svc.updateApiCreator(convertedMethods, outDir+apiCreatorF); err != nil {
			return
		}
	}

	for _, method := range convertedMethods {
		filePath := path.Join(outDir, svc.Name, utils.ToLowerCamel(method.Name))
		if _, err = os.Stat(filePath); os.IsNotExist(err) {
			if err = os.Mkdir(filePath, os.ModePerm); err != nil {
				return
			}
		}
		responseSchemaName := "response" + svc.Name + method.Name

		if err = svc.genResponseSchema(sw.schemas[responseSchemaName], sw, filePath+responseSchemaF); err != nil {
			return
		}
		if err = svc.genMakeRequestConfig(method, filePath+makeRequestConfigF); err != nil {
			return
		}
		if err = svc.genTypes(method, sw, filePath+typesF); err != nil {
			return
		}
		if err = svc.genMethodIndex(method, filePath+indexF); err != nil {
			return
		}
		if err = svc.genRequest(method, filePath+requestF); err != nil {
			return
		}
	}

	return
}

func (svc *service) genResponseSchema(schema swSchema, sw *swagger, path string) (err error) {

	srcFile := NewFile()

	srcFile.Line().Add(Id("import Joi from '@hapi/joi';"))

	srcFile.Export().Const().Id("responseSchema").E().Id("Joi.object").Params(
		ValuesFunc(func(group *Group) {
			for prName, v := range schema.Properties {
				switch v.Type {
				case "":
					group.Id(prName).T().Id("Joi.object").Params(ValuesFunc(
						renderResponseSchema(v, sw),
					)).Dot("unknown").Call()
				case "array":
					processArrayForResponseScheme(group, v, prName, sw)
				default:
					group.Id(prName).T().Id("Joi").Dot(v.Type).Call().Dot("required").Call()
				}
			}
		}),
	).Dot("unknown").Call().Op(";")

	return srcFile.Save(path)
}

func (svc *service) genMakeRequestConfig(method *templateMethod, path string) (err error) {

	srcFile := NewFile()

	srcFile.Import("./response-schema", "responseSchema")
	srcFile.Import("./types", "RequestParamsType")

	srcFile.Const().Id("ENDPOINT").E().Id("'/" + method.InterfaceBase.Name + "/" + utils.ToLowerCamel(method.Name) + "'")
	srcFile.Line()
	srcFile.Export().Const().Id("makeRequestConfig").E().Params(
		Values(
			Id("additionalFetchParams"),
			Id("bodyParams"),
		).T().Id("RequestParamsType"),
	).Op("=>").Params(Values(
		Id("endpoint").T().Id("ENDPOINT"),
		Id("responseSchema"),
		Id("body").T().Values(
			Id("params").T().Id("bodyParams"),
		),
		Id("...additionalFetchParams"),
	)).Op(";")

	return srcFile.Save(path)
}

func (svc *service) genMethodIndex(method *templateMethod, path string) (err error) {

	srcFile := NewFile()

	srcFile.Export().Op("* ").Op("from ").Id("'./request'")
	srcFile.Export().Op("* ").Op("from ").Id("'./types'")

	return srcFile.Save(path)
}

func (svc *service) genTypes(method *templateMethod, sw *swagger, path string) (err error) {

	srcFile := NewFile()

	srcFile.Import(
		"@mihanizm56/fetch-api",

		"IResponse",
		"TranslateFunction",
		"ExtraValidationCallback",
		"ProgressOptions",
		"CustomSelectorDataType",
	)

	srcFile.Type().Id("ParamsType").E().BlockFunc(func(group *Group) {
		for id, v := range method.Args {
			if v.Type.String() == "object" {
				group.Id(v.Name).T().BlockFunc(func(group *Group) {
					for prName, pr := range sw.schemas[method.ArgsTypesMap[id]].Properties {
						if pr.Type == "" {
							group.Id(prName).T().BlockFunc(renderTypesSchema(pr, sw))
						} else if pr.Type == "array" || pr.Type == "object" {
							if pr.Nullable {
								group.Id(prName).T().Id("Array").Op("<").Id("{}").Op("|").Id("null").Op(">").Op(";")
							} else {
								group.Id(prName).T().Id("Array").Op("<").Id("{}").Op(">").Op(";")
							}
						} else if pr.Nullable {
							group.Id(prName).T().Id(pr.Type).Op("|").Id("null").Op(";")
						} else {
							group.Id(prName).T().Id(pr.Type).Op(";")
						}
					}
				})
			} else if v.Type.String() == "array" {
				group.Id(v.Name).T().Id("Array").Op("<").Id("{}").Op(">").Op(";")
			} else {
				group.Id(v.Name).T().Id(v.Type.String()).Op(";")
			}
		}
	}).Op(";")
	srcFile.Line()
	srcFile.Type().Id("FetchParamsType").E().Block(
		Id("translateFunction").Op("?").T().Id("TranslateFunction").Op(";"),
		Id("isErrorTextStraightToOutput").Op("?").T().Id("boolean").Op(";"),
		Id("extraValidationCallback").Op("?").T().Id("ExtraValidationCallback").Op(";"),
		Id("customTimeout").Op("?").T().Id("number").Op(";"),
		Id("abortRequestId").Op("?").T().Id("string").Op(";"),
		Id("progressOptions").Op("?").T().Id("ProgressOptions").Op(";"),
		Id("customSelectorData").Op("?").T().Id("CustomSelectorDataType").Op(";"),
		Id("selectData").Op("?").T().Id("string").Op(";"),
	).Op(";")
	srcFile.Line()
	srcFile.Export().Type().Id("RequestParamsType").E().Block(
		Id("bodyParams").T().Id("ParamsType").Op(";"),
		Id("additionalFetchParams").Op("?").T().Id("FetchParamsType").Op(";"),
	).Op(";")
	srcFile.Line()
	srcFile.Export().Type().Id("ResponseType").E().Id("IResponse").Op("&").Block(
		Id("data").T().BlockFunc(func(group *Group) {
			for id, v := range method.Results {
				if v.Type.String() == "object" {
					group.Id(v.Name).T().BlockFunc(func(group *Group) {
						for prName, pr := range sw.schemas[method.ResultsTypesMap[id]].Properties {
							if pr.Type == "" {
								group.Id(prName).T().BlockFunc(renderTypesSchema(pr, sw))
							} else if pr.Type == "array" || pr.Type == "object" {
								if pr.Nullable {
									group.Id(prName).T().Id("Array").Op("<").Id("{}").Op("|").Id("null").Op(">").Op(";")
								} else {
									group.Id(prName).T().Id("Array").Op("<").Id("{}").Op(">").Op(";")
								}
							} else if pr.Nullable {
								group.Id(prName).T().Id(pr.Type).Op("|").Id("null").Op(";")
							} else {
								group.Id(prName).T().Id(pr.Type).Op(";")
							}
						}
					})
				} else if v.Type.String() == "array" {
					group.Id(v.Name).T().Id("Array").Op("<").Id("{}").Op(">").Op(";")
				} else {
					group.Id(v.Name).T().Id(v.Type.String()).Op(";")
				}
			}
		}).Op(";"),
	).Op(";")

	return srcFile.Save(path)
}

func (svc *service) genRequest(method *templateMethod, path string) (err error) {

	srcFile := NewFile()

	srcFile.Import("@mihanizm56/fetch-api", "JSONRPCRequest", "IResponse")
	srcFile.Import("./make-request-config", "makeRequestConfig")
	srcFile.Import("./types", "RequestParamsType")

	srcFile.Export().Const().Id("request").E().Params(Id("values").T().Id("RequestParamsType")).T().Id("Promise").Generic("IResponse").Op("=>")
	srcFile.New(Id("JSONRPCRequest").Call().Dot("makeRequest").Call(Id("makeRequestConfig").Params(Id("values")))).Op(";")

	return srcFile.Save(path)
}

func (svc *service) genIndex(methods []*templateMethod, path string) (err error) {

	srcFile := NewFile()

	for _, method := range methods {
		srcFile.Export().Values(
			Id("RequestParamsType").As().Id(method.InterfaceBase.Name+method.Name+"ParamsType"),
			Id("ResponseType").As().Id(method.InterfaceBase.Name+method.Name+"ResponseType"),
		).Op(" from ").Id("'./" + method.InterfaceBase.Name + "/" + utils.ToLowerCamel(method.Name) + "'").Op(";")
	}

	srcFile.Export().Block(Id("APICreator").As().Op("default")).Op(" from ").Id("'./api-creator'").Op(";")

	return srcFile.Save(path)
}

func (svc *service) updateIndex(methods []*templateMethod, path string) (err error) {

	srcFile := NewFile()

	baseNames := []string{}
	for _, method := range methods {
		baseName := method.InterfaceBase.Name + method.Name
		baseNames = append(baseNames, baseName)

		srcFile.Export().Values(
			Id("RequestParamsType").As().Id(baseName+"ParamsType"),
			Id("ResponseType").As().Id(baseName+"ResponseType"),
		).Op(" from ").Id("'./" + method.InterfaceBase.Name + "/" + utils.ToLowerCamel(method.Name) + "'").Op(";")
	}

	if err = srcFile.CheckRepetition(path, baseNames); err != nil {
		return err
	}

	return srcFile.AppendAfter(path, "Type", ";")
}

func (svc *service) genApiCreator(methods []*templateMethod, path string) (err error) {

	srcFile := NewFile()

	for _, method := range methods {
		reqName := method.InterfaceBase.Name + method.Name + "Request"
		srcFile.Import("./"+method.InterfaceBase.Name+"/"+utils.ToLowerCamel(method.Name), "request as "+reqName)
	}
	srcFile.Line().Add()
	srcFile.Export().Const().Id("APICreator").E().BlockFunc(func(group *Group) {
		prevIface := methods[0].InterfaceBase.Name
		for i := 0; i < len(methods); i++ {
			group.Id(methods[i].InterfaceBase.Name).T().BlockFunc(func(group *Group) {
				for i < len(methods) {
					if methods[i].InterfaceBase.Name != prevIface {
						prevIface = methods[i].InterfaceBase.Name
						break
					}
					group.Id(methods[i].Name + "Request").T().Id(methods[i].InterfaceBase.Name + methods[i].Name + "Request").Op(",")
					i += 1
				}
			}).Op(",")
		}
	})

	return srcFile.Save(path)
}

func (svc *service) updateApiCreator(methods []*templateMethod, path string) (err error) {

	srcFile := NewFile()

	baseNames := []string{}
	for _, method := range methods {
		baseName := method.InterfaceBase.Name + method.Name
		baseNames = append(baseNames, baseName)

		reqName := baseName + "Request"
		srcFile.Import("./"+method.InterfaceBase.Name+"/"+utils.ToLowerCamel(method.Name), "request as "+reqName)
	}

	if err = srcFile.CheckRepetition(path, baseNames); err != nil {
		return err
	}

	if err = srcFile.AppendAfter(path, "import", "\n"); err != nil {
		return err
	}

	srcFile = NewFile()
	prevIface := methods[0].InterfaceBase.Name
	for i := 0; i < len(methods); i++ {
		srcFile.Id(methods[i].InterfaceBase.Name).T().BlockFunc(func(group *Group) {
			for i < len(methods) {
				if methods[i].InterfaceBase.Name != prevIface {
					prevIface = methods[i].InterfaceBase.Name
					break
				}
				group.Id(methods[i].Name + "Request").T().Id(methods[i].InterfaceBase.Name + methods[i].Name + "Request").Op(",")
				i += 1
			}
		}).Op(",")
	}

	return srcFile.AppendAfter(path, ":", "},")
}

func renderResponseSchema(schema swSchema, sw *swagger) func(group *Group) {
	return func(group *Group) {

		if schema.Properties == nil && schema.Type != "" {
			if schema.Nullable {
				group.Id("Joi").Dot(schema.Type).Call().Dot("required").Call().Dot("allow").Call(Id("null"))
			} else {
				group.Id("Joi").Dot(schema.Type).Call().Dot("required").Call()
			}
			return
		}

		schemaRefParts := strings.Split(schema.Ref, "/")
		schemaName := schemaRefParts[len(schemaRefParts)-1]
		schema = sw.schemas[schemaName]
		for prName, pr := range schema.Properties {
			switch pr.Type {
			case "":

				group.Id(prName).T().Id("Joi.object").Params(ValuesFunc(renderResponseSchema(pr, sw))).Dot("unknown").Call()
			case "array":
				processArrayForResponseScheme(group, pr, prName, sw)
			case "object":
				if pr.Nullable {
					group.Id(prName).T().Id("Joi.object").Params(ValuesFunc(renderResponseSchema(schema, sw))).Dot("unknown").Call().Dot("allow").Call(Id("null"))
				} else {
					group.Id(prName).T().Id("Joi.object").Params(ValuesFunc(renderResponseSchema(schema, sw))).Dot("unknown").Call()
				}
			default:
				if pr.Nullable {
					group.Id(prName).T().Id("Joi").Dot(pr.Type).Call().Dot("required").Call().Dot("allow").Call(Id("null"))
				} else {
					group.Id(prName).T().Id("Joi").Dot(pr.Type).Call().Dot("required").Call()
				}
			}
		}
		return
	}
}

func renderTypesSchema(schema swSchema, sw *swagger) func(group *Group) {
	return func(group *Group) {
		schemaRefParts := strings.Split(schema.Ref, "/")
		schemaName := schemaRefParts[len(schemaRefParts)-1]
		schema = sw.schemas[schemaName]
		for prName, pr := range schema.Properties {
			if pr.Type == "" {
				group.Id(prName).T().BlockFunc(renderTypesSchema(pr, sw))
			} else if pr.Type == "array" || pr.Type == "object" {
				if pr.Nullable {
					group.Id(prName).T().Id("Array").Op("<").Id("{}").Op("|").Id("null").Op(">").Op(";")
				} else {
					group.Id(prName).T().Id("Array").Op("<").Id("{}").Op(">").Op(";")
				}
			} else if pr.Nullable {
				group.Id(prName).T().Id(pr.Type).Op("|").Id("null").Op(";")
			} else {
				group.Id(prName).T().Id(pr.Type).Op(";")
			}
		}
		return
	}
}

func processArrayForResponseScheme(group *Group, pr swSchema, prName string, sw *swagger) {
	if isBasicType(pr.Items.Type) {
		if pr.Nullable {
			group.Id(prName).T().Id("Joi.array").Call().
				Dot("items").ParamsFunc(renderResponseSchema(*pr.Items, sw)).Dot("required").Call().Dot("allow").Call(Id("null"))
		} else {
			group.Id(prName).T().Id("Joi.array").Call().
				Dot("items").ParamsFunc(renderResponseSchema(*pr.Items, sw)).Dot("required").Call()
		}
	} else {
		if pr.Nullable {
			group.Id(prName).T().Id("Joi.array").Call().
				Dot("items").Params(
				Id("Joi.object").Call(ValuesFunc(
					renderResponseSchema(*pr.Items, sw)),
				)).Dot("required").Call().Dot("allow").Call(Id("null"))
		} else {
			group.Id(prName).T().Id("Joi.array").Call().
				Dot("items").Params(
				Id("Joi.object").Call(ValuesFunc(
					renderResponseSchema(*pr.Items, sw)),
				)).Dot("required").Call()
		}
	}
}

func isBasicType(name string) bool {
	return TSBasicTypes[name]
}
