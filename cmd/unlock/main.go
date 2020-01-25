package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
)

func main() {
	flagInput := flag.String("i", "", "Iteration input as hexadecimal value.")
	flagNonce := flag.String("n", "", "Iteration nonce as hexadecimal value.")
	flagCiphertext := flag.String("c", "", "Iteration ciphertext to be decrypted as hexadecimal value.")
	flagSolution := flag.String("s", "", "Puzzle solution for iteration as hexadecimal value.")
	flagKey := flag.String("k", "", "Final decryption key used for decrypting ciphertext.")
	flag.Parse()

	require(len(*flagCiphertext) > 0, "Please provide ciphertext.")
	require(len(*flagInput) > 0, "Please provide input.")
	require(len(*flagNonce) > 0, "Please provide nonce.")
	require(len(*flagSolution) > 0, "Please provide solution.")

	input, err := hex.DecodeString(*flagInput)
	requireSuccess(err, "Iteration input is invalid")
	nonce, err := hex.DecodeString(*flagNonce)
	requireSuccess(err, "Iteration nonce is invalid")
	ciphertext, err := hex.DecodeString(*flagCiphertext)
	requireSuccess(err, "Iteration ciphertext is invalid")
	solution, err := hex.DecodeString(*flagSolution)
	requireSuccess(err, "Iteration solution is invalid")

	associated := append(solution, input...)
	var key []byte
	if *flagKey == "" {
		hashed := sha256.Sum256(associated)
		key = hashed[:]
	} else {
		key, err = hex.DecodeString(*flagKey)
		requireSuccess(err, "Failed to decode decryption key.")
	}
	blockcipher, err := aes.NewCipher(key)
	requireSuccess(err, "Failed to construct AES block cipher")
	aead, err := cipher.NewGCM(blockcipher)
	requireSuccess(err, "Failed to construct AEAD cipher")
	plaintext, err := aead.Open(nil, nonce, ciphertext, associated)
	requireSuccess(err, "Failed to decrypt ciphertext")
	os.Stderr.WriteString(fmt.Sprintf("Result: %s\n", string(plaintext)))
}

func requireSuccess(err error, message string) {
	if err != nil {
		panic("Fatal: " + message + ": " + err.Error())
	}
}

func require(condition bool, message string) {
	if !condition {
		panic("Fatal: " + message)
	}
}
