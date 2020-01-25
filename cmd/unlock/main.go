package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"os"

	"github.com/cobratbq/goutils/std/builtin"
	"github.com/cobratbq/goutils/std/errors"
)

func main() {
	flagInput := flag.String("i", "", "Iteration input as hexadecimal value.")
	flagNonce := flag.String("n", "", "Iteration nonce as hexadecimal value.")
	flagCiphertext := flag.String("c", "", "Iteration ciphertext to be decrypted as hexadecimal value.")
	flagSolution := flag.String("s", "", "Puzzle solution for iteration as hexadecimal value.")
	flagKey := flag.String("k", "", "Final decryption key used for decrypting ciphertext.")
	flag.Parse()

	builtin.Require(len(*flagCiphertext) > 0, "Please provide ciphertext.")
	builtin.Require(len(*flagInput) > 0, "Please provide input.")
	builtin.Require(len(*flagNonce) > 0, "Please provide nonce.")
	builtin.Require(len(*flagSolution) > 0, "Please provide solution.")

	input, err := hex.DecodeString(*flagInput)
	errors.RequireSuccess(err, "Iteration input is invalid")
	nonce, err := hex.DecodeString(*flagNonce)
	errors.RequireSuccess(err, "Iteration nonce is invalid")
	ciphertext, err := hex.DecodeString(*flagCiphertext)
	errors.RequireSuccess(err, "Iteration ciphertext is invalid")
	solution, err := hex.DecodeString(*flagSolution)
	errors.RequireSuccess(err, "Iteration solution is invalid")

	associated := append(solution, input...)
	var key []byte
	if *flagKey == "" {
		hashed := sha256.Sum256(associated)
		key = hashed[:]
	} else {
		key, err = hex.DecodeString(*flagKey)
		errors.RequireSuccess(err, "Failed to decode decryption key.")
	}
	blockcipher, err := aes.NewCipher(key)
	errors.RequireSuccess(err, "Failed to construct AES block cipher")
	aead, err := cipher.NewGCM(blockcipher)
	errors.RequireSuccess(err, "Failed to construct AEAD cipher")
	plaintext, err := aead.Open(nil, nonce, ciphertext, associated)
	errors.RequireSuccess(err, "Failed to decrypt ciphertext")
	os.Stderr.WriteString(fmt.Sprintf("Result: %s\n", string(plaintext)))
}
