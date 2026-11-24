//
// Copyright (C) 2026 Holger de Carne
//
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

package files

import (
	"bytes"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
)

func Update(file string, data []byte) (bool, error) {
	logger := slog.With(slog.String("file", file))
	logger.Info("updating file...")
	old, err := os.ReadFile(file)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return false, fmt.Errorf("failed to read file %q (cause: %w)", file, err)
		}
		old = []byte{}
	}
	if bytes.Equal(old, data) {
		logger.Info("file up-to-date - skipping write")
		return false, nil
	}
	fileDir := filepath.Dir(file)
	err = os.MkdirAll(fileDir, 0o755)
	if err != nil {
		return true, fmt.Errorf("failed to create file directory %q (cause: %w)", fileDir, err)
	}
	err = os.WriteFile(file, data, 0o644)
	if err != nil {
		return true, fmt.Errorf("failed to write file %q (cause: %w)", file, err)
	}
	logger.Info("file updated")
	return true, nil
}
