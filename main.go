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

	returnCode := 0

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

		// Handle range statements (for ... range loops)
		if rangeStmt, ok := n.(*ast.RangeStmt); ok {
			// Check the key variable in "for k, v := range ..."
			if key, ok := rangeStmt.Key.(*ast.Ident); ok {
				if !isException(key.Name) && len(key.Name) <= 2 {
					fmt.Printf("Variable '%s' in range loop is too short at position %d\n", key.Name, fset.Position(key.Pos()).Line)
					returnCode = 1
				}
			}

			// Check the value variable in "for k, v := range ..."
			if value, ok := rangeStmt.Value.(*ast.Ident); ok {
				if !isException(value.Name) && len(value.Name) <= 2 {
					fmt.Printf("Variable '%s' in range loop is too short at position %d\n", value.Name, fset.Position(value.Pos()).Line)
					returnCode = 1
				}
			}

			// Skip further checks for range loops as we only care about initialization.
			return true
		}

		// Handle standard for-loops (initialization part with ":=")
		if forStmt, ok := n.(*ast.ForStmt); ok {
			// Check if initialization is using ":="
			if initStmt, ok := forStmt.Init.(*ast.AssignStmt); ok && initStmt.Tok == token.DEFINE {
				for _, lhs := range initStmt.Lhs {
					if ident, ok := lhs.(*ast.Ident); ok {
						// If the variable name is 'i', it's fine to ignore
						if isException(ident.Name) {
							return true
						}
						// Check if the variable is too short
						if len(ident.Name) <= 2 {
							fmt.Printf("Variable '%s' in for-loop initialization is too short at position %d\n", ident.Name, fset.Position(ident.Pos()).Line)
							returnCode = 1
						}
					}
				}
			}
		}

		// Handle variable declarations outside of loops (":=" assignments only)
		if decl, ok := n.(*ast.AssignStmt); ok {
			// Ensure we only check ":=" (initializations) and not "=" (assignments)
			if decl.Tok == token.DEFINE {
				for _, lhs := range decl.Lhs {
					if ident, ok := lhs.(*ast.Ident); ok {
						// Check if the variable name is too short and not the exempted 'i'
						if len(ident.Name) <= 2 && !isException(ident.Name) {
							fmt.Printf("Variable '%s' is too short at position %d\n", ident.Name, fset.Position(ident.Pos()).Line)
							returnCode = 1
						}
					}
				}
			}
		}

		// Handle regular variable declarations
		if decl, ok := n.(*ast.ValueSpec); ok {
			for _, name := range decl.Names {
				// Check if the variable name is too short and not the exempted 'i'
				if len(name.Name) <= 2 && !isException(name.Name) {
					fmt.Printf("Variable '%s' is too short at position %d\n", name.Name, fset.Position(name.Pos()).Line)
					returnCode = 1
				}
			}
		}

		return true
	})
	os.Exit(returnCode)
}

func isException(varName string) bool {
	exceptions := []string{"ok", "i", "_", "tx", "wg"}
	for _, exception := range exceptions {
		if varName == exception {
			return true
		}
	}
	return false
}
