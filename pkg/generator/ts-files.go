package generator

import (
	. "github.com/seniorGolang/tg/pkg/typescript"
	"github.com/seniorGolang/tg/pkg/utils"
	"github.com/vetcher/go-astra/types"
)

func (tsDoc *ts) genResponseSchema(schema tsSchema, path string) (err error) {

	srcFile := NewFile()

	srcFile.Line().Add(Id("import Joi from '@hapi/joi';"))

	srcFile.Export().Const().Id("responseSchema").E().Id("Joi.object").Params(
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
	).Dot("unknown").Call().Op(";")

	return srcFile.Save(path)
}

func (tsDoc *ts) genMakeRequestConfig(method *types.Function, path string, svc *service) (err error) {

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

func (tsDoc *ts) genMethodIndex(path string) (err error) {

	srcFile := NewFile()

	srcFile.Export().Op("* ").Op("from ").Id("'./request'")
	srcFile.Export().Op("* ").Op("from ").Id("'./types'")

	return srcFile.Save(path)
}

func (tsDoc *ts) genTypes(request, response tsSchema, path string) (err error) {

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

	for _, method := range methods {
		srcFile.Export().Values(
			Id("RequestParamsType").As().Id(svc.Name+method.Name+"ParamsType"),
			Id("ResponseType").As().Id(svc.Name+method.Name+"ResponseType"),
		).Op(" from ").Id("'./" + svc.Name + "/" + utils.ToLowerCamel(method.Name) + "'").Op(";")
	}

	srcFile.Export().Block(Id("APICreator").As().Op("default")).Op(" from ").Id("'./api-creator'").Op(";")

	return srcFile.Save(path)
}

func (tsDoc *ts) updateIndex(methods []*types.Function, path string, svc *service) (err error) {

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

func (tsDoc *ts) genApiCreator(methods []*types.Function, path string, svc *service) (err error) {

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

func (tsDoc *ts) updateApiCreator(methods []*types.Function, path string, svc *service) (err error) {

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
