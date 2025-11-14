package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/anabellayholman/homevision_unpacker/pkg/cli"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: unpack <file.env> [outdir]")
		os.Exit(1)
	}

	input := os.Args[1]
	files, err := cli.ParseEnv(input)
	if err != nil {
		fmt.Println("read/parse error:", err)
		os.Exit(1)
	}

	outdir := "output"
	if len(os.Args) >= 3 {
		outdir = os.Args[2]
	}

	os.MkdirAll(outdir, 0755)
	err = cli.WriteOutputs(files, outdir)
	if err != nil {
		fmt.Println("write error:", err)
		os.Exit(1)
	}

	fmt.Printf("extracted %d files to %s\n", len(files), outdir)
	for _, f := range files {
		ok := cli.VerifySHA1(f)
		fmt.Printf("- %s (size=%d, ok=%v)\n", f.Name, f.Size, ok)
	}
}
