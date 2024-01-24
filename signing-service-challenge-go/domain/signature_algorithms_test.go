package domain

import (
	"testing"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
)

func TestEccAlgorithm_Name(t *testing.T) {
	ecc := eccAlgorithm{}
	name := ecc.Name()

	if name != "ECC" {
		t.Errorf("expected ECC, got: %s", name)
	}
}

func TestEccAlgorithm_GenerateEncodedPrivateKey(t *testing.T) {
	ecc := eccAlgorithm{}
	encodedPrivateKey, err := ecc.GenerateEncodedPrivateKey()

	if err != nil {
		t.Errorf("expected no error, got %s", err)
	}

	_, err = crypto.NewECCMarshaler().Decode(encodedPrivateKey)
	if err != nil {
		t.Errorf("decode of generated private key failed: %s", err)
	}
}

func TestRsaAlgorithm_Name(t *testing.T) {
	rsa := rsaAlgorithm{}
	name := rsa.Name()

	if name != "RSA" {
		t.Errorf("expected RSA, got: %s", name)
	}
}

func TestRsaAlgorithm_GenerateEncodedPrivateKey(t *testing.T) {
	rsa := rsaAlgorithm{}
	encodedPrivateKey, err := rsa.GenerateEncodedPrivateKey()

	if err != nil {
		t.Errorf("expected no error, got %s", err)
	}

	_, err = crypto.NewRSAMarshaler().Unmarshal(encodedPrivateKey)
	if err != nil {
		t.Errorf("decode of generated private key failed: %s", err)
	}
}

func TestFindSupportedAlgorithm(t *testing.T) {
	t.Run("returns found: false when algorithm does not exist", func(t *testing.T) {
		_, found := FindSupportedAlgorithm("INVALID")
		if found {
			t.Error("expected found: false")
		}
	})

	t.Run("returns algorithm when algorithm exists", func(t *testing.T) {
		algorithm, found := FindSupportedAlgorithm("RSA")
		if !found {
			t.Error("expected found: true")
		}

		if algorithm.Name() != "RSA" {
			t.Errorf("expected %s, got %s", "RSA", algorithm.Name())
		}
	})
}
