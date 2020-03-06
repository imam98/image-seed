package main

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"io/ioutil"
	"log"
)

func main() {
	ciphertext, err := ioutil.ReadFile("cipher.encrypt")
	if err != nil {
		log.Fatalln(err.Error())
	}

	key, err := ioutil.ReadFile("enc.key")
	if err != nil {
		log.Fatalln(err.Error())
	}
	nonce := key[32:]

	block, err := aes.NewCipher(key[:32])
	if err != nil {
		log.Fatalln(err.Error())
	}
	fmt.Println("key:", key[:32])
	fmt.Println("Nonce:", nonce)

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Fatalln(err.Error())
	}

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		log.Fatalln(err.Error())
	}

	if err := ioutil.WriteFile("message.txt", plaintext, 0644); err != nil {
		log.Fatalln(err.Error())
	}
	fmt.Println(string(plaintext))
}
