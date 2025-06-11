package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"io"
)

type EncryptionService interface {
	EncryptWithAES(plaintext, key string) (string, error)
	DecryptWithAES(ciphertext, key string) (string, error)
	EncryptWithRSA(plaintext, publicKeyPEM string) (string, error)
	DecryptWithRSA(ciphertext, privateKeyPEM string) (string, error)
	GenerateAESKey() (string, error)
}

type encryptionService struct{}

func NewEncryptionService() EncryptionService {
	return &encryptionService{}
}

// EncryptWithAES encrypts plaintext using AES-GCM
func (s *encryptionService) EncryptWithAES(plaintext, key string) (string, error) {
	if len(key) != 32 {
		return "", errors.New("key must be 32 bytes for AES-256")
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptWithAES decrypts ciphertext using AES-GCM
func (s *encryptionService) DecryptWithAES(ciphertext, key string) (string, error) {
	if len(key) != 32 {
		return "", errors.New("key must be 32 bytes for AES-256")
	}

	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertextBytes := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// EncryptWithRSA encrypts plaintext using RSA public key
func (s *encryptionService) EncryptWithRSA(plaintext, publicKeyPEM string) (string, error) {
	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil {
		return "", errors.New("failed to parse PEM block containing the public key")
	}

	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", err
	}

	rsaPublicKey, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		return "", errors.New("not an RSA public key")
	}

	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, rsaPublicKey, []byte(plaintext), nil)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptWithRSA decrypts ciphertext using RSA private key
func (s *encryptionService) DecryptWithRSA(ciphertext, privateKeyPEM string) (string, error) {
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return "", errors.New("failed to parse PEM block containing the private key")
	}

	privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}

	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	plaintext, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privKey, data, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// GenerateAESKey generates a random 32-byte AES key
func (s *encryptionService) GenerateAESKey() (string, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(key), nil
}

