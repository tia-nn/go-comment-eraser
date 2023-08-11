package main

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"os"
	"strings"
)

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
		return copyFile(src, dst)
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

func copyFile(src, dst string) error {
	r, err := os.Open(src)
	if err != nil {
		return err
	}
	defer r.Close()

	stat, err := r.Stat()
	if err != nil {
		return err
	}

	w, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, stat.Mode())
	if err != nil {
		return err
	}
	_, err = io.Copy(w, r)
	if err1 := w.Close(); err1 != nil && err == nil {
		err = err1
	}
	return err
}
