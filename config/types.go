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

import "github.com/tforceaio/tf-prism/db"

// DbType is the type of database engine.
type DbType = db.DbType

// RootConfig contains all available configurations for the application.
type RootConfig struct {
	ConfigDir  string          `koanf:"-"`
	ConfigFile string          `koanf:"-"`
	IsPortable bool            `koanf:"-"`
	Database   *DatabaseConfig `koanf:"db"`
	Logger     *LoggerConfig   `koanf:"log"`
	Server     *ServerConfig   `koanf:"server"`
}

// DatabaseConfig contains configurations for database connection.
type DatabaseConfig struct {
	Type DbType `koanf:"type"` // Supported: postgres, mysql
	Uri  string `koanf:"uri"`
}

// HTTPConfig contains configurations for the HTTP server.
type HTTPConfig struct {
	Port int `koanf:"port"`
}

// LoggerConfig contains configurations for logging.
type LoggerConfig struct {
	Level string `koanf:"level"`
}

// ServerConfig contains configurations for the server component.
type ServerConfig struct {
	HTTP *HTTPConfig `koanf:"http"`
}
