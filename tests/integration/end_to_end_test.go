package integration

import (
	"path/filepath"
	"testing"

	cli "github.com/anabellayholman/homevision_unpacker/pkg/cli"
)

func TestIntegrationSampleEnv(t *testing.T) {
	envPath := filepath.Join("tests", "sample.env")
	files, err := cli.ParseEnv(envPath)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if len(files) == 0 {
		t.Fatal("no files extracted")
	}
}
