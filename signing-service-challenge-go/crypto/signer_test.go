package crypto

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/sha256"
	"testing"
)

func TestRSASigner_Sign(t *testing.T) {
	generator := RSAGenerator{}
	keyPair, err := generator.Generate()
	if err != nil {
		t.Fatal(err)
	}

	dataToBeSigned := "some-data"
	signer := RSASigner{keyPair: *keyPair}
	signature, err := signer.Sign([]byte(dataToBeSigned))
	if err != nil {
		t.Fatal(err)
	}

	hash := sha256.New()
	_, err = hash.Write([]byte(dataToBeSigned))
	if err != nil {
		t.Fatal(err)
	}
	digest := hash.Sum(nil)

	err = rsa.VerifyPKCS1v15(keyPair.Public, crypto.SHA256, digest, signature)
	if err != nil {
		t.Errorf("signature verification failed: %s", err)
	}
}

func TestECCSigner_Sign(t *testing.T) {
	generator := ECCGenerator{}
	keyPair, err := generator.Generate()
	if err != nil {
		t.Fatal(err)
	}

	dataToBeSigned := "some-data"
	signer := ECCSigner{keyPair: *keyPair}
	signature, err := signer.Sign([]byte(dataToBeSigned))
	if err != nil {
		t.Fatal(err)
	}

	digest, err := computeDigestWithHashFunction([]byte(dataToBeSigned))
	if err != nil {
		t.Fatal(err)
	}

	result := ecdsa.VerifyASN1(keyPair.Public, digest, signature)
	if !result {
		t.Errorf("signature verification failed: %s", err)
	}
}