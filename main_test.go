package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEraseComment(t *testing.T) {
	tmpdir := t.TempDir()
	src := filepath.Join(tmpdir, "src.go")
	os.WriteFile(src, []byte(`// This is a comment.
package main

func main() {
	// This is another comment.
}
`), 0644)

	if err := eraseComment(src); err != nil {
		t.Error(err)
	}

	got, err := os.ReadFile(src)
	if err != nil {
		t.Error(err)
	}
	want := `package main

func main() {

}
`
	if string(got) != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestEraseComment_LeaveSpecialComments(t *testing.T) {
	tmpdir := t.TempDir()
	src := filepath.Join(tmpdir, "src.go")
	os.WriteFile(src, []byte(`//go:build something

//go:generate bar

package main

//go:inline
func main() {
}
`), 0644)

	if err := eraseComment(src); err != nil {
		t.Error(err)
	}

	got, err := os.ReadFile(src)
	if err != nil {
		t.Error(err)
	}
	want := `//go:build something

//go:generate bar

package main

//go:inline
func main() {
}
`
	if string(got) != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestEraseComment_GeneratedFile(t *testing.T) {
	tmpdir := t.TempDir()
	src := filepath.Join(tmpdir, "src.go")
	os.WriteFile(src, []byte(`// Code generated by foo. DO NOT EDIT.

package main

func main() {
	// the eraser does nothing in the automatically generated file.
}
`), 0644)

	if err := eraseComment(src); err != nil {
		t.Error(err)
	}
}

func TestEraseComment_IncludeGeneratedFileFlag(t *testing.T) {
	tmpdir := t.TempDir()
	src := filepath.Join(tmpdir, "src.go")
	os.WriteFile(src, []byte(`// Code generated by baz. DO NOT EDIT.

package main

func main() {
	// this comment is erased if generated flag is true.
}
`), 0644)

	generatedFlag = true

	if err := eraseComment(src); err != nil {
		t.Error(err)
	}

	got, err := os.ReadFile(src)
	if err != nil {
		t.Error(err)
	}
	want := `package main

func main() {

}
`
	if string(got) != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
