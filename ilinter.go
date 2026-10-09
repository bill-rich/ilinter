// Package ilinter defines an analyzer that reports variables and constants
// whose names are shorter than a minimum length.
//
// Every variable and constant declared in a package is checked, whether at
// package level or inside a function: var and const declarations, short
// variable declarations, range and select clause variables, and type switch
// guards. Function parameters, results, receivers and struct fields are not
// checked. The blank identifier is never reported, and a configurable list of
// names such as "ok" is always allowed.
package ilinter

import (
	"go/ast"
	"go/types"
	"strings"
	"unicode/utf8"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer reports variable and constant names shorter than -min characters.
var Analyzer = &analysis.Analyzer{
	Name:     "ilinter",
	Doc:      "report variable and constant names shorter than a minimum length",
	URL:      "https://github.com/bill-rich/ilinter",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

var (
	minLength = 3
	allow     = "ok,i,tx,wg"
)

func init() {
	Analyzer.Flags.IntVar(&minLength, "min", minLength, "minimum number of characters in a name")
	Analyzer.Flags.StringVar(&allow, "allow", allow, "comma-separated names that are always allowed")
}

func run(pass *analysis.Pass) (any, error) {
	allowed := map[string]bool{"_": true}
	for _, name := range strings.Split(allow, ",") {
		if name = strings.TrimSpace(name); name != "" {
			allowed[name] = true
		}
	}

	report := func(ident *ast.Ident, kind string) {
		if allowed[ident.Name] || utf8.RuneCountInString(ident.Name) >= minLength {
			return
		}
		pass.Reportf(ident.Pos(), "%s name %q is shorter than %d characters", kind, ident.Name, minLength)
	}

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	filter := []ast.Node{(*ast.File)(nil), (*ast.Ident)(nil), (*ast.TypeSwitchStmt)(nil)}
	insp.Nodes(filter, func(node ast.Node, push bool) bool {
		if !push {
			return true
		}
		switch node := node.(type) {
		case *ast.File:
			// Generated code, such as the test main that go test writes, is
			// not ours to name. Returning false skips the whole file.
			return !ast.IsGenerated(node)
		case *ast.Ident:
			switch obj := pass.TypesInfo.Defs[node].(type) {
			case *types.Var:
				if kind := obj.Kind(); kind == types.LocalVar || kind == types.PackageVar {
					report(node, "variable")
				}
			case *types.Const:
				report(node, "constant")
			}
		case *ast.TypeSwitchStmt:
			// The guard of a type switch is a symbolic variable: each clause
			// gets its own implicit object and Defs maps the identifier to nil.
			if assign, ok := node.Assign.(*ast.AssignStmt); ok {
				if ident, ok := assign.Lhs[0].(*ast.Ident); ok {
					report(ident, "variable")
				}
			}
		}
		return true
	})
	return nil, nil
}
