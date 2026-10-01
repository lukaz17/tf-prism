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

package engine

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/tforceaio/tf-prism/config"
)

// ConfigModule handles configuration management commands.
type ConfigModule struct {
	cfg    *config.RootConfig
	logger zerolog.Logger
}

// Return new ConfigModule instance.
func NewConfigModule(c *Controller) *ConfigModule {
	return &ConfigModule{
		cfg:    c.Config,
		logger: c.ModuleLogger("config"),
	}
}

// Export active configuration to stdout or file.
func (m *ConfigModule) Export(outputPath string, diffOnly bool) error {
	m.logger.Info().
		Bool("diffOnly", diffOnly).
		Msg("Start export configuration.")

	k := koanf.New(".")
	k.Load(structs.Provider(m.cfg, "koanf"), nil)

	if diffOnly {
		defaults := config.DefaultConfig()
		for key, defVal := range defaults.All() {
			if val := k.Get(key); val != nil && reflect.DeepEqual(val, defVal) {
				k.Delete(key)
			}
		}
	}
	data, _ := k.Marshal(yaml.Parser())
	if outputPath == "" {
		fmt.Print(string(data))
	} else {
		err := m.writeFile(outputPath, data)
		if err != nil {
			return err
		}
		m.logger.Info().Str("outputPath", outputPath).
			Msg("Configuration exported successfully.")
	}

	return nil
}

// Get value of key from active configuration.
func (m *ConfigModule) Get(key string) error {
	if err := ValidateString(key, "key"); err != nil {
		return err
	}

	k := koanf.New(".")
	k.Load(structs.Provider(m.cfg, "koanf"), nil)

	if val := k.Get(key); val != nil {
		fmt.Printf("%v\n", val)
	}
	return nil
}

// Write value of key to active YAML config file.
func (m *ConfigModule) Set(key, value string) error {
	if err := ValidateString(key, "key"); err != nil {
		return err
	}
	if err := ValidateString(value, "value"); err != nil {
		return err
	}

	k := koanf.New(".")
	if _, err := os.Stat(m.cfg.ConfigFile); err == nil {
		err = k.Load(file.Provider(m.cfg.ConfigFile), yaml.Parser())
		if err != nil {
			return err
		}
	}
	err := k.Set(key, value)
	if err != nil {
		return err
	}
	data, _ := k.Marshal(yaml.Parser())
	err = m.writeFile(m.cfg.ConfigFile, data)
	if err != nil {
		return err
	}

	m.logger.Info().Str("key", key).Str("value", value).
		Msg("Configuration updated successfully.")
	return nil
}

// Decorator to log error occurred when calling handlers.
func (m *ConfigModule) logError(err error) {
	if err != nil {
		m.logger.Err(err).Msg("Unexpected error has occurred.")
	}
}

// Write data to a file safely.
func (m *ConfigModule) writeFile(path string, data []byte) error {
	dir := filepath.Dir(path)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", path, err)
	}

	return nil
}

// Define Cobra Command for Config module.
func ConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage application configuration.",
	}

	exportCmd := &cobra.Command{
		Use:   "export [output]",
		Short: "Export active configuration to stdout or file.",
		Run: func(cmd *cobra.Command, args []string) {
			c := InitApp()
			defer c.Close()
			m := NewConfigModule(c)
			output, _ := cmd.Flags().GetString("output")
			diff, _ := cmd.Flags().GetBool("diff")
			if len(args) >= 1 && output == "" {
				output = args[0]
			}
			m.logError(m.Export(output, diff))
		},
	}
	exportCmd.Flags().BoolP("diff", "d", false, "Export differences from default configuration only.")
	exportCmd.Flags().StringP("output", "o", "", "Output file path. Empty will output to stdout.")

	getCmd := &cobra.Command{
		Use:   "get <key>",
		Short: "Get the active value of a configuration key.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			c := NewController(true)
			defer c.Close()
			m := NewConfigModule(c)
			m.logError(m.Get(args[0]))
		},
	}

	setCmd := &cobra.Command{
		Use:   "set <key> <value>",
		Short: "Set a configuration value in the active YAML file.",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			c := NewController(true)
			defer c.Close()
			m := NewConfigModule(c)
			m.logError(m.Set(args[0], args[1]))
		},
	}

	cmd.AddCommand(exportCmd, getCmd, setCmd)

	return cmd
}
