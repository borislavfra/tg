package generator

import (
	"github.com/vetcher/go-astra/types"
	"reflect"
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
	}
)

type templateMethod struct {
	types.Base
	InterfaceBase types.Base
	Args          []types.Variable `json:"args,omitempty"`
	Results       []types.Variable `json:"results,omitempty"`
}

func tsFunctionConverter(function *types.Function, interfaceBase types.Base) (converterFunction *templateMethod) {
	newArgs := []types.Variable{}
	for _, v := range function.Args {
		if typeMap[v.Type.String()] != "" {
			v.Type = types.TName{TypeName: typeMap[v.Type.String()]}
			newArgs = append(newArgs, v)
		}
	}

	newResults := []types.Variable{}
	for _, v := range function.Results {
		if typeMap[v.Type.String()] != "" {
			v.Type = types.TName{TypeName: typeMap[v.Type.String()]}
			newResults = append(newResults, v)
		}
	}

	return &templateMethod{
		Base:          function.Base,
		InterfaceBase: interfaceBase,
		Args:          newArgs,
		Results:       newResults,
	}
}
