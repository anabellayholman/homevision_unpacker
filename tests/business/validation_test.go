package business

import (
	"testing"

	cli "github.com/anabellayholman/homevision_unpacker/pkg/cli"
)

func TestBusinessRulesOnSample(t *testing.T) {
	files := []cli.File{
		{
			Name: "image1.jpg",
			Size: 5,
			Hash: "",
			Data: []byte{1, 2, 3, 4, 5},
		},
		{
			Name: "config.yaml",
			Size: 4,
			Hash: "",
			Data: []byte("test"),
		},
	}

	names := make(map[string]bool)

	for i, f := range files {
		if f.Name == "" {
			t.Fatalf("empty name at idx %d", i)
		}
		if names[f.Name] {
			t.Fatalf("duplicate name %s", f.Name)
		}
		names[f.Name] = true

		if f.Size == 0 {
			t.Fatalf("zero size for %s", f.Name)
		}

		if f.Hash != "" && !cli.VerifySHA1(f) {
			t.Fatalf("sha1 mismatch %s", f.Name)
		}

		if len(f.Data) != f.Size {
			t.Fatalf("size mismatch %s: expected %d, got %d", f.Name, f.Size, len(f.Data))
		}
	}
}
