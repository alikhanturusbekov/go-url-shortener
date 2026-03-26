package resetgen

import (
	"go/ast"
	"strings"
	"testing"
)

// TestGenerate_ResetableStruct verifies Reset method generation
func TestGenerate_ResetableStruct(t *testing.T) {
	structs := []StructInfo{
		{
			Name: "ResetableStruct",
			Fields: []FieldInfo{
				{"i", &ast.Ident{Name: "int"}},
				{"str", &ast.Ident{Name: "string"}},
				{"strP", &ast.StarExpr{X: &ast.Ident{Name: "string"}}},
				{"s", &ast.ArrayType{}},
				{"m", &ast.MapType{}},
				{"child", &ast.StarExpr{X: &ast.Ident{Name: "ResetableStruct"}}},
			},
		},
	}

	code := string(Generate("testdata", structs))

	tests := []string{
		"r.i = 0",
		`r.str = ""`,
		"if r.strP != nil",
		"r.s = r.s[:0]",
		"clear(r.m)",
		"resetter.Reset()",
	}

	for _, expected := range tests {
		if !strings.Contains(code, expected) {
			t.Fatalf("generated code missing: %s\nCode:\n%s", expected, code)
		}
	}
}
