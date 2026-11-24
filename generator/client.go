//
// Copyright (C) 2026 Holger de Carne
//
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

package generator

import (
	_ "embed"
	"fmt"
	"go/types"
	"log/slog"
	"strings"

	"github.com/tdrn-org/go-oapi-tools/internal/codegen"
)

func RunClientConfig(config *ClientConfig) (bool, error) {
	return GenerateClientWrapper(config.BasePkg, config.ClientType, config.WrapperType, config.WrapperFile)
}

func GenerateClientWrapper(basePkg, clientType, wrapperType, wrapperFile string) (bool, error) {
	slog.Info("generating client wrapper...", slog.Any("basePkg", basePkg), slog.String("clientType", clientType), slog.String("wrapperType", wrapperType))
	code, err := codegen.Generate(basePkg, clientType, wrapperType, &wrapperGenerator{})
	if err != nil {
		return false, err
	}
	return codegen.FormatAndUpdate(wrapperFile, code)
}

type wrapperGenerator struct {
}

//go:embed client_preamble.txt
var wrapperPreamble string

func (g *wrapperGenerator) Preamble(out *strings.Builder, ctx *codegen.GeneratorContext) error {
	_, err := fmt.Fprintf(out, wrapperPreamble,
		ctx.SrcTypeInfo.Name,
		ctx.SrcTypeInfo.FullName,
		ctx.SrcTypeInfo.PkgName,
		ctx.SrcTypeInfo.PkgPath,
		ctx.DstTypeInfo.Name,
		ctx.DstTypeInfo.FullName,
		ctx.DstTypeInfo.PkgName,
		ctx.DstTypeInfo.PkgPath)
	return err
}

//go:embed client_method.txt
var wrapperMethod string

func (g *wrapperGenerator) Method(out *strings.Builder, ctx *codegen.GeneratorContext, method *types.Func, signature *types.Signature) error {
	methodName := method.Name()
	funcSignature := codegen.FuncSignature(&ctx.DstTypeInfo.PkgInfo, stripWithResponseSuffix(methodName), signature)
	funcCall := codegen.FuncCall(&ctx.DstTypeInfo.PkgInfo, methodName, signature)
	_, err := fmt.Fprintf(out, wrapperMethod,
		ctx.SrcTypeInfo.Name,
		ctx.SrcTypeInfo.FullName,
		ctx.SrcTypeInfo.PkgName,
		ctx.SrcTypeInfo.PkgPath,
		ctx.DstTypeInfo.Name,
		ctx.DstTypeInfo.FullName,
		ctx.DstTypeInfo.PkgName,
		ctx.DstTypeInfo.PkgPath,
		funcSignature,
		funcCall)
	return err
}

func stripWithResponseSuffix(s string) string {
	stripped, _ := strings.CutSuffix(s, "WithResponse")
	return stripped
}

//go:embed client_epilogue.txt
var wrapperEpilogue string

func (g *wrapperGenerator) Epilogue(out *strings.Builder, ctx *codegen.GeneratorContext) error {
	_, err := fmt.Fprintf(out, wrapperEpilogue,
		ctx.SrcTypeInfo.Name,
		ctx.SrcTypeInfo.FullName,
		ctx.SrcTypeInfo.PkgName,
		ctx.SrcTypeInfo.PkgPath,
		ctx.DstTypeInfo.Name,
		ctx.DstTypeInfo.FullName,
		ctx.DstTypeInfo.PkgName,
		ctx.DstTypeInfo.PkgPath)
	return err
}
