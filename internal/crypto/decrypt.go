package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

func Decrypt(key, nonce, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("error creating cipher block from key: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("error wrapping block with GCM: %w", err)
	}

	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("error decrypting ciphertext: %w", err)
	}

	return plaintext, nil
}
