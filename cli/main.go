package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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

	if strings.TrimSpace(outdir) == "" {
		fmt.Println("output directory cannot be empty")
		os.Exit(1)
	}

	if err := os.MkdirAll(outdir, 0755); err != nil {
		fmt.Println("cannot create output directory:", err)
		os.Exit(1)
	}

	fmt.Printf("processing %d files from: %s\n", len(files), input)

	for _, f := range files {
		name := filepath.Clean(f.Name)

		if filepath.Ext(name) == "" && f.Ext != "" {
			name = name + "." + strings.TrimPrefix(f.Ext, ".")
		}

		outpath := filepath.Join(outdir, name)

		if err := os.MkdirAll(filepath.Dir(outpath), 0755); err != nil {
			fmt.Println("cannot create directory for file:", err)
			os.Exit(1)
		}

		if err := os.WriteFile(outpath, f.Data, 0644); err != nil {
			fmt.Println("write error:", err)
			os.Exit(1)
		}

		ok := cli.VerifySHA1(f)
		fmt.Printf("- %s (size=%d, ok=%v)\n", name, f.Size, ok)
	}

	fmt.Printf("extracted %d files to %s\n", len(files), outdir)
}
