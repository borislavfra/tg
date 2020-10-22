package generator

import (
	"os"
	"path"
	"strings"

	. "github.com/seniorGolang/tg/pkg/typescript"
	"github.com/seniorGolang/tg/pkg/utils"
	"github.com/vetcher/go-astra/types"
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

	if _, err = os.Stat(outDir + indexF); os.IsNotExist(err) {
		if err = svc.genIndex(svc.Methods, outDir+indexF); err != nil {
			return
		}
	} else {
		if err = svc.updateIndex(svc.Methods, outDir+indexF); err != nil {
			return
		}
	}
	if _, err = os.Stat(outDir + apiCreatorF); os.IsNotExist(err) {
		if err = svc.genApiCreator(svc.Methods, outDir+apiCreatorF); err != nil {
			return
		}
	} else {
		if err = svc.updateApiCreator(svc.Methods, outDir+apiCreatorF); err != nil {
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
		responseSchemaName := "response" + svc.Name + method.Name
		requestSchemaName := "request" + svc.Name + method.Name

		if err = svc.genResponseSchema(sw.schemas[responseSchemaName], sw, filePath+responseSchemaF); err != nil {
			return
		}
		if err = svc.genMakeRequestConfig(method, filePath+makeRequestConfigF); err != nil {
			return
		}
		if err = svc.genTypes(sw.schemas[requestSchemaName], sw.schemas[responseSchemaName], sw, filePath+typesF); err != nil {
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

func (svc *service) genMakeRequestConfig(method *types.Function, path string) (err error) {

	srcFile := NewFile()

	srcFile.Import("./response-schema", "responseSchema")
	srcFile.Import("./types", "RequestParamsType")

	srcFile.Const().Id("ENDPOINT").E().Id("'/" + svc.Name + "/" + utils.ToLowerCamel(method.Name) + "'")
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

func (svc *service) genMethodIndex(method *types.Function, path string) (err error) {

	srcFile := NewFile()

	srcFile.Export().Op("* ").Op("from ").Id("'./request'")
	srcFile.Export().Op("* ").Op("from ").Id("'./types'")

	return srcFile.Save(path)
}

func (svc *service) genTypes(request, response swSchema, sw *swagger, path string) (err error) {

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
		for prName, schema := range request.Properties {
			switch schema.Type {
			case "":
				group.Id(prName).T().BlockFunc(renderTypesSchema(schema, sw)).Op(";")
			case "array":
				if schema.Nullable {
					group.Id(prName).T().Id("Array").Op("<").Id("{}").Op("|").Id("null").Op(">").Op(";")
				} else {
					group.Id(prName).T().Id("Array").Op("<").Id("{}").Op(">").Op(";")
				}
			default:
				group.Id(prName).T().Id(schema.Type).Op(";")
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
			for prName, schema := range response.Properties {
				switch schema.Type {
				case "":
					group.Id(prName).T().BlockFunc(renderTypesSchema(schema, sw)).Op(";")
				case "array":
					if schema.Nullable {
						group.Id(prName).T().Id("Array").Op("<").Id("{}").Op("|").Id("null").Op(">").Op(";")
					} else {
						group.Id(prName).T().Id("Array").Op("<").Id("{}").Op(">").Op(";")
					}
				default:
					group.Id(prName).T().Id(schema.Type).Op(";")
				}
			}
		}).Op(";"),
	).Op(";")

	return srcFile.Save(path)
}

func (svc *service) genRequest(method *types.Function, path string) (err error) {

	srcFile := NewFile()

	srcFile.Import("@mihanizm56/fetch-api", "JSONRPCRequest", "IResponse")
	srcFile.Import("./make-request-config", "makeRequestConfig")
	srcFile.Import("./types", "RequestParamsType")

	srcFile.Export().Const().Id("request").E().Params(Id("values").T().Id("RequestParamsType")).T().Id("Promise").Generic("IResponse").Op("=>")
	srcFile.New(Id("JSONRPCRequest").Call().Dot("makeRequest").Call(Id("makeRequestConfig").Params(Id("values")))).Op(";")

	return srcFile.Save(path)
}

func (svc *service) genIndex(methods []*types.Function, path string) (err error) {

	srcFile := NewFile()

	for _, method := range methods {
		srcFile.Export().Values(
			Id("RequestParamsType").As().Id(svc.Name+method.Name+"ParamsType"),
			Id("ResponseType").As().Id(svc.Name+method.Name+"ResponseType"),
		).Op(" from ").Id("'./" + svc.Name + "/" + utils.ToLowerCamel(method.Name) + "'").Op(";")
	}

	srcFile.Export().Block(Id("APICreator").As().Op("default")).Op(" from ").Id("'./api-creator'").Op(";")

	return srcFile.Save(path)
}

func (svc *service) updateIndex(methods []*types.Function, path string) (err error) {

	srcFile := NewFile()

	baseNames := []string{}
	for _, method := range methods {
		baseName := svc.Name + method.Name
		baseNames = append(baseNames, baseName)

		srcFile.Export().Values(
			Id("RequestParamsType").As().Id(baseName+"ParamsType"),
			Id("ResponseType").As().Id(baseName+"ResponseType"),
		).Op(" from ").Id("'./" + svc.Name + "/" + utils.ToLowerCamel(method.Name) + "'").Op(";")
	}

	if err = srcFile.CheckRepetition(path, baseNames); err != nil {
		return err
	}

	return srcFile.AppendAfter(path, "Type", ";")
}

func (svc *service) genApiCreator(methods []*types.Function, path string) (err error) {

	srcFile := NewFile()

	for _, method := range methods {
		reqName := svc.Name + method.Name + "Request"
		srcFile.Import("./"+svc.Name+"/"+utils.ToLowerCamel(method.Name), "request as "+reqName)
	}
	srcFile.Line().Add()
	srcFile.Export().Const().Id("APICreator").E().Block(
		Id(svc.Name).T().BlockFunc(func(group *Group) {
			for i := range methods {
				group.Id(methods[i].Name + "Request").T().Id(svc.Name + methods[i].Name + "Request").Op(",")
			}
		}).Op(","),
	).Op(";")

	return srcFile.Save(path)
}

func (svc *service) updateApiCreator(methods []*types.Function, path string) (err error) {

	srcFile := NewFile()

	baseNames := []string{}
	for _, method := range methods {
		baseName := svc.Name + method.Name
		baseNames = append(baseNames, baseName)

		reqName := baseName + "Request"
		srcFile.Import("./"+svc.Name+"/"+utils.ToLowerCamel(method.Name), "request as "+reqName)
	}

	if err = srcFile.CheckRepetition(path, baseNames); err != nil {
		return err
	}

	if err = srcFile.AppendAfter(path, "import", "\n"); err != nil {
		return err
	}

	srcFile = NewFile()

	srcFile.Id(svc.Name).T().BlockFunc(func(group *Group) {
		for i := range methods {
			group.Id(methods[i].Name + "Request").T().Id(svc.Name + methods[i].Name + "Request").Op(",")
		}
	}).Op(",")

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
			switch pr.Type {
			case "":
				group.Id(prName).T().BlockFunc(renderTypesSchema(pr, sw)).Op(";")
			case "object":
				fallthrough
			case "array":
				if pr.Nullable {
					group.Id(prName).T().Id("Array").Op("<").Id("{}").Op("|").Id("null").Op(">").Op(";")
				} else {
					group.Id(prName).T().Id("Array").Op("<").Id("{}").Op(">").Op(";")
				}
			default:
				if pr.Nullable {
					group.Id(prName).T().Id(pr.Type).Op("|").Id("null").Op(";")
				} else {
					group.Id(prName).T().Id(pr.Type).Op(";")
				}
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
