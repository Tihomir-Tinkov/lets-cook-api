package adapters

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

/*
	TODO: check correctness
	1. validate version
	2. validate decoded parameters against limits
	3. make parameter parsing stricter
	4. add NeedsRehash
	5. add tests for malformed hashes and DoS cases
*/

var (
	ErrInvalidPassword = errors.New("invalid password")
	ErrInvalidHash     = errors.New("invalid password hash")
)

type PasswordHasher struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	keyLength   uint32
	saltLength  uint32
}

func NewPasswordHasher() *PasswordHasher {
	return &PasswordHasher{
		memory:      64 * 1024, // 64 MB
		iterations:  3,
		parallelism: 4,
		keyLength:   32,
		saltLength:  16,
	}
}

func equalBytes(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}

	return subtle.ConstantTimeCompare(a, b) == 1
}

func parseHash(encoded string) (
	memory uint32,
	iterations uint32,
	parallelism uint8,
	salt []byte,
	hash []byte,
	err error,
) {

	parts := strings.Split(encoded, "$")

	if len(parts) != 6 {
		err = ErrInvalidHash
		return
	}

	if parts[1] != "argon2id" {
		err = ErrInvalidHash
		return
	}

	_, err = fmt.Sscanf(
		parts[3],
		"m=%d,t=%d,p=%d",
		&memory,
		&iterations,
		&parallelism,
	)

	if err != nil {
		return
	}

	salt, err = base64.RawStdEncoding.DecodeString(parts[4])

	if err != nil {
		return
	}

	hash, err = base64.RawStdEncoding.DecodeString(parts[5])

	return
}

func (a *PasswordHasher) Hash(password string) (string, error) {
	salt := make([]byte, a.saltLength)

	_, err := rand.Read(salt)

	if err != nil {
		return "", err
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		a.iterations,
		a.memory,
		a.parallelism,
		a.keyLength,
	)

	return fmt.Sprintf(
		"$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		a.memory,
		a.iterations,
		a.parallelism,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash),
	), nil
}

func (a *PasswordHasher) Compare(
	password string,
	encodedHash string,
) error {

	memory, iterations, parallelism, salt, hash, err :=
		parseHash(encodedHash)

	if err != nil {
		return ErrInvalidHash
	}

	comparisonHash := argon2.IDKey(
		[]byte(password),
		salt,
		iterations,
		memory,
		parallelism,
		uint32(len(hash)),
	)

	if !equalBytes(hash, comparisonHash) {
		return ErrInvalidPassword
	}

	return nil
}
