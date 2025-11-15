package cli

import (
	"bytes"
	"errors"
	"strings"
)

type EnvFile struct {
	Name string
	Ext  string
	Hash string
	Data []byte
}

func ParseEnvBytes(b []byte) ([]EnvFile, error) {
	lines, data := splitHeaderAndData(b)
	if len(lines) == 0 {
		return nil, errors.New("no header")
	}

	f := EnvFile{}

	for _, line := range lines {
		if strings.HasPrefix(line, "FILENAME/") {
			f.Name = strings.TrimPrefix(line, "FILENAME/")
		} else if strings.HasPrefix(line, "EXT/") {
			f.Ext = strings.TrimPrefix(line, "EXT/")
		} else if strings.HasPrefix(line, "SHA1/") {
			f.Hash = strings.TrimPrefix(line, "SHA1/")
		}
	}

	f.Data = data

	return []EnvFile{f}, nil
}

func splitHeaderAndData(b []byte) ([]string, []byte) {
	parts := bytes.SplitN(b, []byte("\n"), 4)

	if len(parts) < 4 {
		return nil, nil
	}

	lines := []string{
		string(parts[0]),
		string(parts[1]),
		string(parts[2]),
	}

	return lines, parts[3]
}
