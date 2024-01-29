package domain

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type KeyPair interface {
	Sign(dataToBeSigned []byte) (signature []byte, err error)
	EncodedPublicKey() (string, error)
	AlgorithmName() string
}

type KeyPairGenerator interface {
	AlgorithmName() string
	Generate() (KeyPair, error)
}

type SignatureDevice struct {
	ID      uuid.UUID
	KeyPair KeyPair
	// (optional) user provided string to be displayed in the UI
	Label string
	// track the last signature created with this device
	Base64EncodedLastSignature string
	// track how many signatures have been created with this device
	SignatureCounter uint
}

func (device SignatureDevice) Sign(dataToBeSigned string) ([]byte, error) {
	return device.KeyPair.Sign([]byte(dataToBeSigned))
}

func BuildSignatureDevice(id uuid.UUID, generator KeyPairGenerator, label ...string) (SignatureDevice, error) {
	keyPair, err := generator.Generate()
	if err != nil {
		err = errors.New(fmt.Sprintf("key pair generation failed: %s", err.Error()))
		return SignatureDevice{}, err
	}

	device := SignatureDevice{
		ID:      id,
		KeyPair: keyPair,
	}

	if len(label) > 0 {
		device.Label = label[0]
	}

	return device, nil
}

// WARNING:
// All operations must be executed inside WriteTx() or ReadTx(),
// as Go maps are not safe for concurrent use.
// For this reason, do not expose the `SignatureDeviceRepository` directly
// to the `api` package.
// Instead, expose the `SignatureDeviceRepositoryProvider`, which ensures
// that every operation will be executed inside a transaction.
type SignatureDeviceRepository interface {
	Create(device SignatureDevice) error
	// Increment the signatureCounter, and update the lastSignature
	MarkSignatureCreated(deviceID uuid.UUID, newSignature string) error
	Find(id uuid.UUID) (SignatureDevice, bool, error)
	List() ([]SignatureDevice, error)
}

type SignatureDeviceRepositoryProvider interface {
	WriteTx(func(SignatureDeviceRepository) error) error
	ReadTx(func(SignatureDeviceRepository) error) error
}
