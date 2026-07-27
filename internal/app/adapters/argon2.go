package adapters

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"golang.org/x/crypto/argon2"
)

var (
	ErrInvalidPassword = errors.New("invalid password")
	ErrInvalidHash     = errors.New("invalid password hash")
)

const (
	// These limits protect against attacker-controlled hashes causing excessive
	// memory/CPU usage during verification.
	maxMemoryKiB  = 1024 * 1024 // 1 GiB in KiB (Argon2 uses KiB)
	maxIterations = 20
	maxParallel   = 32

	minSaltLength = 16
	maxSaltLength = 64

	minKeyLength = 16
	maxKeyLength = 128

	maxEncodedHashLength = 512
	maxPasswordLength    = 1024
)

type Argon2Hasher struct {
	config Argon2Config
}

type hashParams struct {
	Version     int
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	Salt        []byte
	Hash        []byte
}

type Argon2Config struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

func NewArgon2Hasher(config Argon2Config) (*Argon2Hasher, error) {
	if config.Memory == 0 {
		return nil, errors.New("argon2 memory must be greater than zero")
	}

	if config.Iterations == 0 {
		return nil, errors.New("argon2 iterations must be greater than zero")
	}

	if config.Parallelism == 0 {
		return nil, errors.New("argon2 parallelism must be greater than zero")
	}

	if config.SaltLength < minSaltLength {
		return nil, fmt.Errorf("argon2 salt length must be at least %d", minSaltLength)
	}

	if config.KeyLength < minKeyLength {
		return nil, fmt.Errorf("argon2 key length must be at least %d", minKeyLength)
	}

	if config.SaltLength > maxSaltLength {
		return nil, fmt.Errorf("argon2 salt length must be at most %d", maxSaltLength)
	}

	if config.KeyLength > maxKeyLength {
		return nil, fmt.Errorf("argon2 key length must be at most %d", maxKeyLength)
	}

	if config.Memory > maxMemoryKiB {
		return nil, fmt.Errorf("argon2 memory exceeds maximum")
	}

	if config.Iterations > maxIterations {
		return nil, fmt.Errorf("argon2 iterations exceeds maximum")
	}

	if config.Parallelism > maxParallel {
		return nil, fmt.Errorf("argon2 parallelism exceeds maximum")
	}

	return &Argon2Hasher{
		config: config,
	}, nil
}

func (a *Argon2Hasher) Hash(password string) (string, error) {
	if len(password) > maxPasswordLength {
		return "", ErrInvalidPassword
	}
	salt := make([]byte, a.config.SaltLength)

	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		a.config.Iterations,
		a.config.Memory,
		a.config.Parallelism,
		a.config.KeyLength,
	)

	return fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		a.config.Memory,
		a.config.Iterations,
		a.config.Parallelism,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash),
	), nil
}

func (a *Argon2Hasher) Compare(password, encodedHash string) error {
	if len(password) > maxPasswordLength {
		return ErrInvalidPassword
	}
	params, err := parseHash(encodedHash)

	// log?
	if err != nil {
		return ErrInvalidPassword
	}

	comparisonHash := argon2.IDKey(
		[]byte(password),
		params.Salt,
		params.Iterations,
		params.Memory,
		params.Parallelism,
		uint32(len(params.Hash)),
	)

	if subtle.ConstantTimeCompare(params.Hash, comparisonHash) != 1 {
		return ErrInvalidPassword
	}

	return nil
}

// NeedsRehash reports whether the stored hash should be regenerated using the current hasher configuration.
func (a *Argon2Hasher) NeedsRehash(encodedHash string) bool {
	params, err := parseHash(encodedHash)
	if err != nil {
		return true
	}

	return params.Version != argon2.Version ||
		params.Memory != a.config.Memory ||
		params.Iterations != a.config.Iterations ||
		params.Parallelism != a.config.Parallelism ||
		uint32(len(params.Salt)) != a.config.SaltLength ||
		uint32(len(params.Hash)) != a.config.KeyLength
}

func parseHash(encoded string) (*hashParams, error) {
	if len(encoded) > maxEncodedHashLength {
		return nil, ErrInvalidHash
	}
	parts := strings.Split(encoded, "$")

	// "", "argon2id", "v=19", "m=...,t=...,p=...", salt, hash
	if len(parts) != 6 {
		return nil, ErrInvalidHash
	}

	if parts[1] != "argon2id" {
		return nil, ErrInvalidHash
	}

	params := &hashParams{}

	if !strings.HasPrefix(parts[2], "v=") {
		return nil, ErrInvalidHash
	}

	version, err := strconv.Atoi(strings.TrimPrefix(parts[2], "v="))
	if err != nil {
		return nil, ErrInvalidHash
	}

	params.Version = version

	if params.Version != argon2.Version {
		return nil, ErrInvalidHash
	}

	if err := parseArgon2Params(parts[3], params); err != nil {
		return nil, err
	}

	if params.Memory > maxMemoryKiB ||
		params.Iterations > maxIterations ||
		params.Parallelism > maxParallel {
		return nil, ErrInvalidHash
	}

	var saltErr error

	params.Salt, saltErr = base64.RawStdEncoding.DecodeString(parts[4])
	if saltErr != nil ||
		len(params.Salt) < minSaltLength ||
		len(params.Salt) > maxSaltLength {
		return nil, ErrInvalidHash
	}

	var hashErr error

	params.Hash, hashErr = base64.RawStdEncoding.DecodeString(parts[5])
	if hashErr != nil ||
		len(params.Hash) < minKeyLength ||
		len(params.Hash) > maxKeyLength {
		return nil, ErrInvalidHash
	}

	return params, nil
}

func parseArgon2Params(input string, params *hashParams) error {
	seen := map[string]bool{}

	for _, part := range strings.Split(input, ",") {
		key, value, ok := strings.Cut(part, "=")
		if !ok || seen[key] {
			return ErrInvalidHash
		}

		seen[key] = true

		switch key {
		case "m":
			v, err := strconv.ParseUint(value, 10, 32)
			if err != nil {
				return ErrInvalidHash
			}
			params.Memory = uint32(v)

		case "t":
			v, err := strconv.ParseUint(value, 10, 32)
			if err != nil {
				return ErrInvalidHash
			}
			params.Iterations = uint32(v)

		case "p":
			v, err := strconv.ParseUint(value, 10, 8)
			if err != nil {
				return ErrInvalidHash
			}
			params.Parallelism = uint8(v)

		default:
			return ErrInvalidHash
		}
	}

	if !seen["m"] || !seen["t"] || !seen["p"] {
		return ErrInvalidHash
	}

	if params.Memory == 0 ||
		params.Iterations == 0 ||
		params.Parallelism == 0 {
		return ErrInvalidHash
	}

	return nil
}
