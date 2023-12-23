# go-comment-eraser

Remove all comments from Go source code.

## SYNOPSIS

Install the command.

```console
go install github.com/tia-nn/go-comment-eraser@latest
```

```console
go-comment-eraser [-generated] $PATH_TO_SOURCE_DIRECTORY
```

## Examples

```console
go-comment-eraser github.com/tia-nn/go-comment-eraser
```

```diff
diff --git a/main.go b/main.go
index 5087d08..913debc 100644
--- a/main.go
+++ b/main.go
@@ -65,22 +65,18 @@ func eraseComment(src string) error {
 	return nil
 }

-// parseFile parses the Go source code file and returns the Go source
-// that is modified to erase all comments.
 func parseFile(src string) ([]byte, error) {
-	// Parse the Go source code file
+
 	fset := token.NewFileSet()
 	node, err := parser.ParseFile(fset, src, nil, parser.ParseComments)
 	if err != nil {
 		return nil, err
 	}

-	// If the AST is generated, just copy the file
 	if ast.IsGenerated(node) {
 		return nil, errSkipErase
 	}

-	// Erase all comments from the AST
 	ast.Inspect(node, func(n ast.Node) bool {
 		switch n := n.(type) {
 		case *ast.CommentGroup:
@@ -109,7 +105,6 @@ func parseFile(src string) ([]byte, error) {
 	}
 	node.Comments = list

-	// Print the modified Go source code
 	var buf bytes.Buffer
 	if err := format.Node(&buf, fset, node); err != nil {
 		return nil, err
```
