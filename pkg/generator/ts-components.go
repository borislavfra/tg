package generator

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/seniorGolang/tg/pkg/mod"
	"github.com/seniorGolang/tg/pkg/tags"
	"github.com/vetcher/go-astra"
	"github.com/vetcher/go-astra/types"
)

func (tsDoc *ts) registerStruct(name, pkgPath string, mTags tags.DocTags, fields []types.StructField) {

	if len(fields) == 0 {
		tsDoc.schemas[name] = tsSchema{Type: "object"}
		return
	}

	if tsDoc.schemas == nil {
		tsDoc.schemas = make(tsSchemas)
	}

	structType := types.Struct{
		Base: types.Base{Name: name, Docs: mTags.ToDocs()},
	}

	for _, field := range fields {
		field.Docs = mTags.Sub(field.Name).ToDocs()
		structType.Fields = append(structType.Fields, field)
	}
	tsDoc.schemas[name] = tsDoc.walkVariable(name, pkgPath, structType, mTags)
}

func (tsDoc *ts) walkVariable(typeName, pkgPath string, varType types.Type, varTags tags.DocTags) (schema tsSchema) {

	if len(varTags) > 0 {

		schema.Description = varTags.Value(tagDesc)
		if example := varTags.Value(tagExample); example != "" {

			var value interface{} = example
			_ = json.Unmarshal([]byte(example), &value)
			schema.Example = value
		}

		if format := varTags.Value(tagFormat); format != "" {
			schema.Format = format
		}

		if newType := varTags.Value(tagType); newType != "" {
			schema.Type = newType
			return
		}
	}

	if newType, format := castType(varType.String()); newType != varType.String() {

		schema.Type = newType
		schema.Format = format
		return
	}

	switch vType := varType.(type) {

	case types.TName:

		if types.IsBuiltin(varType) {
			schema.Type = vType.TypeName
			return
		}

		schema.Example = nil
		schema.Ref = fmt.Sprintf("#/components/schemas/%s", vType.String())

		if nextType := tsDoc.searchType(pkgPath, vType.TypeName); nextType != nil {
			if tsDoc.knownCount(vType.TypeName) < 2 {
				tsDoc.knownInc(vType.TypeName)
				tsDoc.schemas[vType.TypeName] = tsDoc.walkVariable(typeName, pkgPath, nextType, varTags)
			}
		}

	case types.TMap:

		schema.Type = "object"
		schema.AdditionalProperties = tsDoc.walkVariable(typeName, pkgPath, vType.Value, nil)

	case types.TArray:

		schema.Type = "array"
		schema.Maximum = vType.ArrayLen
		schema.Nullable = vType.IsSlice
		itemSchema := tsDoc.walkVariable(typeName, pkgPath, vType.Next, nil)
		schema.Items = &itemSchema

	case types.Struct:

		schema.Type = "object"
		schema.Properties = make(tsProperties)

		for _, field := range vType.Fields {
			if fieldName := jsonName(field); fieldName != "-" {
				schema.Properties[fieldName] = tsDoc.walkVariable(field.Name, pkgPath, field.Type, tags.ParseTags(field.Docs))
			}
		}

	case types.TImport:

		schema.Example = nil
		schema.Ref = fmt.Sprintf("#/components/schemas/%s", vType.Next)

		if nextType := tsDoc.searchType(vType.Import.Package, vType.Next.String()); nextType != nil {
			if tsDoc.knownCount(vType.Next.String()) < 2 {
				tsDoc.knownInc(vType.Next.String())
				tsDoc.schemas[vType.Next.String()] = tsDoc.walkVariable(typeName, vType.Import.Package, nextType, varTags)
			}
		}

	case types.TEllipsis:

		schema.Type = "array"
		itemSchema := tsDoc.walkVariable(typeName, pkgPath, vType.Next, varTags)
		schema.Items = &itemSchema

	case types.TPointer:

		return tsDoc.walkVariable(typeName, pkgPath, vType.Next, nil)

	case types.TInterface:

		schema.Type = "object"
		schema.Nullable = true

	default:
		tsDoc.log.WithField("type", vType).Error("unknown type")
		return
	}
	return
}

func (tsDoc *ts) searchType(pkg, name string) (retType types.Type) {

	if retType = tsDoc.parseType(pkg, name); retType == nil {

		pkgPath := mod.PkgModPath(pkg)

		if retType = tsDoc.parseType(pkgPath, name); retType == nil {

			pkgPath = path.Join("./vendor", pkg)

			if retType = tsDoc.parseType(pkgPath, name); retType == nil {

				pkgPath = tsDoc.trimLocalPkg(pkg)
				retType = tsDoc.parseType(pkgPath, name)
			}
		}
	}
	return
}

func (tsDoc *ts) parseType(relPath, name string) (retType types.Type) {

	pkgPath, _ := filepath.Abs(relPath)

	_ = filepath.Walk(pkgPath, func(filePath string, info os.FileInfo, err error) (retErr error) {

		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if !strings.HasSuffix(info.Name(), ".go") {
			return nil
		}

		var srcFile *types.File
		if srcFile, err = astra.ParseFile(filePath, astra.IgnoreConstants, astra.IgnoreMethods); err != nil {
			retErr = errors.Wrap(err, fmt.Sprintf("%s,%s", relPath, name))
			tsDoc.log.WithError(err).Errorf("parse file %s", filePath)
			return err
		}

		for _, typeInfo := range srcFile.Interfaces {

			if typeInfo.Name == name {
				retType = types.TInterface{Interface: &typeInfo}
				return
			}
		}

		for _, typeInfo := range srcFile.Types {

			if typeInfo.Name == name {
				retType = typeInfo.Type
				return
			}
		}

		for _, structInfo := range srcFile.Structures {

			if structInfo.Name == name {
				retType = structInfo
				return
			}
		}
		return
	})
	return
}

func (tsDoc *ts) trimLocalPkg(pkg string) (pgkPath string) {

	module := tsDoc.getModName()

	if module == "" {
		return pkg
	}

	moduleTokens := strings.Split(module, "/")
	pkgTokens := strings.Split(pkg, "/")

	if len(pkgTokens) < len(moduleTokens) {
		return pkg
	}

	pgkPath = path.Join(strings.Join(pkgTokens[len(moduleTokens):], "/"))
	return
}

func (tsDoc *ts) getModName() (module string) {

	modFile, err := os.OpenFile("go.mod", os.O_RDONLY, os.ModePerm)

	if err != nil {
		return
	}
	defer modFile.Close()

	rd := bufio.NewReader(modFile)
	if module, err = rd.ReadString('\n'); err != nil {
		return ""
	}
	module = strings.Trim(module, "\n")

	moduleTokens := strings.Split(module, " ")

	if len(moduleTokens) == 2 {
		module = strings.TrimSpace(moduleTokens[1])
	}
	return
}

func (tsDoc *ts) knownCount(typeName string) int {
	if _, found := tsDoc.knownTypes[typeName]; !found {
		return 0
	}
	return tsDoc.knownTypes[typeName]
}

func (tsDoc *ts) knownInc(typeName string) {
	if _, found := tsDoc.knownTypes[typeName]; !found {
		tsDoc.knownTypes[typeName] = 0
	}
	tsDoc.knownTypes[typeName]++
}
