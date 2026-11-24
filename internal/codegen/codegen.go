//
// Copyright (C) 2026 Holger de Carne
//
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

package codegen

import (
	"fmt"
	"go/format"
	"go/types"
	"strings"

	"github.com/tdrn-org/go-oapi-tools/internal/files"
	"golang.org/x/tools/go/packages"
)

func LoadType(pkgPattern, typeName string) (*types.Interface, error) {
	config := &packages.Config{
		Tests: true,
		Mode:  packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles | packages.NeedImports | packages.NeedTypes | packages.NeedTypesSizes | packages.NeedSyntax | packages.NeedTypesInfo,
	}
	pkgs, err := packages.Load(config, pkgPattern)
	if err != nil {
		return nil, fmt.Errorf("failed to load package pattern %q (cause: %w)", pkgPattern, err)
	}
	for _, pkg := range pkgs {
		obj := pkg.Types.Scope().Lookup(typeName)
		if obj == nil {
			continue
		}
		src, ok := obj.Type().Underlying().(*types.Interface)
		if !ok {
			return nil, fmt.Errorf("type %q is not an interface", typeName)
		}
		return src, nil
	}
	return nil, fmt.Errorf("failed to load type %q", typeName)
}

type GeneratorContext struct {
	SrcTypeInfo TypeInfo
	DstTypeInfo TypeInfo
}

type Generator interface {
	Preamble(out *strings.Builder, ctx *GeneratorContext) error
	Method(out *strings.Builder, ctx *GeneratorContext, method *types.Func, signature *types.Signature) error
	Epilogue(out *strings.Builder, ctx *GeneratorContext) error
}

func Generate(basePkg, srcType, dstType string, generator Generator) (string, error) {
	srcTypeInfo, err := splitType(srcType)
	if err != nil {
		return "", err
	}
	dstTypeInfo, err := splitType(dstType)
	if err != nil {
		return "", err
	}
	srcPkgPattern := resolvePkgPattern(basePkg, srcTypeInfo.PkgPath)
	srcInterface, err := LoadType(srcPkgPattern, srcTypeInfo.Name)
	if err != nil {
		return "", err
	}
	out := &strings.Builder{}
	ctx := &GeneratorContext{
		SrcTypeInfo: *srcTypeInfo,
		DstTypeInfo: *dstTypeInfo,
	}
	err = generator.Preamble(out, ctx)
	if err != nil {
		return "", err
	}
	for method := range srcInterface.Methods() {
		signature := method.Type().(*types.Signature)
		err = generator.Method(out, ctx, method, signature)
		if err != nil {
			return "", err
		}
	}
	err = generator.Epilogue(out, ctx)
	if err != nil {
		return "", err
	}
	return out.String(), nil
}

func splitType(t string) (*TypeInfo, error) {
	basePathLen := strings.LastIndex(t, "/")
	pathLen := strings.LastIndex(t, ".")
	if basePathLen < 0 || pathLen < 0 {
		return nil, fmt.Errorf("invalid type name %q", t)
	}
	typeInfo := &TypeInfo{
		Name:     t[pathLen+1:],
		FullName: t[basePathLen+1:],
		PkgInfo: PkgInfo{
			PkgName: t[basePathLen+1 : pathLen],
			PkgPath: t[:pathLen],
		},
	}
	return typeInfo, nil
}

func resolvePkgPattern(basePkg, srcFullPkg string) string {
	baseParts := strings.Split(basePkg, "/")
	srcParts := strings.Split(srcFullPkg, "/")
	basePartsLen := len(baseParts)
	srcPartsLen := len(srcParts)
	maxLen := max(basePartsLen, srcPartsLen)
	resolved := "."
	for i := range maxLen {
		if i < basePartsLen {
			if i < srcPartsLen {
				basePart := baseParts[i]
				srcPart := srcParts[i]
				if basePart != srcPart {
					resolved = joinPkgPath("..", resolved, srcPart)
				}
			} else /* i >= srcPartsLen */ {
				resolved = joinPkgPath(baseParts[i], resolved, "")
			}
		} else /* i < srcPartsLen */ {
			resolved = joinPkgPath("", resolved, srcParts[i])
		}
	}
	return resolved
}

func joinPkgPath(base, resolved, src string) string {
	joined := base
	if resolved != "" && resolved != "." {
		if joined != "" {
			joined = joined + "/" + resolved
		} else {
			joined = resolved
		}
	}
	if src != "" {
		if joined != "" {
			joined = joined + "/" + src
		} else {
			joined = src
		}
	}
	if joined == "" {
		joined = "."
	}
	return joined
}

func FormatAndUpdate(file, code string) (bool, error) {
	formatted, err := format.Source([]byte(code))
	if err != nil {
		updated, _ := files.Update(file, []byte(code))
		return updated, fmt.Errorf("failed to format generated code (cause: %w)", err)
	}
	return files.Update(file, formatted)
}

func FuncSignature(pkg *PkgInfo, name string, signature *types.Signature) fmt.Stringer {
	return stringerFunc(func() string {
		buffer := &strings.Builder{}
		buffer.WriteString(name)
		writeFuncSignatureParams(buffer, pkg, signature.Params())
		writeFuncSignatureResults(buffer, pkg, signature.Results())
		return buffer.String()
	})
}

func writeFuncSignatureParams(buffer *strings.Builder, pkg *PkgInfo, params *types.Tuple) {
	buffer.WriteString("(")
	if params != nil {
		paramsLen := params.Len()
		for i := range paramsLen {
			if i > 0 {
				buffer.WriteString(", ")
			}
			param := params.At(i)
			buffer.WriteString(param.Name())
			buffer.WriteString(" ")
			buffer.WriteString(localizeFuncSignatureType(pkg, param.Type(), i == paramsLen-1))
		}
	}
	buffer.WriteString(")")
}

func writeFuncSignatureResults(buffer *strings.Builder, pkg *PkgInfo, results *types.Tuple) {
	if results != nil {
		resultsLen := results.Len()
		if resultsLen > 0 {
			buffer.WriteString(" (")
			for i := range resultsLen {
				if i > 0 {
					buffer.WriteString(", ")
				}
				result := results.At(i)
				buffer.WriteString(result.Name())
				buffer.WriteString(" ")
				buffer.WriteString(localizeFuncSignatureType(pkg, result.Type(), false))
			}
			buffer.WriteString(")")
		}
	}
}

func localizeFuncSignatureType(pkg *PkgInfo, t types.Type, ellipsis bool) string {
	arrayOrSlicePrefix := "[]"
	if ellipsis {
		arrayOrSlicePrefix = "..."
	}
	switch t := t.(type) {
	case *types.Pointer:
		return "*" + localizeFuncSignatureType(pkg, t.Elem(), ellipsis)
	case *types.Slice:
		return arrayOrSlicePrefix + localizeFuncSignatureType(pkg, t.Elem(), ellipsis)
	case *types.Array:
		return arrayOrSlicePrefix + localizeFuncSignatureType(pkg, t.Elem(), ellipsis)
	default:
		{
			typeString := t.String()
			typeInfo, _ := splitType(typeString)
			if typeInfo == nil {
				return typeString
			}
			if typeInfo.PkgName == pkg.PkgName {
				return typeInfo.Name
			}
			return typeInfo.FullName
		}
	}
}

func FuncCall(pkg *PkgInfo, name string, signature *types.Signature) fmt.Stringer {
	return stringerFunc(func() string {
		buffer := &strings.Builder{}
		buffer.WriteString(name)
		buffer.WriteString("(")
		params := signature.Params()
		if params != nil {
			paramsLen := params.Len()
			for i := range paramsLen {
				if i > 0 {
					buffer.WriteString(", ")
				}
				param := params.At(i)
				writeFuncCallParam(buffer, param, i == paramsLen-1)
			}
		}
		buffer.WriteString(")")
		return buffer.String()
	})
}

func writeFuncCallParam(buffer *strings.Builder, param *types.Var, ellipsis bool) {
	buffer.WriteString(param.Name())
	if ellipsis {
		switch param.Type().(type) {
		case *types.Array, *types.Slice:
			buffer.WriteString("...")
		}
	}
}

type stringerFunc func() string

func (f stringerFunc) String() string {
	return f()
}
