package integration

import "testing"

type File struct {
	Name string
	Size int
	Hash string
	Data []byte
}

func TestIntegrationSampleEnv(t *testing.T) {
	files := []File{
		{
			Name: "image1.jpg",
			Size: 5,
			Data: []byte{1, 2, 3, 4, 5},
		},
		{
			Name: "config.yaml",
			Size: 4,
			Data: []byte("test"),
		},
	}

	if len(files) == 0 {
		t.Fatal("no files extracted")
	}
}
