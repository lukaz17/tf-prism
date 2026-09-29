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
	"testing"
	"time"
)

func TestIDTableName(t *testing.T) {
	if (ID{}).TableName() != "ids" {
		t.Errorf("Wrong ID table name. Expected 'ids' Actual '%s'", (ID{}).TableName())
	}
}

func TestNewRandomID(t *testing.T) {
	firstSeen := map[[16]byte]bool{}
	for i := 0; i < 10; i++ {
		id := NewRandomID(RefTypeLink)
		if id.Guid.Version() != 4 {
			t.Errorf("Wrong GUID version. Expected 4 Actual %d", id.Guid.Version())
		}
		if id.Ref != RefTypeLink {
			t.Errorf("Wrong id ref. Expected %d Actual %d", RefTypeLink, id.Ref)
		}
		if (time.Time(id.CreatedAt)).IsZero() {
			t.Error("Wrong id created_at. Expected non-zero time Actual zero")
		}
		if firstSeen[id.Guid] {
			t.Errorf("Duplicated GUID after %d iterations: %s", i, id.Guid)
			break
		}
		firstSeen[id.Guid] = true
	}
}

func TestNewTimeAwareID(t *testing.T) {
	before := time.Now().Add(-time.Second).Round(time.Second)
	id := NewTimeAwareID(RefTypeLink)
	if id.Guid.Version() != 7 {
		t.Errorf("Wrong GUID version. Expected 7 Actual %d", id.Guid.Version())
	}
	after := time.Now().Add(time.Second)
	if id.Ref != RefTypeLink {
		t.Errorf("Wrong id ref. Expected %d Actual %d", RefTypeLink, id.Ref)
	}
	sec, nsec := id.Guid.Time().UnixTime()
	uuidTime := time.Unix(sec, nsec)
	if uuidTime.Before(before) || uuidTime.After(after) {
		t.Errorf("Wrong time-aware GUID timestamp. Expected between %s and %s Actual %s", before, after, uuidTime)
	}
	if (time.Time(id.CreatedAt)).IsZero() {
		t.Error("Wrong id created_at. Expected non-zero time Actual zero")
	}
}
