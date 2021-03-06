package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

//CreateHash create the hask for the encryption
func CreateHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

//CheckFileExtension checks if the file extension exists
func CheckFileExtension(filename string) bool {
	if filename[len(filename)-4:] == ".txt" {
		return true
	}
	return false
}

//Encrypt encrypts the data
func Encrypt(data []byte, passphrase string) []byte {
	block, err := aes.NewCipher([]byte(CreateHash(passphrase)))
	if err != nil {
		panic(err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext
}

//Decrypt decrypts the data
func Decrypt(data []byte, passphrase string) []byte {
	key := []byte(CreateHash(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}
	return plaintext
}

//EncryptFile encrypts the file
func EncryptFile(filename string, data []byte, passphrase string) {
	f, err := os.Create(filename)
	if err != nil {
		panic(err.Error())
	}
	defer f.Close()
	f.Write(Encrypt(data, passphrase))
}

//DecryptFile decrypts the file
func DecryptFile(filename string, passphrase string) []byte {
	inFile, err := os.Open(filename)
	if err != nil {
		panic(err.Error())
	}
	defer inFile.Close()
	cipherText, err := ioutil.ReadAll(inFile)
	if err != nil {
		panic(err.Error())
	}
	return Decrypt(cipherText, passphrase)
}

//ScanUserInput scans the user input
func ScanUserInput(prompt string) string {
	var input string
	fmt.Print(prompt)
	fmt.Scanln(&input)
	return input
}

//ChooseEncryptOrDecrypt allows user to choose encryption or decryption
func ChooseEncryptOrDecrypt() string {
	var input string
	fmt.Print("Encrypt or Decrypt? (e/d): ")
	fmt.Scanln(&input)
	if input != "e" && input != "d" {
		fmt.Println("Invalid input. Try again.")
		return ChooseEncryptOrDecrypt()
	}
	return input
}
