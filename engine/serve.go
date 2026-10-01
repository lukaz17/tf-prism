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

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/tforceaio/tf-prism/config"
	"github.com/tforceaio/tf-prism/db"
	"github.com/tforceaio/tf-prism/server"
)

// ServeModule handles serving the HTTP server.
type ServeModule struct {
	cfg     *config.RootConfig
	httpSrv *server.Server
	logger  zerolog.Logger
}

// Return new ServeModule.
func NewServeModule(logger zerolog.Logger, cfg *config.RootConfig) *ServeModule {
	return &ServeModule{
		cfg:    cfg,
		logger: logger.With().Str("module", "server").Str("cmd", "serve").Logger(),
	}
}

// Start listen for HTTP requests. Blocks until the server is stopped.
func (m *ServeModule) Start() error {
	ctx, err := db.Connect(m.cfg.Database.Type, m.cfg.Database.Uri)
	if err != nil {
		return fmt.Errorf("connect to database: %w", err)
	}
	defer ctx.Disconnect()

	if err := ctx.Migrate(); err != nil {
		return fmt.Errorf("run database migration: %w", err)
	}

	addr := fmt.Sprintf(":%d", m.cfg.Server.HTTP.Port)
	m.logger.Info().
		Str("addr", addr).
		Msg("Start listening for HTTP requests.")
	m.httpSrv = server.New(addr, ctx)
	return m.httpSrv.Start()
}

// Decorator to log error occurred when calling handlers.
func (m *ServeModule) logError(err error) {
	logProgramError(m.logger, err)
}

// Define Cobra Command for Serve module.
func ServeCmd() *cobra.Command {
	serveCmd := &cobra.Command{
		Use:   "serve",
		Short: "Start HTTP server.",
		Run: func(cmd *cobra.Command, args []string) {
			c := InitApp()
			defer c.Close()
			m := NewServeModule(c.Logger, c.Config)
			m.logError(m.Start())
		},
	}
	return serveCmd
}
