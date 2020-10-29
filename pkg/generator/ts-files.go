package generator

import (
	. "github.com/seniorGolang/tg/pkg/typescript"
	"github.com/seniorGolang/tg/pkg/utils"
	"github.com/vetcher/go-astra/types"
	"os"
	"path"
	"strings"
)

func (tsDoc *ts) genMakeRequestConfig(method *types.Function, path string, svc *service) (err error) {

	srcFile := NewFile()

	srcFile.Import("../_schemas", "SCHEMAS")
	srcFile.Import("./types", "RequestParamsType")

	srcFile.Const().Id("ENDPOINT").E().Id("'/" + svc.Name + "/" + utils.ToLowerCamel(method.Name) + "'")
	srcFile.Line()
	srcFile.Export().Const().Id("makeRequestConfig").E().Params(
		Values(
			Id("additionalFetchParams"),
			Id("bodyParams"),
		).T().Id("RequestParamsType")).Op("=>").Params(Values(
		Id("endpoint").T().Id("ENDPOINT"),
		Id("responseSchema").T().Id("SCHEMAS").Dot(utils.ToLowerCamel(method.Name)),
		Id("body").T().Values(
			Id("params").T().Id("bodyParams"),
		),
		Id("...additionalFetchParams"),
	)).Op(";")

	return srcFile.Save(path)
}

func (tsDoc *ts) genMethodIndex(path string) (err error) {

	srcFile := NewFile()

	srcFile.Export().Op("* ").Op("from ").Id("'./request'")
	srcFile.Export().Op("* ").Op("from ").Id("'./types'")

	return srcFile.Save(path)
}

func (tsDoc *ts) genMethodTypes(path string, svcName, methodName string) (err error) {

	requestParamsType := strings.Title(svcName) + methodName + "ParamsType"
	srcFile := NewFile()

	srcFile.Import(
		"@mihanizm56/fetch-api",

		"TranslateFunction",
		"ExtraValidationCallback",
		"ProgressOptions",
		"CustomSelectorDataType",
	)
	srcFile.Import("../_types", requestParamsType)
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
		Id("bodyParams").T().Id(requestParamsType).Op(";"),
		Id("additionalFetchParams").Op("?").T().Id("FetchParamsType").Op(";"),
	).Op(";")

	return srcFile.Save(path)
}

func (tsDoc *ts) genRequest(path string) (err error) {

	srcFile := NewFile()

	srcFile.Import("@mihanizm56/fetch-api", "JSONRPCRequest", "IResponse")
	srcFile.Import("./make-request-config", "makeRequestConfig")
	srcFile.Import("./types", "RequestParamsType")

	srcFile.Export().Const().Id("request").E().Params(Id("values").T().Id("RequestParamsType")).T().Id("Promise").Generic("IResponse").Op("=>")
	srcFile.New(Id("JSONRPCRequest").Call().Dot("makeRequest").Call(Id("makeRequestConfig").Params(Id("values")))).Op(";")

	return srcFile.Save(path)
}

func (tsDoc *ts) genIndex(methods []*types.Function, path string, svc *service) (err error) {

	srcFile := NewFile()

	srcFile.Export().Block(Id("APICreator").As().Op("default")).Op(" from ").Id("'./api-creator'").Op(";")

	srcFile.Export().ValuesFunc(func(group *Group) {
		for _, method := range methods {
			group.Id(strings.Title(svc.Name) + method.Name + "ParamsType")
			group.Id(strings.Title(svc.Name) + method.Name + "ResponseType")
		}
	}).Op(" from ").Id("'./" + svc.Name + "/" + "_types" + "'").Op(";")

	return srcFile.Save(path)
}

func (tsDoc *ts) updateIndex(methods []*types.Function, path string, svc *service) (err error) {

	srcFile := NewFile()

	baseNames := []string{}
	srcFile.Export().ValuesFunc(func(group *Group) {
		for _, method := range methods {
			baseName := svc.Name + method.Name
			baseNames = append(baseNames, baseName)

			group.Id(strings.Title(svc.Name) + method.Name + "ParamsType")
			group.Id(strings.Title(svc.Name) + method.Name + "ResponseType")
		}
	}).Op(" from ").Id("'./" + svc.Name + "/" + "_types" + "'").Op(";")

	if err = srcFile.CheckRepetition(path, baseNames); err != nil {
		return err
	}

	return srcFile.AppendAfter(path, "from", ";")
}

func (tsDoc *ts) genApiCreator(methods []*types.Function, path string, svc *service) (err error) {

	srcFile := NewFile()

	tmpMethods := append(methods, &types.Function{
		Base: types.Base{
			Name: "Batched",
		},
	})

	for i := range tmpMethods {
		reqName := svc.Name + tmpMethods[i].Name + "Request"
		srcFile.Import("./"+svc.Name+"/"+utils.ToLowerCamel(tmpMethods[i].Name), "request as "+reqName)
	}
	srcFile.Line().Add()
	srcFile.Export().Const().Id("APICreator").E().Block(
		Id(svc.Name).T().BlockFunc(func(group *Group) {
			for i := range tmpMethods {
				group.Id(tmpMethods[i].Name + "Request").T().Id(svc.Name + tmpMethods[i].Name + "Request").Op(",")
			}
		}).Op(","),
	).Op(";")

	return srcFile.Save(path)
}

func (tsDoc *ts) updateApiCreator(methods []*types.Function, path string, svc *service) (err error) {

	srcFile := NewFile()

	tmpMethods := append(methods, &types.Function{
		Base: types.Base{
			Name: "Batched",
		},
	})

	baseNames := []string{}
	for _, method := range tmpMethods {
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
		for i := range tmpMethods {
			group.Id(tmpMethods[i].Name + "Request").T().Id(svc.Name + tmpMethods[i].Name + "Request").Op(",")
		}
	}).Op(",")

	return srcFile.AppendAfter(path, ":", "},")
}

func (tsDoc *ts) genMethods(svc *service, path string) (err error) {

	srcFile := NewFile()

	srcFile.Export().Const().Id("METHODS").E().ValuesFunc(func(group *Group) {
		for i := range svc.methods {
			lowCaseName := utils.ToLowerCamel(svc.methods[i].Name)
			group.Id(lowCaseName).T().SingleQ(lowCaseName)
		}
	}).Op(";")

	return srcFile.Save(path)
}

func (tsDoc *ts) genSchemas(svc *service, path string) (err error) {
	srcFile := NewFile()

	schemas := map[string]tsSchema{}
	for i := range svc.methods {
		responseSchemaName := "response" + svc.Name + svc.methods[i].Name
		if schema, ok := tsDoc.schemas[responseSchemaName]; ok {
			schemas[svc.methods[i].Name] = schema
		}
	}
	srcFile.Line().Add(Id("import Joi from '@hapi/joi';"))
	srcFile.Export().Const().Id("SCHEMAS").E().ValuesFunc(func(group *Group) {
		for methodName, schema := range schemas {
			group.Id(utils.ToLowerCamel(methodName)).T().Id("Joi.object").Params(
				ValuesFunc(func(group *Group) {
					for prName, v := range schema.Properties {
						switch v.Type {
						case "":
							group.Id(prName).T().Id("Joi.object").Params(ValuesFunc(
								renderResponseSchema(v, tsDoc),
							)).Dot("unknown").Call()
						case "array":
							processArrayForResponseScheme(group, v, prName, tsDoc)
						default:
							group.Id(prName).T().Id("Joi").Dot(v.Type).Call().Dot("required").Call()
						}
					}
				}),
			).Dot("unknown").Call()
		}
	}).Op(";")

	return srcFile.Save(path)
}

func (tsDoc *ts) genTypes(svc *service, path string) (err error) {

	srcFile := NewFile()

	responseSchemas := map[string]tsSchema{}
	requestSchemas := map[string]tsSchema{}
	for i := range svc.methods {
		responseSchemaName := "response" + svc.Name + svc.methods[i].Name
		requestSchemaName := "request" + svc.Name + svc.methods[i].Name
		if schema, ok := tsDoc.schemas[responseSchemaName]; ok {
			responseSchemas[svc.methods[i].Name] = schema
		}
		if schema, ok := tsDoc.schemas[requestSchemaName]; ok {
			requestSchemas[svc.methods[i].Name] = schema
		}
	}

	srcFile.Import(
		"@mihanizm56/fetch-api",

		"IResponse",
	)

	for methodName, request := range requestSchemas {
		srcFile.Export().Type().Id(strings.Title(svc.Name) + methodName + "ParamsType").E().BlockFunc(func(group *Group) {
			for prName, schema := range request.Properties {
				switch schema.Type {
				case "":
					group.Id(prName).T().BlockFunc(renderTypesSchema(schema, tsDoc)).Op(";")
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
	}

	for methodName, response := range responseSchemas {
		srcFile.Export().Type().Id(strings.Title(svc.Name) + methodName + "ResponseType").E().Id("IResponse").Op("&").Block(
			Id("data").T().BlockFunc(func(group *Group) {
				for prName, schema := range response.Properties {
					switch schema.Type {
					case "":
						group.Id(prName).T().BlockFunc(renderTypesSchema(schema, tsDoc)).Op(";")
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
		srcFile.Line()
	}

	return srcFile.Save(path)
}

func (tsDoc *ts) genBatched(svc *service, batchedPath string) (err error) {

	err = tsDoc.genMethodIndex(batchedPath + indexF)
	if err != nil {
		return
	}
	err = tsDoc.genRequest(batchedPath + requestF)
	if err != nil {
		return
	}

	mrcFile := NewFile()
	mrcFile.Import("./types", "RequestParamsType")
	mrcFile.Import("./_utils/get-schemas", "getSchemas")
	mrcFile.Line()
	mrcFile.Const().Id("ENDPOINT").E().SingleQ("/" + svc.Name).Op(";")
	mrcFile.Line()
	mrcFile.Export().Const().Id("makeRequestConfig").E().Params(
		Values(
			Id("additionalFetchParams"),
			Id("bodyParams"),
		).T().Id("RequestParamsType")).Op("=>").Params(Values(
		Id("endpoint").T().Id("ENDPOINT"),
		Id("responseSchema").T().Id("getSchemas").Call(Id("bodyParams")),
		Id("body").T().Id("bodyParams"),
		Id("isBatchRequest").T().True(),
		Id("...additionalFetchParams"),
	)).Op(";")
	if err = mrcFile.Save(batchedPath + makeRequestConfigF); err != nil {
		return
	}

	typesFile := NewFile()
	typesFile.Import(
		"@mihanizm56/fetch-api",

		"IResponse",
		"TranslateFunction",
		"ExtraValidationCallback",
		"ProgressOptions",
		"CustomSelectorDataType",
	)
	typesFile.Import("../_methods", "METHODS")
	requestParams, responseParams := func() (requests, responses []string) {
		svcCap := strings.Title(svc.Name)
		for i := range svc.methods {
			requests = append(requests, svcCap+svc.methods[i].Name+"ParamsType")
			responses = append(responses, svcCap+svc.methods[i].Name+"ResponseType")
		}
		return
	}()
	typesFile.Import(
		"../_types",
		append(requestParams, responseParams...)...,
	)
	typesFile.Line()
	typesFile.Export().Type().Id("BatchedParamsType").E().Block(
		Id("method").T().Id("keyof typeof ").Id("METHODS").Op(";"),
		Id("params").T().UnionFunc(func(group *Group) {
			for i := range requestParams {
				group.Id(requestParams[i])
			}
		}).Op(";"),
	).Op(";")
	typesFile.Line()
	typesFile.Type().Id("FetchParamsType").E().Block(
		Id("translateFunction").Op("?").T().Id("TranslateFunction").Op(";"),
		Id("isErrorTextStraightToOutput").Op("?").T().Id("boolean").Op(";"),
		Id("extraValidationCallback").Op("?").T().Id("ExtraValidationCallback").Op(";"),
		Id("customTimeout").Op("?").T().Id("number").Op(";"),
		Id("abortRequestId").Op("?").T().Id("string").Op(";"),
		Id("progressOptions").Op("?").T().Id("ProgressOptions").Op(";"),
		Id("customSelectorData").Op("?").T().Id("CustomSelectorDataType").Op(";"),
		Id("selectData").Op("?").T().Id("string").Op(";"),
	).Op(";")
	typesFile.Line()
	typesFile.Export().Type().Id("RequestParamsType").E().Block(
		Id("bodyParams").T().Id("Array").Op("<").Id("BatchedParamsType").Op(">").Op(";"),
		Id("additionalFetchParams").Op("?").T().Id("FetchParamsType").Op(";"),
	).Op(";")
	typesFile.Line()
	typesFile.Export().Type().Id("ResponseType").E().Id("IResponse").Op("&").Block(
		Id("data").T().Id("Array").Op("<").UnionFunc(func(group *Group) {
			for i := range responseParams {
				group.Id(responseParams[i])
			}
		}).Op(">").Op(";"),
	).Op(";")
	if err = typesFile.Save(batchedPath + typesF); err != nil {
		return
	}

	filePath := path.Join(batchedPath, "_utils")
	if _, err = os.Stat(filePath); os.IsNotExist(err) {
		if err = os.Mkdir(filePath, os.ModePerm); err != nil {
			return
		}
	}
	gsFile := NewFile()
	gsFile.Import("../../_schemas", "SCHEMAS")
	gsFile.Import("../types", "BatchedParamsType")
	gsFile.Line()
	gsFile.Export().Const().Id("getSchemas").
		E().Params(Id("options").T().Id("Array").Op("<").Id("BatchedParamsType").Op(">")).Op("=>").
		Id("options").Dot("map").Params(Params(Block(Id("method"))).Op("=>").Id("SCHEMAS").Index(Id("method"))).Op(";")
	if err = gsFile.Save(filePath + getSchemasF); err != nil {
		return
	}

	return
}
