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

package engine

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateFileExists(t *testing.T) {
	existing := filepath.Join(t.TempDir(), "a.txt")
	if err := os.WriteFile(existing, []byte("x"), 0664); err != nil {
		t.Fatal(err)
	}

	if err := ValidateFileExists(existing); err != nil {
		t.Errorf("expected nil, got %v", err)
	}
	if err := ValidateFileExists(filepath.Join(t.TempDir(), "missing.txt")); err == nil {
		t.Error("expected error for missing path, got nil")
	}
}

func TestValidateString(t *testing.T) {
	if err := ValidateString("value", "label"); err != nil {
		t.Errorf("expected nil, got %v", err)
	}
	if err := ValidateString("", "label"); err == nil {
		t.Error("expected error for empty value, got nil")
	}
}
