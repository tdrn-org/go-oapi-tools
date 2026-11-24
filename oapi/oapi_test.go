//
// Copyright (C) 2026 Holger de Carne
//
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

package oapi_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tdrn-org/go-oapi-tools/oapi"
)

func TestGenerate(t *testing.T) {
	config, err := oapi.LoadConfig("testdata/oapi.conf")
	require.NoError(t, err)
	updated, err := oapi.RunConfig(config)
	require.NoError(t, err)
	require.False(t, updated)
}
