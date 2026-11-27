//
// Copyright (C) 2026 Holger de Carne
//
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

package generator_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tdrn-org/go-oapi-tools/generator"
)

func TestLoadClientConfig(t *testing.T) {
	config, err := generator.LoadClientConfig("testdata/client.conf")
	require.NoError(t, err)
	require.NotNil(t, config)
}

func TestLoadMockConfig(t *testing.T) {
	config, err := generator.LoadClientConfig("testdata/mock.conf")
	require.NoError(t, err)
	require.NotNil(t, config)
}
