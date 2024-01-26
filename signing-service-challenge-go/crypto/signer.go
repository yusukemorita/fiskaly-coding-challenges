package crypto

import "crypto/rand"

// Signer defines a contract for different types of signing implementations.
type Signer interface {
	Sign(dataToBeSigned []byte) ([]byte, error)
}

// TODO: implement RSA and ECDSA signing ...
type RSASigner struct {
	keyPair RSAKeyPair
}

func (signer RSASigner) Sign(dataToBeSigned []byte) ([]byte, error) {
	// TODO: use PSS or PKCS?
	return signer.keyPair.Private.Sign(rand.Reader, dataToBeSigned, nil)
}

type ECCSigner struct {
	keyPair ECCKeyPair
}

func (signer ECCSigner) Sign(dataToBeSigned []byte) ([]byte, error) {
	return signer.keyPair.Private.Sign(rand.Reader, dataToBeSigned, nil)
}
