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

package db

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DbType selects which database server to connect to.
type DbType string

// RefType is defined using CRC-32/ISO-HDLC result of the type name in lowercase.
type RefType uint32

// Return underlying database type for GORM.
func (RefType) GormDataType() string {
	return "int unsigned"
}

// Timestamp is a timestamp stored in database as milliseconds since Unix epoch.
type Timestamp time.Time

// Return underlying database type for GORM.
func (Timestamp) GormDataType() string {
	return "bigint"
}

// Return new Timestamp instance, truncated to millisecond precision.
func NewTimestamp() Timestamp {
	return Timestamp(time.Now().UTC().Truncate(time.Millisecond))
}

// Return new Timestamp instance using given time `t`, truncated to millisecond precision.
func NewTimestampFromTime(t time.Time) Timestamp {
	return Timestamp(t.Truncate(time.Millisecond))
}

// Decode database value into Timestamp.
func (c *Timestamp) Scan(src interface{}) error {
	switch v := src.(type) {
	case nil:
		*c = Timestamp(time.Time{})
	case int64:
		*c = Timestamp(time.UnixMilli(v))
	case time.Time:
		*c = Timestamp(v)
	case []byte:
		return fmt.Errorf("cannot scan []byte into Timestamp: %s", v)
	default:
		return fmt.Errorf("cannot scan %T into Timestamp", src)
	}
	return nil
}

// Return value for database.
func (c Timestamp) Value() (driver.Value, error) {
	return int64(time.Time(c).Truncate(time.Millisecond).UnixMilli()), nil
}

// Determine if the error is record not found.
func (c *DbContext) isEmptyResultError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return errStr == "record not found"
}

// Return new ID using UUID version 4, and include it to the transaction.
func (c *DbContext) newRandomID(tx *gorm.DB, typ RefType) uuid.UUID {
	id := NewRandomID(typ)
	tx.Create(id)
	return id.Guid
}

// Return new ID using UUID version 7, and include it to the transaction.
func (c *DbContext) newTimeAwareID(tx *gorm.DB, typ RefType) uuid.UUID {
	id := NewTimeAwareID(typ)
	tx.Create(id)
	return id.Guid
}
