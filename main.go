package main

import (
	"bytes"
	"errors"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

var errSkipErase = errors.New("skip erase")

func main() {
	eraseComment("main.go", "/dev/stdout")
}

func eraseComment(src, dst string) error {
	// Parse the Go source code file
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, src, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	// If the AST is generated, just copy the file
	if ast.IsGenerated(node) {
		return errSkipErase
	}

	// Erase all comments from the AST
	ast.Inspect(node, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.CommentGroup:
			list := n.List[:0]
			for _, c := range n.List {
				if isSpecialComment(c) {
					list = append(list, c)
				}
			}
			n.List = list
		}
		return true
	})
	list := node.Comments[:0]
	for _, g := range node.Comments {
		group := g.List[:0]
		for _, c := range g.List {
			if isSpecialComment(c) {
				group = append(group, c)
			}
		}
		if len(group) != 0 {
			g.List = group
			list = append(list, g)
		}
	}
	node.Comments = list

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

func isSpecialComment(c *ast.Comment) bool {
	return strings.HasPrefix(c.Text, "//go:")
}
