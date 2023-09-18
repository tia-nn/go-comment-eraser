package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var errSkipErase = errors.New("skip erase")

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
	}

	dir, err := filepath.Abs(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			if strings.HasPrefix(info.Name(), ".") || info.Name() == "testdata" {
				return filepath.SkipDir
			}
			return nil
		}

		if strings.HasSuffix(info.Name(), ".go") {
			if err := eraseComment(path); err != nil {
				log.Println(path, ":", err)
			}
		}
		return nil
	})
}

func eraseComment(src string) error {
	data, err := parseFile(src)
	if errors.Is(err, errSkipErase) {
		return nil
	}
	if err != nil {
		return err
	}

	tmp := fmt.Sprintf(src+".tmp.%d", os.Getpid())
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return err
	}
	defer os.ReadFile(tmp)
	if err := os.Rename(tmp, src); err != nil {
		return err
	}
	return nil
}

// parseFile parses the Go source code file and returns the Go source
// that is modified to erase all comments.
func parseFile(src string) ([]byte, error) {
	// Parse the Go source code file
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, src, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	// If the AST is generated, just copy the file
	if ast.IsGenerated(node) {
		return nil, errSkipErase
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
		return nil, err
	}

	return buf.Bytes(), nil
}

func isSpecialComment(c *ast.Comment) bool {
	return strings.HasPrefix(c.Text, "//go:")
}
