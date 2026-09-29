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

//go:build db

package db

import (
	"os"
	"testing"
)

// Return new DbContext from TFPRISM_TEST_DB_TYPE and TFPRISM_TEST_DB_URI.
// Skip the test when either variable is not set.
func connect(t *testing.T) *DbContext {
	t.Helper()
	dbType := DbType(os.Getenv("TFPRISM_TEST_DB_TYPE"))
	uri := os.Getenv("TFPRISM_TEST_DB_URI")
	if dbType == "" || uri == "" {
		t.Skip("TFPRISM_TEST_DB_TYPE or TFPRISM_TEST_DB_URI not set; skipping db integration test")
	}
	c, err := Connect(dbType, uri)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	return c
}

func TestDbContextConnect_Integration(t *testing.T) {
	c := connect(t)
	defer c.Disconnect()
	if c == nil {
		t.Fatal("Nil DbContext after successful connect")
	}
}

// Count and Truncate are method of DbContext. InsertLinks is tested here
// because tests of these methods require a populated table.
func TestDbContextMethods_Integration(t *testing.T) {
	c := connect(t)
	defer c.Disconnect()

	// Prepare tables expected by the methods under test.
	if err := c.db.AutoMigrate(&ID{}, &Link{}); err != nil {
		t.Fatalf("Failed to migrate tables: %v", err)
	}

	links := []*Link{
		{URL: "https://tforce.io/prism", Title: "TFprism"},
		{URL: "https://example.com/dev-null", Title: "Dev Null"},
	}
	if err := c.InsertLinks(links); err != nil {
		t.Fatalf("Failed to insert links: %v", err)
	}

	count, err := c.Count(&Link{}, nil, nil)
	if err != nil {
		t.Fatalf("Failed to count links: %v", err)
	}
	if count != 2 {
		t.Errorf("Wrong link count. Expected 2 Actual %d", count)
	}

	count, err = c.Count(&Link{}, "url = ?", "https://tforce.io/prism")
	if err != nil {
		t.Fatalf("Failed to count links with condition: %v", err)
	}
	if count != 1 {
		t.Errorf("Wrong filtered link count. Expected 1 Actual %d", count)
	}

	c.Truncate(&Link{})
	count, err = c.Count(&Link{}, nil, nil)
	if err != nil {
		t.Fatalf("Failed to count links after truncate: %v", err)
	}
	if count != 0 {
		t.Errorf("Wrong link count after truncate. Expected 0 Actual %d", count)
	}
}
