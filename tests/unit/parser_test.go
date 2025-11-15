package cli

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

type File struct {
	GUID string
	Name string
	Ext  string
	Hash string
	Size int
	Data []byte
}

func VerifySHA1(f File) bool {
	if f.Hash == "" {
		return true
	}
	sum := sha1.Sum(f.Data)
	return fmt.Sprintf("%x", sum) == f.Hash
}

func ParseEnv(path string) ([]File, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return ParseEnvBytes(b)
}

func ParseEnvBytes(b []byte) ([]File, error) {
	if len(b) == 0 {
		return nil, errors.New("empty input")
	}
	if bytes.HasPrefix(b, []byte("FILENAME/")) {
		return parseSimpleHeaderFormat(b)
	}
	return parseContainerFormat(b)
}

func parseSimpleHeaderFormat(b []byte) ([]File, error) {
	parts := bytes.SplitN(b, []byte("\n"), 4)
	if len(parts) < 4 {
		return nil, errors.New("invalid simple header")
	}
	filenameLine := string(parts[0])
	extLine := string(parts[1])
	sha1Line := string(parts[2])
	data := parts[3]
	var name, ext, hash string
	if strings.HasPrefix(filenameLine, "FILENAME/") {
		name = strings.TrimPrefix(filenameLine, "FILENAME/")
	}
	if strings.HasPrefix(extLine, "EXT/") {
		ext = strings.TrimPrefix(extLine, "EXT/")
	}
	if strings.HasPrefix(sha1Line, "SHA1/") {
		hash = strings.TrimPrefix(sha1Line, "SHA1/")
	}
	f := File{
		Name: name,
		Ext:  ext,
		Hash: hash,
		Size: len(data),
		Data: data,
	}
	return []File{f}, nil
}

func parseContainerFormat(b []byte) ([]File, error) {
	sep := []byte{0x2A, 0x2A, 0x25, 0x25} // **%%
	parts := bytes.Split(b, sep)
	filesByGUID := map[string]*File{}
	order := []string{}
	lastGUID := ""
	for _, part := range parts {
		if len(part) == 0 {
			continue
		}
		if looksMostlyBinary(part) {
			attachBinaryPart(part, filesByGUID, &order, &lastGUID)
			continue
		}
		text := sanitizeText(part)
		guid, name, ext, hash, doctype := extractFields(text)
		if guid != "" {
			lastGUID = guid
			f, ok := filesByGUID[guid]
			if !ok {
				f = &File{GUID: guid}
				filesByGUID[guid] = f
				order = append(order, guid)
			}
			if name != "" {
				f.Name = name
			}
			if ext != "" {
				f.Ext = ext
			}
			if hash != "" {
				f.Hash = hash
			}
			if doctype != "" && f.Ext == "" {
				f.Ext = extFromDoctype(doctype)
			}
			continue
		}
		if name != "" || ext != "" || hash != "" || doctype != "" {
			var target *File
			if lastGUID != "" {
				target = filesByGUID[lastGUID]
			}
			if target == nil {
				f := &File{Name: name, Ext: ext, Hash: hash}
				anon := fmt.Sprintf("__anon_%p", f)
				filesByGUID[anon] = f
				order = append(order, anon)
				lastGUID = anon
			} else {
				if name != "" {
					target.Name = name
				}
				if ext != "" {
					target.Ext = ext
				}
				if hash != "" {
					target.Hash = hash
				}
				if doctype != "" && target.Ext == "" {
					target.Ext = extFromDoctype(doctype)
				}
			}
			continue
		}
	
		if idx := indexBinarySignature(part); idx >= 0 {
			bin := part[idx:]
			attachBinaryPart(bin, filesByGUID, &order, &lastGUID)
		}
	}
	out := []File{}
	for _, g := range order {
		f := filesByGUID[g]
		if f == nil {
			continue
		}
		f.Size = len(f.Data)
		out = append(out, *f)
	}
	return out, nil
}

func sanitizeText(b []byte) string {
	s := string(b)
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\x00", "")
	return s
}

func extractFields(s string) (guid, name, ext, hash, doctype string) {
	lines := strings.Split(s, "\n")
	for _, ln := range lines {
		ln = strings.TrimSpace(ln)
		if ln == "" {
			continue
		}
		if strings.HasPrefix(ln, "GUID/") {
			guid = strings.TrimPrefix(ln, "GUID/")
			guid = strings.TrimSpace(guid)
		} else if strings.HasPrefix(ln, "FILENAME/") {
			name = strings.TrimPrefix(ln, "FILENAME/")
			name = strings.TrimSpace(name)
		} else if strings.HasPrefix(ln, "EXT/") {
			ext = strings.TrimPrefix(ln, "EXT/")
			ext = strings.TrimSpace(ext)
		} else if strings.HasPrefix(ln, "SHA1/") {
			hash = strings.TrimPrefix(ln, "SHA1/")
			hash = strings.TrimSpace(hash)
		} else if strings.HasPrefix(ln, "DOCTYPE/") || strings.HasPrefix(ln, "TYPE/") {
			if strings.HasPrefix(ln, "DOCTYPE/") {
				doctype = strings.TrimPrefix(ln, "DOCTYPE/")
			} else {
				doctype = strings.TrimPrefix(ln, "TYPE/")
			}
			doctype = strings.TrimSpace(doctype)
		} else if strings.HasPrefix(ln, "_SIG/") {
			// ignore signature marker
		}
	}
	return
}

func extFromDoctype(dt string) string {
	dt = strings.ToUpper(dt)
	if strings.Contains(dt, "IMAGE") || strings.Contains(dt, "IMAG") {
		return "jpg"
	}
	if strings.Contains(dt, "JPG") || strings.Contains(dt, "JPEG") {
		return "jpg"
	}
	if strings.Contains(dt, "PNG") {
		return "png"
	}
	if strings.Contains(dt, "PDF") {
		return "pdf"
	}
	return ""
}

func looksMostlyBinary(b []byte) bool {
	if len(b) == 0 {
		return false
	}
	printable := 0
	for _, c := range b {
		if c >= 0x20 && c <= 0x7E {
			printable++
		} else if c == '\n' || c == '\r' || c == '\t' {
			printable++
		}
	}
	return float64(printable)/float64(len(b)) < 0.6
}

func indexBinarySignature(b []byte) int {
	sigs := [][]byte{
		[]byte{0xFF, 0xD8, 0xFF},       
		[]byte{0x89, 'P', 'N', 'G'},   
		[]byte{'R', 'I', 'F', 'F'},      
		[]byte{'P', '%', 'P', 'S'},     
		[]byte{'%', 'P', 'D', 'F'},     
		[]byte{'B', 'M'},               
		[]byte{'W', 'E', 'B', 'P'},      
	}
	for i := 0; i < len(b); i++ {
		for _, s := range sigs {
			if i+len(s) <= len(b) && bytes.Equal(b[i:i+len(s)], s) {
				return i
			}
		}
	}
	return -1
}

func attachBinaryPart(part []byte, filesByGUID map[string]*File, order *[]string, lastGUID *string) {
	if len(part) == 0 {
		return
	}
	binIdx := indexBinarySignature(part)
	var data []byte
	if binIdx >= 0 {
		data = part[binIdx:]
	} else {
		data = part
	}
	if *lastGUID != "" {
		f := filesByGUID[*lastGUID]
		if f == nil {
			f = &File{GUID: *lastGUID}
			filesByGUID[*lastGUID] = f
			*order = append(*order, *lastGUID)
		}
		f.Data = append(f.Data, data...)
		f.Size = len(f.Data)
		return
	}
	anon := &File{Data: data, Size: len(data)}
	key := fmt.Sprintf("__anon_%p", anon)
	filesByGUID[key] = anon
	*order = append(*order, key)
}

func hexPreview(b []byte) string {
	l := 32
	if len(b) < l {
		l = len(b)
	}
	return hex.EncodeToString(b[:l])
}
