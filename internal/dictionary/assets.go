package dictionary

import (
	"bytes"
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"encoding/json"
	"fmt"
	"os"
	"path"
)

type (
	reader struct {
		filename, dir string
		key, nonce    []byte
	}
)

func ReadAsset(filename, dir string, key, nonce []byte) *reader {
	return &reader{filename, dir, key, nonce}
}

func (r *reader) To(target any) error {
	ciphertext, err := os.ReadFile(path.Join(r.dir, r.filename))
	if err != nil {
		return fmt.Errorf("os.ReadFile: %w", err)
	}

	plaintext, err := r.decrypt(r.key, r.nonce, ciphertext)
	if err != nil {
		return fmt.Errorf("r.decrypt: %w", err)
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

func (r *reader) decrypt(key, nonce, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("aes.NewCipher: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("cipher.NewGCM: %w", err)
	}

	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("aesGCM.Open: %w", err)
	}

	return plaintext, nil
}
