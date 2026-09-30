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
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Prepare tables required by link tests.
func migrateLinkTables(t *testing.T, c *DbContext) {
	t.Helper()
	require.NoError(t, c.db.AutoMigrate(&ID{}, &Link{}), "Failed to migrate tables")
}

// Populate the links table with fixture links and return them.
func seedLinks(t *testing.T, c *DbContext) []*Link {
	t.Helper()
	links := []*Link{
		{URL: "https://tforce.io/prism", Title: "TFprism"},
		{URL: "https://example.com/dev-null", Title: "Dev Null"},
	}
	require.NoError(t, c.InsertLinks(links), "Failed to insert links")
	return links
}

// Remove all records from tables used by link tests.
func truncateLinkTables(c *DbContext) {
	c.Truncate(&Link{})
	c.Truncate(&ID{})
}

func TestDbContextGetLink_Integration(t *testing.T) {
	c := connect(t)
	defer c.Disconnect()
	migrateLinkTables(t, c)
	seeded := seedLinks(t, c)
	defer truncateLinkTables(c)

	for _, in := range seeded {
		got, err := c.GetLink(in.Id)
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, *in, *got)
	}
	unknown, err := c.GetLink(uuid.New())
	require.NoError(t, err)
	assert.Nil(t, unknown)
}

func TestDbContextGetLinkByURL_Integration(t *testing.T) {
	c := connect(t)
	defer c.Disconnect()
	migrateLinkTables(t, c)
	seeded := seedLinks(t, c)
	defer truncateLinkTables(c)

	for _, in := range seeded {
		got, err := c.GetLinkByURL(in.URL)
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, *in, *got)
	}
	unknown, err := c.GetLinkByURL("https://example.com/unknown")
	require.NoError(t, err)
	assert.Nil(t, unknown)
}

func TestDbContextGetLinks_Integration(t *testing.T) {
	c := connect(t)
	defer c.Disconnect()
	migrateLinkTables(t, c)
	seeded := seedLinks(t, c)
	defer truncateLinkTables(c)

	ids := []uuid.UUID{seeded[0].Id, seeded[1].Id, uuid.New()}
	got, err := c.GetLinks(ids)
	require.NoError(t, err)
	assert.ElementsMatch(t, seeded, got)
}

func TestDbContextGetLinksByURL_Integration(t *testing.T) {
	c := connect(t)
	defer c.Disconnect()
	migrateLinkTables(t, c)
	seeded := seedLinks(t, c)
	defer truncateLinkTables(c)

	urls := []string{seeded[0].URL, "https://example.com/unknown"}
	got, err := c.GetLinksByURL(urls)
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, *seeded[0], *got[0])
}

func TestDbContextInsertLink_Integration(t *testing.T) {
	c := connect(t)
	defer c.Disconnect()
	migrateLinkTables(t, c)
	defer truncateLinkTables(c)

	in := &Link{URL: "https://tforce.io/prism", Title: "TFprism"}
	require.NoError(t, c.InsertLink(in))
	assert.NotEqual(t, uuid.Nil, in.Id)
	var got Link
	require.NoError(t, c.db.Where("id = ?", in.Id).First(&got).Error, "Failed to query inserted link")
	assert.Equal(t, *in, got)
}

func TestDbContextInsertLinks_Integration(t *testing.T) {
	c := connect(t)
	defer c.Disconnect()
	migrateLinkTables(t, c)
	defer truncateLinkTables(c)

	inserted := []*Link{
		{URL: "https://tforce.io/prism", Title: "TFprism"},
		{URL: "https://example.com/dev-null", Title: "Dev Null"},
	}
	require.NoError(t, c.InsertLinks(inserted))
	for _, in := range inserted {
		assert.NotEqual(t, uuid.Nil, in.Id)
		var got Link
		require.NoError(t, c.db.Where("id = ?", in.Id).First(&got).Error, "Failed to query inserted link")
		assert.Equal(t, *in, got)
	}
}

func TestDbContextUpdateLink_Integration(t *testing.T) {
	c := connect(t)
	defer c.Disconnect()
	migrateLinkTables(t, c)
	seeded := seedLinks(t, c)
	defer truncateLinkTables(c)

	in := seeded[0]
	in.URL = "https://tforce.io/prism-v2"
	in.Title = "TFprism v2"
	require.NoError(t, c.UpdateLink(in.Id, in))
	got, err := c.GetLink(in.Id)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, *in, *got)

	require.NoError(t, c.UpdateLink(uuid.New(), in))
	var count int64
	require.NoError(t, c.db.Model(&Link{}).Count(&count).Error)
	assert.Equal(t, int64(2), count)
}

func TestDbContextDeleteLink_Integration(t *testing.T) {
	c := connect(t)
	defer c.Disconnect()
	migrateLinkTables(t, c)
	seeded := seedLinks(t, c)
	defer truncateLinkTables(c)

	require.NoError(t, c.DeleteLink(seeded[0].Id))
	got, err := c.GetLink(seeded[0].Id)
	require.NoError(t, err)
	assert.Nil(t, got)

	require.NoError(t, c.DeleteLink(uuid.New()))
}
