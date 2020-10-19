package sec

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"

	"golang.org/x/crypto/pbkdf2"
)

const (
	iterCount = 10
)

// Crypter is an interface of a encrypter/decrypter
type Crypter interface {
	Encrypt([]byte) ([]byte, error)
	Decrypt([]byte) ([]byte, error)
}

// AESCrypter is a crypting server implementing Crypter interface
type AESCrypter struct {
	gcm cipher.AEAD
}

// NewAES creates a new AESCrypter
func NewAES(key []byte) (*AESCrypter, error) {
	c := new(AESCrypter)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	c.gcm = gcm

	return c, nil
}

// Encrypt encrypts a plaintext data
func (c *AESCrypter) Encrypt(plaintext []byte) ([]byte, error) {
	nonce := make([]byte, c.gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	data := []byte(plaintext)
	ciphertext := c.gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

// Decrypt decrypts a previously encrypted data
func (c *AESCrypter) Decrypt(ciphertext []byte) ([]byte, error) {
	nonceSize := c.gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("encrypted data too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := c.gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// CreatePassKey creates an AES-compatible key from a passphrase
func CreatePassKey(passphrase []byte) ([]byte, error) {
	hasher := md5.New()
	hasher.Write(passphrase)
	salt := hasher.Sum(nil)

	phash := pbkdf2.Key(passphrase, salt, iterCount, 4096, sha1.New)
	hasher.Reset()
	hasher.Write(phash)
	pwdKey := hex.EncodeToString(hasher.Sum(nil))
	return []byte(pwdKey), nil
}
