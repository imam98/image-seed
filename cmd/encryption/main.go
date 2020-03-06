package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"log"
)

func main() {
	key := make([]byte, 32)
	rand.Read(key)

	plaintext, err := ioutil.ReadFile("message.txt")
	if err != nil {
		log.Fatalln(err.Error())
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatalln(err.Error())
	}

	nonce := make([]byte, 12)
	if _, err := rand.Read(nonce); err != nil {
		log.Fatalln(err.Error())
	}
	fmt.Println("Key:", key)
	fmt.Println("Nonce:", nonce)
	key = append(key, nonce...)

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Println(err.Error())
	}

	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)
	fmt.Println(ciphertext)

	if err := ioutil.WriteFile("enc.key", key, 0644); err != nil {
		log.Fatalln(err.Error())
	}

	if err := ioutil.WriteFile("cipher.encrypt", ciphertext, 0644); err != nil {
		log.Fatalln(err.Error())
	}
}
