//
// Copyright (C) 2026 Holger de Carne
//
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

package generator

import (
	"fmt"
	"os"

	"go.yaml.in/yaml/v4"
)

type ClientConfig struct {
	BasePkg     string `yaml:"base_pkg"`
	ClientType  string `yaml:"client_type"`
	WrapperType string `yaml:"wrapper_type"`
	WrapperFile string `yaml:"wrapper_file"`
}

func LoadClientConfig(file string) (*ClientConfig, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read client config %q (cause: %w)", file, err)
	}
	config := &ClientConfig{}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal client config %q (cause: %w)", file, err)
	}
	return config, nil
}

type MockConfig struct {
	BasePkg     string `yaml:"base_pkg"`
	ServiceType string `yaml:"service_type"`
	MockType    string `yaml:"mock_type"`
	MockFile    string `yaml:"mock_file"`
}

func LoadMockConfig(file string) (*MockConfig, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read mock config %q (cause: %w)", file, err)
	}
	config := &MockConfig{}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal mock config %q (cause: %w)", file, err)
	}
	return config, nil
}
