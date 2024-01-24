package domain

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type SignatureDevice struct {
	id                uuid.UUID
	algorithmName     string
	encodedPrivateKey []byte
	// (optional) user provided string to be displayed in the UI
	label string
	// track the last signature created with this device
	lastSignature string
	// track how many signatures have been created with this device
	signatureCounter uint
}

func BuildSignatureDevice(id uuid.UUID, algorithm signatureAlgorithm, label ...string) (SignatureDevice, error) {
	encodedPrivateKey, err := algorithm.GenerateEncodedPrivateKey()
	if err != nil {
		err = errors.New(fmt.Sprintf("private key generation failed: %s", err.Error()))
		return SignatureDevice{}, err
	}

	device := SignatureDevice{
		id:                id,
		algorithmName:     algorithm.Name(),
		encodedPrivateKey: encodedPrivateKey,
	}

	if len(label) > 0 {
		device.label = label[0]
	}

	return device, nil
}
