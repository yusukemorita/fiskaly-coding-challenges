package domain

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type SignatureAlgorithm string

const (
	RSA SignatureAlgorithm = "RSA"
	ECC SignatureAlgorithm = "ECC"
)

var supportedAlgorithms = []SignatureAlgorithm{RSA, ECC}

type SignatureDevice struct {
	id                uuid.UUID
	algorithm         SignatureAlgorithm
	encodedPrivateKey []byte
	// (optional) user provided string to be displayed in the UI
	label string
	// track the last signature created with this device
	lastSignature string
	// track how many signatures have been created with this device
	signatureCounter uint
}

func BuildSignatureDevice(id string, algorithm string, label ...string) (SignatureDevice, error) {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		err = errors.New(fmt.Sprintf("invalid uuid: %s", err.Error()))
		return SignatureDevice{}, err
	}

	var parsedAlgorithm SignatureAlgorithm
	for _, alg := range supportedAlgorithms {
		if alg == SignatureAlgorithm(algorithm) {
			parsedAlgorithm = alg
			break
		}
	}
	if parsedAlgorithm == "" {
		return SignatureDevice{}, errors.New("invalid algorithm")
	}

	// TODO: generate key pair and set value of encodedPrivateKey
	// TODO: set label when provided
	return SignatureDevice{
		id:        parsedId,
		algorithm: parsedAlgorithm,
	}, nil
}
