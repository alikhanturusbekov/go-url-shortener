package resetgen

import (
	"go/format"
	"os"
	"path/filepath"
)

// WriteFile formats generated code and writes it to reset.gen.go
// inside the target directory.
func WriteFile(dir string, code []byte) error {
	formatted, err := format.Source(code)
	if err != nil {
		return err
	}

	path := filepath.Join(dir, "reset.gen.go")
	return os.WriteFile(path, formatted, 0o644)
}
