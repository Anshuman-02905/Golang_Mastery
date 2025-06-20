package encrypt

import "crypto/aes"

type AES struct{
	key string
	aes.NewCipher
}
