// Copyright (C) 2025  T-Force I/O
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, version 3 of the License.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg, err := buildConfig(true, filepath.Join(t.TempDir(), "prism.yaml"))
	if err != nil {
		t.Error("Error get default config", err)
	}
	if cfg.Logger.Level != "info" {
		t.Errorf("Wrong Logger.Level. Expected '%s' Actual '%s'", "info", cfg.Logger.Level)
	}
	if cfg.Server.HTTP.Port != 8080 {
		t.Errorf("Wrong Server.HTTP.Port. Expected '%d' Actual '%d'", 8080, cfg.Server.HTTP.Port)
	}
}

func TestFileConfig(t *testing.T) {
	contents := []string{
		"log:",
		"  level: debug",
	}
	f := filepath.Join(t.TempDir(), "prism.yaml")
	if err := os.WriteFile(f, []byte(strings.Join(contents, "\n")), 0664); err != nil {
		t.Fatal(err)
	}

	cfg, err := buildConfig(true, f)
	if err != nil {
		t.Error("Error get config from file", err)
	}
	if cfg.Logger.Level != "debug" {
		t.Errorf("Wrong Logger.Level. Expected '%s' Actual '%s'", "debug", cfg.Logger.Level)
	}
	if cfg.Server.HTTP.Port != 8080 {
		t.Errorf("Wrong Server.HTTP.Port. Expected '%d' Actual '%d'", 8080, cfg.Server.HTTP.Port)
	}
}

func TestEnvConfig(t *testing.T) {
	t.Setenv("TFPRISM_LOG_LEVEL", "warn")
	cfg, err := buildConfig(true, filepath.Join(t.TempDir(), "prism.yaml"))
	if err != nil {
		t.Error("Error get config from env", err)
	}
	if cfg.Logger.Level != "warn" {
		t.Errorf("Wrong Logger.Level. Expected '%s' Actual '%s'", "warn", cfg.Logger.Level)
	}
}
