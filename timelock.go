package timelock

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"fmt"
	"os"

	"github.com/cobratbq/goutils/std/builtin"
	"github.com/cobratbq/goutils/std/crypto/rand"
)

const (
	inputSize = 20
)

// Timelock encrypts plaintext data using a simple (as-of-yet unproven, probably
// totally insecure) time-lock encryption mechanism.
func Timelock(plaintext []byte, n, complexity int) [][]byte {
	input := generateRandom(inputSize)
	os.Stdout.WriteString(fmt.Sprintf("%0x\n", input))

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
		os.Stdout.WriteString(fmt.Sprintf("%0x %0x\n", nonce, interim))
		os.Stderr.WriteString(fmt.Sprintf("KEEP SECRET: Puzzle %d: %0x input %d: %0x\n", i, puzzle, i, input))
	}

	// final time-lock iteration
	puzzle := generateRandom(complexity)
	associated = append(puzzle, input...)
	key = sha256.Sum256(associated)
	secretKey := generateKey()
	lastInterim, secretKeyNonce := sealPayload(secretKey[:], associated, key)
	os.Stdout.WriteString(fmt.Sprintf("%0x %0x\n", secretKeyNonce, lastInterim))
	os.Stderr.WriteString(fmt.Sprintf("KEEP SECRET: Puzzle %d: %0x\n", n-1, puzzle))

	// sealing away the actual plaintext
	ciphertext, nonce := sealPayload(plaintext, associated, secretKey)

	os.Stdout.WriteString(fmt.Sprintf("%0x %0x\n", nonce, ciphertext))

	return nil
}

func sealPayload(plaintext, associatedData []byte, key [32]byte) ([]byte, [12]byte) {
	aead := newAes256GCM(key)
	nonce := generateNonce()
	sealed := aead.Seal(nil, nonce[:], plaintext, associatedData)
	return sealed, nonce
}

func newAes256GCM(key [32]byte) cipher.AEAD {
	blockcipher, err := aes.NewCipher(key[:])
	builtin.RequireSuccess(err, "failed to construct AES block cipher")
	aead, err := cipher.NewGCM(blockcipher)
	builtin.RequireSuccess(err, "failed to construct AES-based AEAD cipher")
	return aead
}

func generateKey() [32]byte {
	var value [32]byte
	rand.MustReadBytes(value[:])
	return value
}

func generateNonce() [12]byte {
	var value [12]byte
	rand.MustReadBytes(value[:])
	return value
}

func generateRandom(size int) []byte {
	return rand.RandomizeBytes(make([]byte, size))
}
