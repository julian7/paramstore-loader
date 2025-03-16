package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

const RedactedMagic = "\000REDACTED\000"

type Secret string

type secret struct {
	Filename string `json:"file"`
}

var ErrRedacted = errors.New("file is redacted")

func (s *Secret) UnmarshallJSON(b []byte) error {
	if b[0] == '"' {
		var directsecret string
		if err := json.Unmarshal(b, &directsecret); err != nil {
			return err
		}

		*s = Secret(directsecret)
		return nil
	}

	var filesecret secret
	if err := json.Unmarshal(b, &filesecret); err != nil {
		return err
	}

	data, err := readFile(filesecret.Filename)
	if err != nil {
		return err
	}
	*s = Secret(string(data))
	return nil
}

func readFile(fname string) ([]byte, error) {
	data, err := os.ReadFile(fname)
	if err != nil {
		return nil, err
	}

	if bytes.HasPrefix(data, []byte(RedactedMagic)) {
		return nil, fmt.Errorf("%q: %w", fname, ErrRedacted)
	}

	return data, nil
}
