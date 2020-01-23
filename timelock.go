package timelock

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"os"
)

const (
	inputSize = 20
)

// Timelock encrypts plaintext data using a simple (as-of-yet unproven, probably
// totally insecure) time-lock encryption mechanism.
func Timelock(plaintext []byte, n, complexity int) [][]byte {
	input := generateRandom(inputSize)
	os.Stdout.WriteString(fmt.Sprintf("Input: %0x\n", input))

	var associated []byte
	var key [32]byte
	for i := 0; i < n-1; i++ {
		// all-but-final time-lock iteration
		puzzle := generateRandom(complexity)
		// FIXME should we inject more data in Sha256 hash function? (we're already adding some additional data)
		associated = append(puzzle, input...)
		key = sha256.Sum256(associated)

		input = generateRandom(inputSize)
		interim, nonce := sealPayload(input, associated, key)
		os.Stdout.WriteString(fmt.Sprintf("Nonce %d: %0x interim %d: %0x\n", i, nonce, i, interim))
		os.Stderr.WriteString(fmt.Sprintf("KEEP SECRET: Puzzle %d: %0x input %d: %0x\n", i, puzzle, i, input))
	}

	// final time-lock iteration
	puzzle := generateRandom(complexity)
	associated = append(puzzle, input...)
	key = sha256.Sum256(associated)
	secretKey := generateKey()
	lastInterim, secretKeyNonce := sealPayload(secretKey[:], associated, key)
	os.Stdout.WriteString(fmt.Sprintf("Nonce %d: %0x interim %d: %0x\n", n-1, secretKeyNonce, n-1, lastInterim))
	os.Stderr.WriteString(fmt.Sprintf("KEEP SECRET: Puzzle %d: %0x\n", n-1, puzzle))

	// sealing away the actual plaintext
	// FIXME should we find some associated data such that it is linked to the last time-lock iteration?
	ciphertext, nonce := sealPayload(plaintext, nil, secretKey)

	os.Stdout.WriteString(fmt.Sprintf("Nonce: %0x ciphertext: %0x\n", nonce, ciphertext))

	return nil
}

func sealPayload(plaintext, associatedData []byte, key [32]byte) ([]byte, [12]byte) {
	aead := newAes256AEAD(key)
	nonce := generateNonce()
	sealed := aead.Seal(nil, nonce[:], plaintext, associatedData)
	return sealed, nonce
}

func newAes256AEAD(key [32]byte) cipher.AEAD {
	blockcipher, err := aes.NewCipher(key[:])
	requireSuccess(err, "failed to construct AES block cipher")
	aead, err := cipher.NewGCM(blockcipher)
	requireSuccess(err, "failed to construct AES-based AEAD cipher")
	return aead
}

func generateKey() [32]byte {
	var value [32]byte
	n, err := rand.Read(value[:])
	requireSuccess(err, "failure while generating random bytes for key")
	require(n == len(value), "Failed to read sufficient random data")
	return value
}

func generateNonce() [12]byte {
	var value [12]byte
	n, err := rand.Read(value[:])
	requireSuccess(err, "failure while generating random bytes for nonce")
	require(n == len(value), "Failed to read sufficient random data")
	return value
}

func generateRandom(size int) []byte {
	data := make([]byte, size)
	n, err := rand.Read(data)
	requireSuccess(err, "failed to read random data")
	require(n == size, "Failed to read sufficient random data")
	return data
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
