package crypto

import (
	stdCrypto "crypto"
)

// At the current RSA key size (512 bits), SHA512 causes a
// `message too long for RSA key size` error from `rsa.SignPSS`
// ref: https://github.com/golang/go/blob/d7df7f4fa01f1b445d835fc908c54448a63c68fb/src/crypto/rsa/pss.go#L304-L308
// Therefore use the next biggest hash function, as a bigger hash
// makes collisions less likely.
const hashFunction = stdCrypto.SHA384

// Signer defines a contract for different types of signing implementations.
type Signer interface {
	Sign(dataToBeSigned []byte) ([]byte, error)
}

func computeDigestWithHashFunction(b []byte) ([]byte, error) {
	hash := hashFunction.New()
	_, err := hash.Write(b)
	if err != nil {
		return nil, err
	}
	return hash.Sum(nil), nil
}
