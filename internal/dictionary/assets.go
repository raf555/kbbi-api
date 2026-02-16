package dictionary

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"time"
)

type (
	reader struct {
		filename, dir string
		url           string

		key, nonce []byte
	}
)

func ReadAsset(filename, dir string, key, nonce []byte) *reader {
	return &reader{filename, dir, "", key, nonce}
}

func ReadAssetFromURL(url string, key, nonce []byte) *reader {
	return &reader{"", "", url, key, nonce}
}

func (r *reader) To(target any) error {
	var (
		ciphertext []byte
		err        error
	)

	if r.url != "" {
		ciphertext, err = r.getCiphertextFromURL(context.TODO())
	} else {
		ciphertext, err = r.getCiphertextFromFile()
	}
	if err != nil {
		return fmt.Errorf("get ciphertext: %w", err)
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

func (r *reader) getCiphertextFromURL(ctx context.Context) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second) // hardcode for now
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, r.url, nil)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http do: %w", err)
	}
	defer func() {
		_ = res.Body.Close()
	}()

	ciphertext, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("io readall: %w", err)
	}

	return ciphertext, nil
}

func (r *reader) getCiphertextFromFile() ([]byte, error) {
	ciphertext, err := os.ReadFile(path.Join(r.dir, r.filename))
	if err != nil {
		return nil, fmt.Errorf("os.ReadFile: %w", err)
	}
	return ciphertext, nil
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
