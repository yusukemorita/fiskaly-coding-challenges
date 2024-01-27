package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

// RSAKeyPair is a DTO that holds RSA private and public keys.
type RSAKeyPair struct {
	Public  *rsa.PublicKey
	Private *rsa.PrivateKey
}

func (keyPair RSAKeyPair) EncodedPrivateKey() ([]byte, error) {
	marshaler := RSAMarshaler{}
	_, private, err := marshaler.Marshal(keyPair)
	return private, err
}

func (keyPair RSAKeyPair) SignTransaction(dataToBeSigned []byte) (signature []byte, err error) {
	digest, err := computeDigestWithHashFunction(dataToBeSigned)
	if err != nil {
		return nil, err
	}

	return rsa.SignPSS(
		rand.Reader,
		keyPair.Private,
		hashFunction,
		digest,
		nil,
	)
}

// RSAMarshaler can encode and decode an RSA key pair.
type RSAMarshaler struct{}

// NewRSAMarshaler creates a new RSAMarshaler.
func NewRSAMarshaler() RSAMarshaler {
	return RSAMarshaler{}
}

// Marshal takes an RSAKeyPair and encodes it to be written on disk.
// It returns the public and the private key as a byte slice.
func (m RSAMarshaler) Marshal(keyPair RSAKeyPair) (encodedPublicKey, encodedPrivateKey []byte, err error) {
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(keyPair.Private)
	publicKeyBytes := x509.MarshalPKCS1PublicKey(keyPair.Public)

	encodedPrivate := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA_PRIVATE_KEY",
		Bytes: privateKeyBytes,
	})

	encodedPublic := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA_PUBLIC_KEY",
		Bytes: publicKeyBytes,
	})

	return encodedPublic, encodedPrivate, nil
}

// Unmarshal takes an encoded RSA private key and transforms it into a rsa.PrivateKey.
func (m RSAMarshaler) Unmarshal(privateKeyBytes []byte) (*RSAKeyPair, error) {
	block, _ := pem.Decode(privateKeyBytes)
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return &RSAKeyPair{
		Private: privateKey,
		Public:  &privateKey.PublicKey,
	}, nil
}

// Implements domain.SignatureAlgorithm for RSA.
// Note that any actual logic is implemented in `RSASigner`, `RSAMarshaller` and `RSAGenerator`,
// and this struct merely acts as a facade to make this logic easier to access in the
// `domain` package.
type RSAAlgorithm struct{}

func (rsa RSAAlgorithm) Name() string {
	return "RSA"
}
