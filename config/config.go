//
// Copyright (C) 2026 Holger de Carne
//
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

package config

import (
	"fmt"
	"os"

	"github.com/tdrn-org/go-oapi-tools/generator"
	"github.com/tdrn-org/go-oapi-tools/oapi"
	"github.com/tdrn-org/go-oapi-tools/sync"
	"go.yaml.in/yaml/v4"
)

type GenerateConfig struct {
	SyncEnabled   bool                   `yaml:"sync_enabled"`
	Sync          sync.Config            `yaml:"sync"`
	OAPIEnabled   bool                   `yaml:"oapi_enabled"`
	OAPI          oapi.Config            `yaml:"oapi"`
	ClientEnabled bool                   `yaml:"client_enabled"`
	Client        generator.ClientConfig `yaml:"client"`
	MockEnabled   bool                   `yaml:"mock_enabled"`
	Mock          generator.MockConfig   `yaml:"mock"`
}

func LoadGenerateConfig(file string) (*GenerateConfig, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read generate config %q (cause: %w)", file, err)
	}
	config := &GenerateConfig{}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal generate config %q (cause: %w)", file, err)
	}
	return config, nil
}
