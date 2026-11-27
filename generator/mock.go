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

func RunMockConfig(config *MockConfig) (bool, error) {
	return GenerateServiceMock(config.BasePkg, config.ServiceType, config.MockType, config.MockFile)
}

func GenerateServiceMock(basePkg, clientType, mockType, mockFile string) (bool, error) {
	slog.Info("generating service mock...", slog.Any("basePkg", basePkg), slog.String("clientType", clientType), slog.String("mockType", mockType))
	code, err := codegen.Generate(basePkg, clientType, mockType, &mockGenerator{})
	if err != nil {
		return false, err
	}
	return codegen.FormatAndUpdate(mockFile, code)
}

type mockGenerator struct {
}

//go:embed mock_preamble.txt
var mockPreamble string

func (g *mockGenerator) Preamble(out *strings.Builder, ctx *codegen.GeneratorContext) error {
	_, err := fmt.Fprintf(out, mockPreamble,
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

//go:embed mock_method.txt
var mockMethod string

func (g *mockGenerator) Method(out *strings.Builder, ctx *codegen.GeneratorContext, method *types.Func, signature *types.Signature) error {
	methodName := method.Name()
	funcSignature := codegen.FuncSignature(&ctx.DstTypeInfo.PkgInfo, methodName, signature)
	reflectCall := codegen.ReflectCall(&ctx.DstTypeInfo.PkgInfo, signature)
	reflectReturn := codegen.ReflectReturn(&ctx.DstTypeInfo.PkgInfo, signature)
	_, err := fmt.Fprintf(out, mockMethod,
		ctx.SrcTypeInfo.Name,
		ctx.SrcTypeInfo.FullName,
		ctx.SrcTypeInfo.PkgName,
		ctx.SrcTypeInfo.PkgPath,
		ctx.DstTypeInfo.Name,
		ctx.DstTypeInfo.FullName,
		ctx.DstTypeInfo.PkgName,
		ctx.DstTypeInfo.PkgPath,
		methodName,
		funcSignature,
		reflectCall,
		reflectReturn)
	return err
}

//go:embed mock_epilogue.txt
var mockEpilogue string

func (g *mockGenerator) Epilogue(out *strings.Builder, ctx *codegen.GeneratorContext) error {
	_, err := fmt.Fprintf(out, mockEpilogue,
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
