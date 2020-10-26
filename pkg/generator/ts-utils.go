package generator

import (
	"strings"

	. "github.com/seniorGolang/tg/pkg/typescript"
)

func renderResponseSchema(schema tsSchema, sw *ts) func(group *Group) {
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

func renderTypesSchema(schema tsSchema, sw *ts) func(group *Group) {
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

func processArrayForResponseScheme(group *Group, pr tsSchema, prName string, tsDoc *ts) {
	if isBasicType(pr.Items.Type) {
		if pr.Nullable {
			group.Id(prName).T().Id("Joi.array").Call().
				Dot("items").ParamsFunc(renderResponseSchema(*pr.Items, tsDoc)).Dot("required").Call().Dot("allow").Call(Id("null"))
		} else {
			group.Id(prName).T().Id("Joi.array").Call().
				Dot("items").ParamsFunc(renderResponseSchema(*pr.Items, tsDoc)).Dot("required").Call()
		}
	} else {
		if pr.Nullable {
			group.Id(prName).T().Id("Joi.array").Call().
				Dot("items").Params(
				Id("Joi.object").Call(ValuesFunc(
					renderResponseSchema(*pr.Items, tsDoc)),
				)).Dot("required").Call().Dot("allow").Call(Id("null"))
		} else {
			group.Id(prName).T().Id("Joi.array").Call().
				Dot("items").Params(
				Id("Joi.object").Call(ValuesFunc(
					renderResponseSchema(*pr.Items, tsDoc)),
				)).Dot("required").Call()
		}
	}
}

func isBasicType(name string) bool {
	return TSBasicTypes[name]
}
