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

package server

import (
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/tforceaio/tf-prism/db"
)

// Integration test: full CRUD flow over the HTTP API, requires real database.
func TestLinkRoutesCRUD_Integration(t *testing.T) {
	ctx := connectLink(t)
	defer ctx.Disconnect()
	s := New("", ctx)
	defer ctx.Truncate(&db.Link{})

	// List initially empty.
	resp, body := doLinkRequest(t, s, http.MethodGet, "/links", nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "[]", string(body))

	// Create.
	created := &db.Link{}
	payload := `{"url":"https://tforce.io/prism-api","title":"TFprism API"}`
	resp, body = doLinkRequest(t, s, http.MethodPost, "/links", []byte(payload))
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	if err := json.Unmarshal(body, created); err != nil {
		t.Fatalf("parse created link: %v", err)
	}
	assert.NotEqual(t, uuid.Nil, created.Id)

	// Get by ID.
	resp, body = doLinkRequest(t, s, http.MethodGet, "/links/"+created.Id.String(), nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Contains(t, string(body), "https://tforce.io/prism-api")

	// List contains the created link.
	resp, body = doLinkRequest(t, s, http.MethodGet, "/links", nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Contains(t, string(body), "https://tforce.io/prism-api")

	// Update.
	payload = `{"url":"https://tforce.io/prism-api-v2","title":"TFprism API v2"}`
	resp, _ = doLinkRequest(t, s, http.MethodPut, "/links/"+created.Id.String(), []byte(payload))
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	resp, body = doLinkRequest(t, s, http.MethodGet, "/links/"+created.Id.String(), nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Contains(t, string(body), "prism-api-v2")

	// Delete.
	resp, _ = doLinkRequest(t, s, http.MethodDelete, "/links/"+created.Id.String(), nil)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	resp, body = doLinkRequest(t, s, http.MethodGet, "/links/"+created.Id.String(), nil)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	assert.Contains(t, string(body), "not found")
}

// Return new DbContext from TFPRISM_TEST_DB_TYPE and TFPRISM_TEST_DB_URI.
// Skip the test when either variable is not set.
func connectLink(t *testing.T) *db.DbContext {
	t.Helper()
	dbType := db.DbType(os.Getenv("TFPRISM_TEST_DB_TYPE"))
	uri := os.Getenv("TFPRISM_TEST_DB_URI")
	if dbType == "" || uri == "" {
		t.Skip("TFPRISM_TEST_DB_TYPE or TFPRISM_TEST_DB_URI not set; skipping db integration test")
	}
	c, err := db.Connect(dbType, uri)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	return c
}
