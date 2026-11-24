//
// Copyright (C) 2026 Holger de Carne
//
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

package oapi

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/codegen"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/util"
	"github.com/tdrn-org/go-oapi-tools/internal/files"
	"go.yaml.in/yaml/v4"
)

func RunConfig(config *Config) (bool, error) {
	return GenerateSpecCode(config.SpecFile, config.OptsFile, config.CodeFile)
}

func GenerateSpecCode(specFile, optsFile, codeFile string) (bool, error) {
	slog.Info("generating spec code...", slog.Any("spec", specFile), slog.String("opts", optsFile))
	opts, err := loadOpts(optsFile)
	if err != nil {
		return false, err
	}
	swagger, err := loadSwagger(specFile, opts)
	if err != nil {
		return false, err
	}
	code, err := codegen.Generate(swagger, *opts)
	if err != nil {
		return false, fmt.Errorf("failed to generate spec code (cause: %w)", err)
	}
	return files.Update(codeFile, []byte(code))
}

func loadOpts(optsFile string) (*codegen.Configuration, error) {
	optsData, err := os.ReadFile(optsFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read codegen options from file %q (cause: %w)", optsFile, err)
	}
	opts := codegen.Configuration{}
	err = yaml.Unmarshal(optsData, &opts)
	if err != nil {
		return nil, fmt.Errorf("failed to decode codegen options from file %q (cause: %w)", optsFile, err)
	}
	opts = opts.UpdateDefaults()
	err = opts.Validate()
	if err != nil {
		return nil, fmt.Errorf("failed to validate codegen options from file %q (cause: %w)", optsFile, err)
	}
	for key, warning := range opts.Generate.Warnings() {
		slog.Warn(warning, slog.String("warning", key), slog.String("file", optsFile))
	}
	return &opts, nil
}

func loadSwagger(specFile string, opts *codegen.Configuration) (*openapi3.T, error) {
	overlayOpts := util.LoadSwaggerWithOverlayOpts{
		Path:                opts.OutputOptions.Overlay.Path,
		Strict:              true,
		LegacyJSONPath:      opts.Compatibility.OverlaysUseYamlJsonpath,
		AllowRemoteRefsFrom: opts.AllowRemoteRefsFrom,
	}
	swagger, err := util.LoadSwaggerWithOverlay(specFile, overlayOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to load spec file %q (cause: %w)", specFile, err)
	}
	return swagger, nil
}
