package cli

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type FileData struct {
	Name string
	Ext  string
	Hash string
	Size int
	Data []byte
	Type string
}

var magics = map[string][]byte{
	"jpg":  {0xFF, 0xD8, 0xFF},
	"png":  {0x89, 'P', 'N', 'G'},
	"gif":  {'G', 'I', 'F', '8'},
	"pdf":  {'%', 'P', 'D', 'F', '-'},
	"xml":  {'<', '?', 'x', 'm', 'l'},
	"zip":  {'P', 'K', 0x03, 0x04},
	"webp": {'R', 'I', 'F', 'F'},
}

func detectMagic(b []byte) (string, int) {
	for typ, sig := range magics {
		idx := bytes.Index(b, sig)
		if idx >= 0 {
			return typ, idx
		}
	}
	idxXML := bytes.Index(b, []byte("<?xml"))
	if idxXML >= 0 {
		return "xml", idxXML
	}
	return "", -1
}

func parseHeaderText(header []byte) map[string]string {
	lines := bytes.Split(header, []byte{'\n'})
	res := make(map[string]string)
	for _, ln := range lines {
		ln = bytes.TrimSpace(ln)
		if len(ln) == 0 {
			continue
		}
		parts := strings.SplitN(string(ln), "/", 2)
		if len(parts) == 2 {
			res[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}
	return res
}

func ParseEnvBytes(data []byte) ([]FileData, error) {
	blocks := bytes.Split(data, []byte("**%%"))
	var files []FileData
	for i, blk := range blocks {
		if len(bytes.TrimSpace(blk)) == 0 {
			continue
		}
		magicType, magicPos := detectMagic(blk)
		if magicPos == -1 {
			continue
		}
		header := blk[:magicPos]
		content := blk[magicPos:]
		meta := parseHeaderText(header)
		name := meta["FILENAME"]
		if name == "" {
			name = fmt.Sprintf("block-%03d", i)
		}
		ext := meta["EXT"]
		if ext == "" {
			ext = magicType
		}
		hash := strings.TrimSpace(meta["SHA1"])
		size := len(content)
		files = append(files, FileData{
			Name: name,
			Ext:  ext,
			Hash: hash,
			Size: size,
			Data: content,
			Type: magicType,
		})
	}
	return files, nil
}

func ParseEnv(path string) ([]FileData, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return ParseEnvBytes(b)
}

func VerifySHA1(f FileData) bool {
	if f.Hash == "" {
		return true
	}
	sum := sha1.Sum(f.Data)
	return strings.EqualFold(hex.EncodeToString(sum[:]), strings.TrimSpace(f.Hash))
}

func WriteOutputs(files []FileData, outdir string) error {
	for _, f := range files {
		safe := filepath.Clean(f.Name)
		if filepath.Ext(safe) == "" && f.Ext != "" {
			safe = safe + "." + strings.TrimPrefix(f.Ext, ".")
		}
		path := filepath.Join(outdir, safe)
		if err := ioutil.WriteFile(path, f.Data, 0644); err != nil {
			return err
		}
	}
	return nil
}
