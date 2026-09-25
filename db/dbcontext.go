// Copyright (C) 2025 T-Force I/O
// This file is part of TFprism
//
// TFprism is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// TFprism is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with TFprism. If not, see <https://www.gnu.org/licenses/>.

package db

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	DbPostgres DbType = "postgres"
	DbMySQL    DbType = "mysql"
)

// DbContext encapsulate all actions related to reading from and writing to database.
type DbContext struct {
	db  *gorm.DB
	uri string
}

// Return new DbContext if the connection is successful.
func Connect(dbType DbType, uri string) (*DbContext, error) {
	var dialector gorm.Dialector
	switch dbType {
	case DbPostgres:
		dialector = postgres.Open(uri)
	case DbMySQL:
		dialector = mysql.Open(uri)
	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}
	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, err
	}
	return &DbContext{
		db:  db,
		uri: uri,
	}, nil
}

// Disconnect from database. Currently used for Gorm.
func (c *DbContext) Disconnect() {
}

// Migrate database schema to match database models.
func (c *DbContext) Migrate() error {
	return c.db.AutoMigrate()
}

// Count number of records in a single table that satisfy provided condition.
func (c *DbContext) Count(model interface{}, query, args interface{}) (int64, error) {
	var count int64
	if query == nil {
		result := c.db.Model(model).Count(&count)
		return count, result.Error
	}
	result := c.db.Model(model).
		Where(query, args).
		Count(&count)
	return count, result.Error
}

// Truncate all tables.
func (c *DbContext) Reset() {
}

// Truncate specified table.
func (c *DbContext) Truncate(model interface{}) {
	c.db.Where("1 = 1").Delete(model)
}
