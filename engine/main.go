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
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/tforce-io/tf-golib/opx"
)

var majorVersion = 0
var minorVersion = 1
var patchVersion = 0
var gitCommit, gitDate, gitBranch string

// Return the current application version with build metadata.
func version() string {
	originDate := time.Date(2026, time.September, 1, 0, 0, 0, 0, time.UTC)
	gitDate2, err := time.Parse("20060102", gitDate)
	buildDate := opx.Ternary(err == nil, gitDate2, time.Now().UTC())
	duration := buildDate.Sub(originDate)
	minor := minorVersion
	patch := strconv.Itoa(patchVersion)
	if gitBranch == "master" || gitBranch == "HEAD" {
		// do nothing
	} else if gitBranch == "release" {
		minor += 1
		patch = patch + "-rc"
	} else if strings.Contains(gitBranch, "feat/") {
		minor += 1
		patch = patch + "-dev"
	} else {
		patch = strconv.Itoa(patchVersion+1) + "-dev"
	}
	if gitCommit != "" && len(gitCommit) >= 8 {
		return fmt.Sprintf("%d.%d.%s.%d-%s", majorVersion, minor, patch, duration.Milliseconds()/int64(86400000), gitCommit[:8])
	}
	return fmt.Sprintf("%d.%d.%s.%d", majorVersion, minor, patch, duration.Milliseconds()/int64(86400000))
}

// Initialize configurations, loggings for internal modules, and display basic
// information about this invocation.
func InitApp() *Controller {
	cfg := NewController(true)

	pwd, _ := os.Getwd()
	exec, _ := os.Executable()

	fmt.Printf("\nTFprism v%s\n", version())
	gitDate2, _ := time.Parse("20060102", gitDate)
	buildDate := opx.Ternary(gitDate == "", time.Now().UTC(), gitDate2)
	fmt.Printf("Copyright (C) %d T-Force I/O\n", buildDate.Year())
	fmt.Printf("Licensed under GPL-3.0 license. See COPYING file along with this program for more details.\n\n")
	cfg.Logger.Info().Msg("-----------------")
	cfg.Logger.Info().Msgf("Working directory %s", pwd)
	cfg.Logger.Info().Msgf("Config directory %s", cfg.Config.ConfigDir)
	cfg.Logger.Info().Msgf("Executable file %s", exec)
	cfg.Logger.Info().Msgf("Portable mode %t", cfg.Config.IsPortable)
	cfg.Logger.Info().Msg("-----------------")

	return cfg
}

// Execute the program.
func Execute() {
	gitDate2, _ := time.Parse("20060102", gitDate)
	buildDate := opx.Ternary(gitDate == "", time.Now().UTC(), gitDate2)

	rootCmd := &cobra.Command{
		Use: "prism",
		Long: fmt.Sprintf(
			`TFprism v%s.
Copyright (C) %d T-Force I/O.
Licensed under GPL-3.0 license. See COPYING file along with this program for more details.`,
			version(),
			buildDate.Year()),
		Short:   "Cross platform command line utility.",
		Version: version(),
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	rootCmd.AddCommand(ServeCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
