//
// Copyright (C) 2026 Holger de Carne
//
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

package main

import (
	"log/slog"
	"os"

	"github.com/tdrn-org/go-oapi-tools/config"
	"github.com/tdrn-org/go-oapi-tools/generator"
	"github.com/tdrn-org/go-oapi-tools/oapi"
	"github.com/tdrn-org/go-oapi-tools/sync"
)

func main() {
	generate(os.Args...)
}

func generate(args ...string) int {
	logger := slog.With(slog.String("cmd", "generate"))
	if len(args) != 2 {
		logger.Error("usage: generate <generate.conf>")
		return 1
	}
	configFile := args[1]
	config, err := config.LoadGenerateConfig(configFile)
	if err != nil {
		logger.Error("unable to load config", slog.String("configFile", configFile), slog.Any("err", err))
		return 1
	}
	if config.SyncEnabled {
		_, err = sync.RunConfig(&config.Sync)
		if err != nil {
			logger.Error("sync failure", slog.Any("err", err))
			return 1
		}
	}
	if config.OAPIEnabled {
		_, err = oapi.RunConfig(&config.OAPI)
		if err != nil {
			logger.Error("oapi failure", slog.Any("err", err))
			return 1
		}
	}
	if config.ClientEnabled {
		_, err = generator.RunClientConfig(&config.Client)
		if err != nil {
			logger.Error("generate client wrapper failure", slog.Any("err", err))
			return 1
		}
	}
	if config.MockEnabled {
		_, err = generator.RunMockConfig(&config.Mock)
		if err != nil {
			logger.Error("generate mock service failure", slog.Any("err", err))
			return 1
		}
	}
	return 0
}
