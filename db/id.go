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
	"time"

	"github.com/google/uuid"
)

// ID represents a 128-bit unique identifier.
type ID struct {
	Guid      uuid.UUID `gorm:"column:guid;primaryKey;size:36"`
	Ref       RefType   `gorm:"column:ref"`
	CreatedAt Timestamp `gorm:"column:created_at"`
}

// Return table name for GORM.
func (ID) TableName() string { return "ids" }

// Return new ID using UUID version 4.
func NewRandomID(typ RefType) *ID {
	newUUID, _ := uuid.NewRandom()
	return &ID{
		Guid:      newUUID,
		Ref:       typ,
		CreatedAt: Timestamp(time.Now()),
	}
}

// Return new ID using UUID version 7.
func NewTimeAwareID(typ RefType) *ID {
	newUUID, _ := uuid.NewV7()
	return &ID{
		Guid:      newUUID,
		Ref:       typ,
		CreatedAt: Timestamp(time.Now()),
	}
}
