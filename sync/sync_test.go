//
// Copyright (C) 2026 Holger de Carne
//
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

package sync_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tdrn-org/go-oapi-tools/sync"
)

func TestSyncSpec(t *testing.T) {
	urls := []string{
		"file://./testdata/spec.yaml",
		"https://raw.githubusercontent.com/tdrn-org/go-paperless-ngx/refs/heads/main/api/paperless-ngx.yaml",
	}
	for _, urlString := range urls {
		tmpDir := t.TempDir()
		file := filepath.Join(tmpDir, "spec.yaml")
		config := &sync.Config{
			URL:  urlString,
			File: file,
		}

		// Initial sync
		updated, err := sync.RunConfig(config)
		require.NoError(t, err)
		require.True(t, updated)

		// Initial sync
		updated, err = sync.RunConfig(config)
		require.NoError(t, err)
		require.False(t, updated)
	}
}
