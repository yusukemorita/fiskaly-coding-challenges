package domain

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

// Defines the algorithm related functions that `domain` package requires.
// These operations will be implemented by algorithm specific structs in the
// `crypto` package.
// e.g. `RSAAlgorithm`, `ECCAlgorithm`
type SignatureAlgorithm interface {
	Name() string
	GenerateEncodedPrivateKey() ([]byte, error)
	SignTransaction(encodedPrivateKey []byte, dataToBeSigned []byte) (signature []byte, err error)
}

type SignatureDevice struct {
	ID                uuid.UUID
	Algorithm         SignatureAlgorithm
	EncodedPrivateKey []byte
	// (optional) user provided string to be displayed in the UI
	Label string
	// track the last signature created with this device
	Base64EncodedLastSignature string
	// track how many signatures have been created with this device
	SignatureCounter uint
}

func (device SignatureDevice) SignTransaction(dataToBeSigned string) ([]byte, error) {
	return device.Algorithm.SignTransaction(
		device.EncodedPrivateKey,
		[]byte(dataToBeSigned),
	)
}

func (device SignatureDevice) SecureDataToBeSigned(data string) string {
	components := []string{
		strconv.Itoa(int(device.SignatureCounter)),
		data,
	}

	if device.SignatureCounter == 0 {
		// when the device has not yet been used, the `lastSignature` is blank,
		// so use the device ID instead
		encodedID := base64.StdEncoding.EncodeToString([]byte(device.ID.String()))
		components = append(components, encodedID)
	} else {
		encodedLastSignature := base64.StdEncoding.EncodeToString([]byte(device.Base64EncodedLastSignature))
		components = append(components, encodedLastSignature)
	}

	return strings.Join(components, "_")
}

func BuildSignatureDevice(id uuid.UUID, algorithm SignatureAlgorithm, label ...string) (SignatureDevice, error) {
	encodedPrivateKey, err := algorithm.GenerateEncodedPrivateKey()
	if err != nil {
		err = errors.New(fmt.Sprintf("private key generation failed: %s", err.Error()))
		return SignatureDevice{}, err
	}

	device := SignatureDevice{
		ID:                id,
		Algorithm:         algorithm,
		EncodedPrivateKey: encodedPrivateKey,
	}

	if len(label) > 0 {
		device.Label = label[0]
	}

	return device, nil
}

type SignatureDeviceRepository interface {
	Create(device SignatureDevice) error
	Update(device SignatureDevice) error
	Find(id uuid.UUID) (SignatureDevice, bool, error)
}
