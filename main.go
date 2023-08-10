package main

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
)

func main() {
	removeComment("main.go", "/dev/stdout")
}

func removeComment(src, dst string) error {
	// Parse the Go source code file
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, src, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	// Remove all comments from the AST
	ast.Inspect(node, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.CommentGroup:
			n.List = nil
		}
		return true
	})
	node.Comments = nil

	// Print the modified Go source code
	var buf bytes.Buffer
	if err := format.Node(&buf, fset, node); err != nil {
		return err
	}

	// Write the modified Go source code to file
	if err := os.WriteFile(dst, buf.Bytes(), 0644); err != nil {
		return err
	}
	return nil
}
