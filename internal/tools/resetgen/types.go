package resetgen

import "go/ast"

// StructInfo describes a reset method generation
type StructInfo struct {
	Package string
	Name    string
	Fields  []FieldInfo
	Dir     string
}

// FieldInfo describes a field and its AST type
type FieldInfo struct {
	Name string
	Type ast.Expr
}
