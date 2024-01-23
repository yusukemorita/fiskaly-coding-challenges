package domain

import (
	"errors"
	"fmt"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
)

type SignatureAlgorithm interface {
	Name() string
	GenerateEncodedPrivateKey() ([]byte, error)
}

type ECC struct{}

func (ecc ECC) Name() string {
	return "ECC"
}

func (ecc ECC) GenerateEncodedPrivateKey() ([]byte, error) {
	generator := crypto.ECCGenerator{}
	keyPair, err := generator.Generate()
	if err != nil {
		err = errors.New(fmt.Sprintf("failed to generate: %s", err.Error()))
		return []byte{}, err
	}

	marshaller := crypto.NewECCMarshaler()
	_, privateKey, err := marshaller.Encode(*keyPair)
	if err != nil {
		err = errors.New(fmt.Sprintf("failed to encode: %s", err.Error()))
		return []byte{}, err
	}

	return privateKey, nil
}

type RSA struct{}

func (rsa RSA) Name() string {
	return "RSA"
}

func (rsa RSA) GenerateEncodedPrivateKey() ([]byte, error) {
	generator := crypto.RSAGenerator{}
	keyPair, err := generator.Generate()
	if err != nil {
		return []byte{}, err
	}

	marshaller := crypto.NewRSAMarshaler()
	_, privateKey, err := marshaller.Marshal(*keyPair)
	if err != nil {
		return []byte{}, err
	}

	return privateKey, nil
}

var supportedAlgorithms = []SignatureAlgorithm{RSA{}, ECC{}}
