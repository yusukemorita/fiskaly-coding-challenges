package domain

import (
	"testing"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
)

func TestECC_Name(t *testing.T) {
	ecc := ECC{}
	name := ecc.Name()

	if name != "ECC" {
		t.Errorf("expected ECC, got: %s", name)
	}
}

func TestECC_GenerateEncodedPrivateKey(t *testing.T) {
	ecc := ECC{}
	encodedPrivateKey, err := ecc.GenerateEncodedPrivateKey()

	if err != nil {
		t.Errorf("expected no error, got %s", err)
	}

	_, err = crypto.NewECCMarshaler().Decode(encodedPrivateKey)
	if err != nil {
		t.Errorf("decode of generated private key failed: %s", err)
	}
}

func TestRSA_Name(t *testing.T) {
	rsa := RSA{}
	name := rsa.Name()

	if name != "RSA" {
		t.Errorf("expected RSA, got: %s", name)
	}
}

func TestRSA_GenerateEncodedPrivateKey(t *testing.T) {
	rsa := RSA{}
	encodedPrivateKey, err := rsa.GenerateEncodedPrivateKey()

	if err != nil {
		t.Errorf("expected no error, got %s", err)
	}

	_, err = crypto.NewRSAMarshaler().Unmarshal(encodedPrivateKey)
	if err != nil {
		t.Errorf("decode of generated private key failed: %s", err)
	}
}
