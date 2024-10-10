package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <filename.go>")
		return
	}

	filename := os.Args[1]

	// Parse the Go file
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, nil, parser.AllErrors)
	if err != nil {
		fmt.Printf("Failed to parse file: %v\n", err)
		return
	}

	// Walk through the AST and check variable names
	ast.Inspect(node, func(n ast.Node) bool {
		// Look for range statements (for loops)
		if rangeStmt, ok := n.(*ast.RangeStmt); ok {
			// Check if the key or value variables are named 'i'
			if ident, ok := rangeStmt.Key.(*ast.Ident); ok && ident.Name == "i" {
				// Exempt the variable 'i' in a range loop
				return true
			}
		}

		// Look for variable declarations
		if decl, ok := n.(*ast.AssignStmt); ok {
			for _, lhs := range decl.Lhs {
				if ident, ok := lhs.(*ast.Ident); ok {
					if len(ident.Name) <= 2 && ident.Name != "i" && ident.Name != "ok" {
						fmt.Printf("Variable '%s' is too short at position %d\n", ident.Name, fset.Position(ident.Pos()).Line)
					}
				}
			}
		}

		if decl, ok := n.(*ast.ValueSpec); ok {
			for _, name := range decl.Names {
				if len(name.Name) <= 2 && name.Name != "i" {
					fmt.Printf("Variable '%s' is too short at position %d\n", name.Name, fset.Position(name.Pos()).Line)
				}
			}
		}

		return true
	})
}
