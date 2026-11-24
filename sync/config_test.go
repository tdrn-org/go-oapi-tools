//
// Copyright (C) 2026 Holger de Carne
//
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

package sync_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tdrn-org/go-oapi-tools/sync"
)

func TestLoadConfig(t *testing.T) {
	config, err := sync.LoadConfig("testdata/sync.conf")
	require.NoError(t, err)
	require.NotNil(t, config)
}
