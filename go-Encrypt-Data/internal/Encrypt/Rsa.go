package encrypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
)

type RSA struct {
	pub  *rsa.PublicKey
	priv *rsa.PrivateKey
}

func NewRSA() *RSA {
	privKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	return &RSA{
		pub:  &privKey.PublicKey,
		priv: privKey,
	}
}

func (r *RSA) Encrypt(msg string) string {
	msgBytes := []byte(msg)
	cipherText, _ := rsa.EncryptOAEP(sha256.New(), rand.Reader, r.pub, msgBytes, nil)
	// Encode ciphertext as base64 string for safe transport/storage
	return base64.StdEncoding.EncodeToString(cipherText)
}
func (r *RSA) Decrypt(msg string) string {
	cipherBytes, _ := base64.StdEncoding.DecodeString(msg)

	plainText, _ := rsa.DecryptOAEP(sha256.New(), rand.Reader, r.priv, cipherBytes, nil)
	return string(plainText)
}
