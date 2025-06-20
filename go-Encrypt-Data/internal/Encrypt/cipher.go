package encrypt


type Cipher interface{
	Encrypt(string) string
	Decrypt(string) string
}
