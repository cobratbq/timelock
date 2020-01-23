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
	flagDecrypt := flag.Bool("d", false, "Decrypt final value.")
	flag.Parse()

	input, err := hex.DecodeString(*flagInput)
	requireSuccess(err, "Iteration input is invalid")
	nonce, err := hex.DecodeString(*flagNonce)
	requireSuccess(err, "Iteration nonce is invalid")
	ciphertext, err := hex.DecodeString(*flagCiphertext)
	requireSuccess(err, "Iteration ciphertext is invalid")
	solution, err := hex.DecodeString(*flagSolution)
	requireSuccess(err, "Iteration solution is invalid")

	if !*flagDecrypt {
		associated := append(solution, input...)
		key := sha256.Sum256(associated)
		blockcipher, err := aes.NewCipher(key[:])
		requireSuccess(err, "Failed to construct AES block cipher")
		aead, err := cipher.NewGCM(blockcipher)
		requireSuccess(err, "Failed to construct AEAD cipher")

		plaintext, err := aead.Open(nil, nonce, ciphertext, associated)
		requireSuccess(err, "Failed to decrypt ciphertext")
		os.Stderr.WriteString(fmt.Sprintf("Result: %0x\n", plaintext))
	} else {
		blockcipher, err := aes.NewCipher(input)
		requireSuccess(err, "Failed to construct AES block cipher")
		aead, err := cipher.NewGCM(blockcipher)
		requireSuccess(err, "Failed to construct AEAD cipher")

		plaintext, err := aead.Open(nil, nonce, ciphertext, nil)
		requireSuccess(err, "Failed to decrypt ciphertext")
		os.Stderr.WriteString(fmt.Sprintf("Plaintext: %s\n", plaintext))
	}
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
