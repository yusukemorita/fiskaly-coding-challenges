package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
)

// RSAGenerator generates a RSA key pair.
type RSAGenerator struct{}

func (g RSAGenerator) AlgorithmName() string {
	return "RSA"
}

// Generate generates a new RSAKeyPair.
func (g RSAGenerator) Generate() (domain.KeyPair, error) {
	// Security has been ignored for the sake of simplicity.
	key, err := rsa.GenerateKey(rand.Reader, 512)
	if err != nil {
		return nil, err
	}

	return &RSAKeyPair{
		Public:  &key.PublicKey,
		Private: key,
	}, nil
}

// ECCGenerator generates an ECC key pair.
type ECCGenerator struct{}

func (g ECCGenerator) AlgorithmName() string {
	return "ECC"
}

// Generate generates a new ECCKeyPair.
func (g ECCGenerator) Generate() (domain.KeyPair, error) {
	// Security has been ignored for the sake of simplicity.
	key, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		return nil, err
	}

	return &ECCKeyPair{
		Public:  &key.PublicKey,
		Private: key,
	}, nil
}
