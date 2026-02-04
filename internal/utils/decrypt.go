package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha512"
	"encoding/hex"
	"errors"

	"golang.org/x/crypto/pbkdf2"
)

// Decryptor decrypts AES-256-CBC data encrypted with PBKDF2-derived key (matches Node.js Encryptor).
type Decryptor struct {
	key []byte
	iv  []byte
}

// NewDecryptor creates a Decryptor from secretKey, iv (hex), salt, iterations, and keylen (same as Node).
func NewDecryptor(secretKey, ivHex, salt string, iterations, keylen int) (*Decryptor, error) {
	iv, err := hex.DecodeString(ivHex)
	if err != nil {
		return nil, err
	}
	if len(iv) != aes.BlockSize {
		return nil, errors.New("iv must be 16 bytes (32 hex chars)")
	}
	key := pbkdf2.Key([]byte(secretKey), []byte(salt), iterations, keylen, sha512.New)
	return &Decryptor{key: key, iv: iv}, nil
}

// Decrypt decrypts hex-encoded ciphertext and returns the UTF-8 plaintext (matches Node decrypt).
func (d *Decryptor) Decrypt(encryptedHex string) (string, error) {
	if encryptedHex == "" {
		return "", nil
	}
	ciphertext, err := hex.DecodeString(encryptedHex)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(d.key)
	if err != nil {
		return "", err
	}
	if len(ciphertext)%aes.BlockSize != 0 {
		return "", errors.New("ciphertext is not a multiple of block size")
	}
	mode := cipher.NewCBCDecrypter(block, d.iv)
	mode.CryptBlocks(ciphertext, ciphertext)
	plaintext, err := unpadPKCS7(ciphertext)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

func unpadPKCS7(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, errors.New("empty data")
	}
	pad := int(data[len(data)-1])
	if pad > len(data) || pad > aes.BlockSize || pad == 0 {
		return nil, errors.New("invalid PKCS7 padding")
	}
	for i := len(data) - pad; i < len(data); i++ {
		if data[i] != byte(pad) {
			return nil, errors.New("invalid PKCS7 padding")
		}
	}
	return data[:len(data)-pad], nil
}
