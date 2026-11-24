//
// Copyright (C) 2026 Holger de Carne
//
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

package sync

import (
	"fmt"
	"os"

	"go.yaml.in/yaml/v4"
)

type Config struct {
	URL  string `yaml:"url"`
	File string `yaml:"file"`
}

func LoadConfig(file string) (*Config, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read sync config %q (cause: %w)", file, err)
	}
	config := &Config{}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal sync config %q (cause: %w)", file, err)
	}
	return config, nil
}
