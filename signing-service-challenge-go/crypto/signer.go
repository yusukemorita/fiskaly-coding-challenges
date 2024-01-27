package crypto

import (
	stdCrypto "crypto"
	"crypto/rand"
	"crypto/rsa"
)

// TODO: is this a good hash to use?
const hashFunction = stdCrypto.SHA256

// Signer defines a contract for different types of signing implementations.
type Signer interface {
	Sign(dataToBeSigned []byte) ([]byte, error)
}

type RSASigner struct {
	keyPair RSAKeyPair
}

func (signer RSASigner) Sign(dataToBeSigned []byte) ([]byte, error) {
	digest, err := computeDigestWithHashFunction(dataToBeSigned)
	if err != nil {
		return nil, err
	}

	return rsa.SignPSS(
		rand.Reader,
		signer.keyPair.Private,
		hashFunction,
		digest,
		nil,
	)
}

type ECCSigner struct {
	keyPair ECCKeyPair
}

func (signer ECCSigner) Sign(dataToBeSigned []byte) ([]byte, error) {
	digest, err := computeDigestWithHashFunction(dataToBeSigned)
	if err != nil {
		return nil, err
	}

	return signer.keyPair.Private.Sign(rand.Reader, digest, nil)
}

func computeDigestWithHashFunction(b []byte) ([]byte, error) {
	hash := hashFunction.New()
	_, err := hash.Write(b)
	if err != nil {
		return nil, err
	}
	return hash.Sum(nil), nil
}
