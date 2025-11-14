package business

import (
"testing"
cli "github.com/anabellayholman/homevision_unpacker/cli"
)

func TestBusinessRulesOnSample(t *testing.T) {
files, err := cli.ParseEnv("../cli/sample.env")
if err != nil { t.Fatalf("parse error: %v", err) }
names := make(map[string]bool)
for i, f := range files {
if f.Name == "" { t.Fatalf("empty name at idx %d", i) }
if names[f.Name] { t.Fatalf("duplicate name %s", f.Name) }
names[f.Name] = true
if f.Size == 0 { t.Fatalf("zero size for %s", f.Name) }
if f.Hash != "" && !cli.VerifySHA1(f) { t.Fatalf("sha1 mismatch %s", f.Name) }
if len(f.Data) != f.Size { t.Fatalf("size mismatch %s", f.Name) }
}
}
