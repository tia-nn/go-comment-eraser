package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRemoveComment(t *testing.T) {
	tmpdir := t.TempDir()
	src := filepath.Join(tmpdir, "src.go")
	dst := filepath.Join(tmpdir, "dst.go")
	os.WriteFile(src, []byte(`// This is a comment.
package main

func main() {
	// This is another comment.
}
`), 0644)

	if err := removeComment(src, dst); err != nil {
		t.Error(err)
	}

	got, err := os.ReadFile(dst)
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
