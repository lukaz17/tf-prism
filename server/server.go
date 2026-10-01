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

package server

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tforceaio/tf-prism/db"
)

// Server is a thin wrapper around http.Server with pre-registered routes.
type Server struct {
	engine *gin.Engine
	http   *http.Server
}

// Entrypoint for creating a new Server listening on addr.
func New(addr string, ctx *db.DbContext) *Server {
	s := &Server{}
	s.engine = s.router(ctx)
	s.http = &http.Server{
		Addr:    addr,
		Handler: s.engine,

		ReadHeaderTimeout: 5 * time.Second,
	}
	return s
}

// Build the route table served by this Server. Callable again to rebuild after
// mutating routes.
func (s *Server) Handler() http.Handler {
	return s.engine
}

// Start listening. Blocks until the server is stopped or fails.
func (s *Server) Start() error {
	err := s.http.ListenAndServe()
	if err == http.ErrServerClosed {
		return nil
	}
	return err
}

// Stop serving and release resources.
func (s *Server) Close() error {
	return s.http.Close()
}

// Build the route table served by this Server. Callable again to rebuild after
// mutating routes.
func (s *Server) router(ctx *db.DbContext) *gin.Engine {
	// gin.New() instead of Default: no built-in logger, project logs via zerolog.
	r := gin.New()
	r.Use(gin.Recovery())
	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})
	registerLinkRoutes(r, ctx)
	return r
}
