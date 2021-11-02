package passwordEncrypt

import (
	"crypto/sha256"
	"golang.org/x/crypto/chacha20poly1305"
	"hcc/flute/lib/config"
)

// EncryptPassword : Encrypt password with chacha20poly1305 algorithm
func EncryptPassword(password string) []byte {
	key := sha256.Sum256([]byte(config.Ipmi.PasswordEncryptSecretKey))
	aead, _ := chacha20poly1305.New(key[:])

	nonce := make([]byte, chacha20poly1305.NonceSize)
	encryptedPassword := aead.Seal(nil, nonce, []byte(password), nil)

	return encryptedPassword
}

// DecryptPassword : Decrypt password with chacha20poly1305 algorithm
func DecryptPassword(encryptedPassword []byte) string {
	key := sha256.Sum256([]byte(config.Ipmi.PasswordEncryptSecretKey))
	aead, _ := chacha20poly1305.New(key[:])

	nonce := make([]byte, chacha20poly1305.NonceSize)
	decryptedPassword, _ := aead.Open(nil, nonce, encryptedPassword, nil)

	return string(decryptedPassword)
}
