package assets

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/raf555/kbbi-api/internal/crypto"
)

type (
	reader struct {
		filename, dir string
		key, nonce    []byte
	}
)

func Read(filename, dir string, key, nonce []byte) *reader {
	return &reader{filename, dir, key, nonce}
}

func (r *reader) To(target any) error {
	ciphertext, err := os.ReadFile(path.Join(r.dir, r.filename))
	if err != nil {
		return fmt.Errorf("os.ReadFile: %w", err)
	}

	plaintext, err := crypto.Decrypt(r.key, r.nonce, ciphertext)
	if err != nil {
		return fmt.Errorf("crypto.Decrypt: %w", err)
	}

	reader := bytes.NewReader(plaintext)
	gz, err := gzip.NewReader(reader)
	if err != nil {
		return fmt.Errorf("gzip.NewReader: %w", err)
	}

	if err = json.NewDecoder(gz).Decode(target); err != nil {
		return fmt.Errorf("json.NewDecoder.Decode: %w", err)
	}

	_ = gz.Close()
	return nil
}
