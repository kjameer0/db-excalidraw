package internal

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/exp/rand"
)

func GenerateShortKey() string {

	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const keyLength = 10

	rand.Seed(uint64(time.Now().UnixNano()))
	shortKey := make([]byte, keyLength)
	for i := range shortKey {
		shortKey[i] = charset[rand.Intn(len(charset))]
	}
	return string(shortKey)
}

// func HashToUUID(hashString string) {

// }
func UuidToHash(uuid uuid.UUID, key []byte) string {
	plaintext := []byte(uuid.String())

	// Create AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// Initialize GCM (Galois Counter Mode)
	nonce := make([]byte, 12) // 12 bytes nonce
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}

	// Encrypt
	ciphertext := aesGCM.Seal(nil, nonce, plaintext, nil)
	fmt.Printf("Encrypted: %s\n", hex.EncodeToString(ciphertext))
	// encryptyed := hex.EncodeToString(ciphertext)
	// Decrypt
	decrypted, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Decrypted: %s\n", decrypted)

	return string(decrypted)
}
