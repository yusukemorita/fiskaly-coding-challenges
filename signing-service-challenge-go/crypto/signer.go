package crypto

import (
	stdCrypto "crypto"
	"crypto/rand"
	"crypto/rsa"
)

// Signer defines a contract for different types of signing implementations.
type Signer interface {
	Sign(dataToBeSigned []byte) ([]byte, error)
}

// TODO: implement RSA and ECDSA signing ...

type RSASigner struct {
	keyPair RSAKeyPair
}

// TODO: test this and ECCSigner.Sign
func (signer RSASigner) Sign(dataToBeSigned []byte) ([]byte, error) {
	sha256 := stdCrypto.SHA256 // TODO: is this a good hash to use?

	// compute digest of `dataToBeSigned` with SHA256
	hash := sha256.New()
	_, err := hash.Write(dataToBeSigned)
	if err != nil {
		return nil, err
	}
	digest := hash.Sum(nil)

	// TODO: use PSS or PKCS?
	return rsa.SignPKCS1v15(
		rand.Reader,
		signer.keyPair.Private,
		sha256,
		digest,
	)

	// return signer.keyPair.Private.Sign(
	// 	rand.Reader,
	// 	digest,
	// 	sha256,
	// )
}

type ECCSigner struct {
	keyPair ECCKeyPair
}

func (signer ECCSigner) Sign(dataToBeSigned []byte) ([]byte, error) {
	return signer.keyPair.Private.Sign(rand.Reader, dataToBeSigned, nil)
}
