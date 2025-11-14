package integration

import (
"testing"
cli "github.com/anabellayholman/homevision_unpacker/cli"
)

func TestIntegrationSampleEnv(t *testing.T) {
files, err := cli.ParseEnv("../cli/sample.env")
if err != nil { t.Fatalf("parse error: %v", err) }
if len(files) == 0 { t.Fatal("no files extracted") }
}
