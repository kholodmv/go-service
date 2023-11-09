// Package exitchecker Module exit checker set check for calling os.Exit()
// via inspecting full AST
package exitchecker

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "exitchecker",
	Doc:  "check for os.Exit() in main",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	check := func(node ast.Node) {
		if expr, ok := node.(*ast.ExprStmt); ok {
			if c, ok := expr.X.(*ast.CallExpr); ok {
				if s, ok := c.Fun.(*ast.SelectorExpr); ok {
					if i, ok := s.X.(*ast.Ident); ok {
						// only for calling the os.Exit() function
						if i.Name == "os" && s.Sel.Name == "Exit" {
							pass.Reportf(s.Pos(), "not allowed using of os.Exit()")
						}
					}
				}
			}
		}
	}

	for _, file := range pass.Files {
		if file.Name.Name == "main" {
			// using ast.Inspect function we go through all AST nodes
			ast.Inspect(file, func(node ast.Node) bool {
				if m, ok := node.(*ast.FuncDecl); ok {
					if m.Name.Name == "main" {
						for _, st := range m.Body.List {
							check(st)
						}
					}
				}
				return true
			})
		}
	}
	return nil, nil
}
