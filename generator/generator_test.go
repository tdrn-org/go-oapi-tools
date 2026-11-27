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

func TestGenerateClientWrapper(t *testing.T) {
	config, err := generator.LoadClientConfig("testdata/client.conf")
	require.NoError(t, err)
	updated, err := generator.RunClientConfig(config)
	require.NoError(t, err)
	require.False(t, updated)
}

func TestGenerateServiceMock(t *testing.T) {
	config, err := generator.LoadMockConfig("testdata/mock.conf")
	require.NoError(t, err)
	updated, err := generator.RunMockConfig(config)
	require.NoError(t, err)
	require.False(t, updated)
}
