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
	if len(b) == 0 {
		return "", -1
	}

	for typ, sig := range magics {
		if len(sig) == 0 || len(b) < len(sig) {
			continue
		}
		idx := bytes.Index(b, sig)
		if
