//
// Copyright (C) 2026 Holger de Carne
//
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

package sync

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/tdrn-org/go-oapi-tools/internal/files"
)

func RunConfig(config *Config) (bool, error) {
	url, err := url.Parse(config.URL)
	if err != nil {
		return false, fmt.Errorf("failed to parse sync URL %q (cause: %w)", config.URL, err)
	}
	return SyncSpec(url, config.File)
}

func SyncSpec(url *url.URL, file string) (bool, error) {
	slog.Info("syncing spec...", slog.Any("url", url), slog.String("file", file))
	src, err := readSrc(url)
	if err != nil {
		return false, err
	}
	return files.Update(file, src)
}

func readSrc(url *url.URL) ([]byte, error) {
	urlString := url.String()
	if url.Scheme == "" || url.Scheme == "file" {
		urlString, _ = strings.CutPrefix(urlString, "file://")
		src, err := os.ReadFile(urlString)
		if err != nil {
			return nil, fmt.Errorf("failed to read spec source %q (cause: %w)", urlString, err)
		}
		return src, nil
	}
	rsp, err := http.Get(urlString)
	if err != nil {
		return nil, fmt.Errorf("failed to access spec source %q (cause: %w)", urlString, err)
	}
	if rsp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to access spec source %q (status: %s)", urlString, rsp.Status)
	}
	defer rsp.Body.Close()
	src, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read spec source %q (cause: %w)", urlString, err)
	}
	return src, nil
}
