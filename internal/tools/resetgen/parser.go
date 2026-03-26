package resetgen

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
)

// ParseFile parses a file and returns all structs marked with
// the `// generate:reset` directive.
func ParseFile(path string) ([]StructInfo, error) {
	fset := token.NewFileSet()

	file, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var result []StructInfo

	for _, decl := range file.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok || gen.Tok != token.TYPE {
			continue
		}

		if gen.Doc == nil {
			continue
		}

		hasResetDirective := false
		for _, comment := range gen.Doc.List {
			if comment.Text == "// generate:reset" {
				hasResetDirective = true
				break
			}
		}

		if !hasResetDirective {
			continue
		}

		for _, spec := range gen.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}

			info := StructInfo{
				Name:    typeSpec.Name.Name,
				Package: file.Name.Name,
				Dir:     filepath.Dir(path),
			}

			for _, field := range structType.Fields.List {
				for _, name := range field.Names {
					info.Fields = append(info.Fields, FieldInfo{
						Name: name.Name,
						Type: field.Type,
					})
				}
			}

			result = append(result, info)
		}
	}

	return result, nil
}
