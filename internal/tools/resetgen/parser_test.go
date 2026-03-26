package resetgen

import (
	"path/filepath"
	"testing"
)

// TestParseFile_ResetableStruct verifies parsing of ResetableStruct
func TestParseFile_ResetableStruct(t *testing.T) {
	path := filepath.Join("testdata", "example.go")

	structs, err := ParseFile(path)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}

	if len(structs) != 1 {
		t.Fatalf("expected 1 struct, got %d", len(structs))
	}

	s := structs[0]

	if s.Name != "ResetableStruct" {
		t.Fatalf("expected struct ResetableStruct, got %s", s.Name)
	}

	if len(s.Fields) != 6 {
		t.Fatalf("expected 6 fields, got %d", len(s.Fields))
	}
}
