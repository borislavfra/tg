package generator

/*
import (
	"reflect"
	"strings"

	"github.com/vetcher/go-astra/types"
)

var (
	typeMap = map[string]string{
		reflect.Float32.String(): "number",
		reflect.Float64.String(): "number",
		reflect.Int.String():     "number",
		reflect.Int8.String():    "number",
		reflect.Int16.String():   "number",
		reflect.Int32.String():   "number",
		reflect.Int64.String():   "number",
		reflect.Uint.String():    "number",
		reflect.Uint8.String():   "number",
		reflect.Uint16.String():  "number",
		reflect.Uint32.String():  "number",
		reflect.Uint64.String():  "number",
		reflect.String.String():  "string",
		reflect.Bool.String():    "boolean",

		"uuid.UUID": "string",
	}
)

type templateMethod struct {
	types.Base
	InterfaceBase   types.Base
	Args            []types.Variable `json:"args,omitempty"`
	Results         []types.Variable `json:"results,omitempty"`
	ArgsTypesMap    map[int]string
	ResultsTypesMap map[int]string
}

func tsFunctionConverter(function *types.Function, interfaceBase types.Base, sw *swagger) (converterFunction *templateMethod) {
	newArgs := []types.Variable{}
	argsTypesMap := map[int]string{}
	for _, v := range function.Args {
		if typeMap[v.Type.String()] != "" {
			v.Type = types.TName{TypeName: typeMap[v.Type.String()]}
			newArgs = append(newArgs, v)
		} else if _, ok := v.Type.(types.TArray); ok {
			v.Type = types.TName{TypeName: "array"}
			newArgs = append(newArgs, v)
		} else {
			tmpType := v.Type.String()
			typeStrings := strings.Split(tmpType, ".")
			if len(typeStrings) == 2 {
				if _, ok := sw.schemas[typeStrings[1]]; ok {
					argsTypesMap[len(newArgs)] = typeStrings[1]
					v.Type = types.TName{TypeName: "object"}
					newArgs = append(newArgs, v)
				}
			}
		}
	}

	newResults := []types.Variable{}
	resultsTypesMap := map[int]string{}
	for _, v := range function.Results {
		if typeMap[v.Type.String()] != "" {
			v.Type = types.TName{TypeName: typeMap[v.Type.String()]}
			newResults = append(newResults, v)
		} else if arrayType, ok := v.Type.(types.TArray); ok {
			if typeMap[arrayType.Next.String()] != "" {
				resultsTypesMap[len(newResults)] = typeMap[arrayType.Next.String()]
				v.Type = types.TName{TypeName: "array"}
				newArgs = append(newArgs, v)
			} else {
				tmpType := v.Type.String()
				typeStrings := strings.Split(tmpType, ".")
				if len(typeStrings) == 2 {
					if _, ok := sw.schemas[typeStrings[1]]; ok {
						resultsTypesMap[len(newResults)] = typeStrings[1]
						v.Type = types.TName{TypeName: "object"}
						newResults = append(newResults, v)
					}
				}
			}
			resultsTypesMap[len(newResults)] = arrayType.Next.String()
			v.Type = types.TName{TypeName: "array"}
			newArgs = append(newArgs, v)
		} else {
			tmpType := v.Type.String()
			typeStrings := strings.Split(tmpType, ".")
			if len(typeStrings) == 2 {
				if _, ok := sw.schemas[typeStrings[1]]; ok {
					resultsTypesMap[len(newResults)] = typeStrings[1]
					v.Type = types.TName{TypeName: "object"}
					newResults = append(newResults, v)
				}
			}
		}
	}

	return &templateMethod{
		Base:            function.Base,
		InterfaceBase:   interfaceBase,
		Args:            newArgs,
		Results:         newResults,
		ArgsTypesMap:    argsTypesMap,
		ResultsTypesMap: resultsTypesMap,
	}
}
*/
