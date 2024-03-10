// Package osexitanalyzer provides an analyzer to check for direct calls to os.Exit() within the main.main() function.
// This analyzer is designed to be used with the golang.org/x/tools/go/analysis package.
//
// Usage:
//
//	import (
//		"golang.org/x/tools/go/analysis"
//		"github.com/example/osexitanalyzer"
//	)
//
//	var Analyzer = osexitanalyzer.OsExitAnalyzer
//
// The OsExitAnalyzer checks for direct calls to os.Exit() within the main.main() function.
// If such calls are found, it reports them as diagnostics.
//
// Example:
//
//	func main() {
//		os.Exit(1) // This call will be reported by the analyzer.
//	}
package osexitanalyzer

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// OsExitAnalyzer checks direct call os.Exit() in main.main() func
var OsExitAnalyzer = &analysis.Analyzer{
	Name: "osexitanalyzer",
	Doc:  "checks direct call os.Exit() in main.main() func",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		// loop by all declarations
		for _, d := range file.Decls {
			// skip not func declarations
			f, ok := d.(*ast.FuncDecl)
			if !ok {
				continue
			}
			// handle func declarations
			// in func main() package main
			if f.Name.Name == "main" && pass.Pkg.Name() == "main" {
				ast.Inspect(f.Body, func(x ast.Node) bool {
					if callExpr, ok := x.(*ast.CallExpr); ok {
						if ident, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
							// here is os.Exit() call
							pkg, ok := ident.X.(*ast.Ident)
							if ok {
								if ident.Sel.Name == "Exit" && pkg.Name == "os" {
									pos := pass.Fset.Position(x.Pos())
									pass.Reportf(x.Pos(), "direct call to os.Exit() found in func main() at %s:%d", pos.Filename, pos.Line)
								}
							}
						}
					}
					return true
				})
			}
		}
	}
	return nil, nil
}
