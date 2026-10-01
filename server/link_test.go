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

package server

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinkRoutesBadRequests(t *testing.T) {
	s := New("", nil)

	_, body := doLinkRequest(t, s, http.MethodGet, "/links/not-a-uuid", nil)
	assert.Contains(t, string(body), "invalid link id")

	_, body = doLinkRequest(t, s, http.MethodPost, "/links", []byte(`{"url":""}`))
	assert.Contains(t, string(body), "error")

	_, body = doLinkRequest(t, s, http.MethodPut, "/links/not-a-uuid", []byte(`{"url":"https://x.io","title":"X"}`))
	assert.Contains(t, string(body), "invalid link id")

	_, body = doLinkRequest(t, s, http.MethodDelete, "/links/not-a-uuid", nil)
	assert.Contains(t, string(body), "invalid link id")
}

func doLinkRequest(t *testing.T, s *Server, method, path string, body []byte) (*http.Response, []byte) {
	t.Helper()
	var reader io.Reader
	if body != nil {
		reader = bytes.NewReader(body)
	}
	req := httptest.NewRequest(method, path, reader)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	s.Handler().ServeHTTP(rec, req)
	respBody, err := io.ReadAll(rec.Body)
	if err != nil {
		t.Fatalf("read response body: %v", err)
	}
	return rec.Result(), respBody
}
