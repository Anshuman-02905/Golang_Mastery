package main

import (
	encrypt "Encrypt/internal/Encrypt"
	"fmt"
)

func main() {
	message := "Hello Golang"
	shift := 3
	caeser := encrypt.NewCaeser(shift)
	fmt.Printf("Before CaeserEncryton %v \n", message)
	caeserEncrypted := caeser.Encrypt(message)
	fmt.Printf("After Caeser Eccrytption %v \n", caeserEncrypted)

	// Initialize RSA

	rsaCipher := encrypt.NewRSA()
	// Step 2: RSA Encryption on Caesar output
	rsaEncryptted := rsaCipher.Encrypt(caeserEncrypted)
	fmt.Printf("After RSA Encryptipn (raw Bytes) %v \n: ", []byte(rsaEncryptted))

	//Step3" RSA DECRYTPTION
	rsaDecrytpted := rsaCipher.Decrypt(rsaEncryptted)
	fmt.Printf("After RSA DECRYPTION %v \n", rsaDecrytpted)

	//Step4 Caeser Decrytption

	original := caeser.Decrypt(rsaDecrytpted)
	fmt.Printf("original %v \n", original)

}
