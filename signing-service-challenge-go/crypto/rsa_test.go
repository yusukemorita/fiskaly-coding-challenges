package crypto

import (
	"crypto/rsa"
	"testing"
)

func TestRSASigner_Sign(t *testing.T) {
	generator := RSAGenerator{}
	keyPair, err := generator.generate()
	if err != nil {
		t.Fatal(err)
	}

	dataToBeSigned := "some-data"
	signature, err := keyPair.Sign([]byte(dataToBeSigned))
	if err != nil {
		t.Fatal(err)
	}

	digest, err := computeDigestWithHashFunction([]byte(dataToBeSigned))
	if err != nil {
		t.Fatal(err)
	}

	err = rsa.VerifyPSS(keyPair.Public, hashFunction, digest, signature, nil)
	if err != nil {
		t.Errorf("signature verification failed: %s", err)
	}
}
