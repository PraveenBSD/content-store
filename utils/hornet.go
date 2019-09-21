package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"log"
)

func decryptPrivateKey(privateKeyFilePath string, password []byte) ([]byte, error) {

	privateKeyBytes, _ := ioutil.ReadFile(privateKeyFilePath)

	block, rest := pem.Decode(privateKeyBytes)
	if len(rest) > 0 {
		err := errors.New("Extra data included in key")
		log.Println(err)
		return nil, err
	}

	if x509.IsEncryptedPEMBlock(block) {
		der, err := x509.DecryptPEMBlock(block, password)
		if err != nil {
			log.Println("Decrypt failed: %v", err)
			return nil, err
		}
		return pem.EncodeToMemory(&pem.Block{Type: block.Type, Bytes: der}), nil
	}
	return privateKeyBytes, nil
}

func importPrivateKey(privateKey []byte) (*rsa.PrivateKey, error) {

	block, _ := pem.Decode(privateKey)
	if block == nil {
		log.Println("failed to decode PEM block containing public key")
		return nil, errors.New("failed to decode PEM block containing public key")
	}
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

func importPublicKey(publicKeyFilePath string) (*rsa.PublicKey, error) {

	publicKeyBytes, err := ioutil.ReadFile(publicKeyFilePath)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(publicKeyBytes)
	if block == nil {
		log.Println("failed to decode PEM block containing public key")
		return nil, errors.New("failed to decode PEM block containing public key")
	}
	return x509.ParsePKCS1PublicKey(block.Bytes)
}

// EncryptFile encrypts given file using specified public key
func EncryptFile(publicKeyFilePath string, message []byte) ([]byte, error) {

	publicKey, err := importPublicKey(publicKeyFilePath)
	if err != nil {
		log.Println("Error importing public key:", err)
		return nil, err
	}

	return rsa.EncryptPKCS1v15(rand.Reader, publicKey, message)
}

// DecryptFile decrypts a file using the specified private key
func DecryptFile(privateKeyFilePath string, encryptedFileInBytes []byte) ([]byte, error) {
	privateKeyPem, err := decryptPrivateKey(privateKeyFilePath, []byte("helloworld"))
	if err != nil {
		log.Println("error in decrypting private key")
		return nil, err
	}
	privateKey, err := importPrivateKey(privateKeyPem)
	if err != nil {
		log.Println("error in importing private key")
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, privateKey, encryptedFileInBytes)
}

// func main() {
// 	content, err := ioutil.ReadFile("/Users/praveen/go/src/github.com/PraveenBSD/content-store/contents/hello.txt")
// 	if err != nil {
// 		fmt.Print("\n\nerror reading file...:\n", err)
// 	}
// 	a, err := EncryptFile("/Users/praveen/go/src/github.com/PraveenBSD/content-store/keys/id_rsa.pub.pem", content)
// 	if err != nil {
// 		fmt.Print("\n\nerror encrypting...:\n", err)
// 	}
// 	fmt.Print("\nencrypted text", a)
// 	// err = ioutil.WriteFile("/Users/praveen/go/src/github.com/PraveenBSD/content-store/contents/Encrypted/9bc0-h892",
// 	// 	a, 777)
// 	// if err != nil {
// 	// 	fmt.Print("\n\nerror writing file...:\n", err)
// 	// }
// 	b, err := DecryptFile("/Users/praveen/go/src/github.com/PraveenBSD/content-store/keys/id_rsa", a)
// 	if err != nil {
// 		fmt.Print("\n\nerror decrypting...:\n", err)
// 	}
// 	fmt.Print("\ndecrypted text: ", string(b), " done!")

// }
