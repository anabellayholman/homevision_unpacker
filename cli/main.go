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

	// Crear el directorio de salida
	err = os.MkdirAll(outdir, 0755)
	if err != nil {
		fmt.Println("cannot create output directory:", err)
		os.Exit(1)
	}

	fmt.Printf("processing %d files from: %s\n", len(files), input)

	// Escribir los archivos extraídos usando filepath
	for _, f := range files {
		safe := filepath.Clean(f.Name)
		if filepath.Ext(safe) == "" && f.Ext != "" {
			safe = safe + "." + strings.TrimPrefix(f.Ext, ".")
		}
		outpath := filepath.Join(outdir, safe)

		err := os.MkdirAll(filepath.Dir(outpath), 0755) // crea subcarpetas si es necesario
		if err != nil {
			fmt.Println("cannot create directory for file:", err)
			os.Exit(1)
		}

		err = os.WriteFile(outpath, f.Data, 0644)
		if err != nil {
			fmt.Println("write error:", err)
			os.Exit(1)
		}

		ok := cli.VerifySHA1(f)
		fmt.Printf("- %s (size=%d, ok=%v)\n", safe, f.Size, ok)
	}

	fmt.Printf("extracted %d files to %s\n", len(files), outdir)
}
