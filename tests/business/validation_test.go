package business

import (
	"crypto/sha1"
	"fmt"
	"testing"
)

type File struct {
	Name string
	Size int
	Hash string
	Data []byte
}

func verifySHA1(f File) bool {
	if f.Hash == "" {
		return true
	}
	sum := sha1.Sum(f.Data)
	return fmt.Sprintf("%x", sum) == f.Hash
}

func TestBusinessRulesOnSample(t *testing.T) {
	files := []File{
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

		if f.Hash != "" && !verifySHA1(f) {
			t.Fatalf("sha1 mismatch %s", f.Name)
		}

		if len(f.Data) != f.Size {
			t.Fatalf("size mismatch %s: expected %d, got %d", f.Name, f.Size, len(f.Data))
		}
	}
}
