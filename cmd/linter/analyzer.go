package main

import (
	"go/ast"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "nopanic",
	Doc:  "checks forbidden panic, log.Fatal and os.Exit usage",
	Run:  run,
}

// run runs the analyzer for forbidden panic, log.Fatal and os.Exit usage
func run(pass *analysis.Pass) (interface{}, error) {
	pkgIsMain := pass.Pkg.Name() == "main"

	for _, file := range pass.Files {
		var stack []ast.Node

		ast.Inspect(file, func(n ast.Node) bool {
			if n == nil {
				stack = stack[:len(stack)-1]
				return false
			}

			stack = append(stack, n)

			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			// panic()
			if ident, ok := call.Fun.(*ast.Ident); ok && ident.Name == "panic" {
				pass.Reportf(call.Pos(), "panic is forbidden")
				return true
			}

			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}

			obj := pass.TypesInfo.Uses[sel.Sel]
			fn, ok := obj.(*types.Func)
			if !ok {
				return true
			}

			pkg := fn.Pkg()
			if pkg == nil {
				return true
			}

			if pkg.Path() == "log" && fn.Name() == "Fatal" && !isAllowedExit(pass, stack, pkgIsMain, call) {
				pass.Reportf(call.Pos(), "log.Fatal is forbidden outside main")
			}

			if pkg.Path() == "os" && fn.Name() == "Exit" && !isAllowedExit(pass, stack, pkgIsMain, call) {
				pass.Reportf(call.Pos(), "os.Exit is forbidden outside main")
			}

			return true
		})
	}

	return nil, nil
}

// isAllowedExit checks where the function call is allowed
func isAllowedExit(pass *analysis.Pass, stack []ast.Node, pkgIsMain bool, node ast.Node) bool {
	if pkgIsMain {
		for i := len(stack) - 1; i >= 0; i-- {
			if fn, ok := stack[i].(*ast.FuncDecl); ok {
				return fn.Name.Name == "main"
			}
		}
	}

	// allow anywhere in *_test.go
	if isTestFile(pass, node) {
		return true
	}

	return false
}

// isTestFile checks whether the code file is *_test.go
func isTestFile(pass *analysis.Pass, node ast.Node) bool {
	pos := pass.Fset.Position(node.Pos())
	return strings.HasSuffix(pos.Filename, "_test.go")
}
