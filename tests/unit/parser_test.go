package unit

import (
	"testing"

	cli "github.com/anabellayholman/homevision_unpacker/pkg/cli"
)

func TestParseEnvBytesBasic(t *testing.T) {
	data := []byte("FILENAME/a.txt\nEXT/txt\nSHA1/\n" + string([]byte{0x41, 0x42, 0x43}))
	files, err := cli.ParseEnvBytes(data)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if len(files) != 1 {
		t.Fatalf("expected 1 got %d", len(files))
	}
	if files[0].Ext != "txt" {
		t.Fatalf("ext expected txt")
	}
}
