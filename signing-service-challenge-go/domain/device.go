package domain

import (
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
	LastSignature string
	// track how many signatures have been created with this device
	SignatureCounter uint
}

func SignTransaction(
	device SignatureDevice,
	deviceRepository SignatureDeviceRepository,
	dataToBeSigned string,
	findSupportedAlgorithm func(name string) (SignatureAlgorithm, bool),
) (
	signature string,
	signedData string,
	err error,
) {
	algorithm, found := findSupportedAlgorithm(device.AlgorithmName)
	// TODO: better error handling?
	if !found {
		return "", "", errors.New("algorithm is not supported")
	}

	// TODO: when signature counter is 0
	// The resulting string (secured_data_to_be_signed) should follow this format:
	// <signature_counter>_<data_to_be_signed>_<last_signature_base64_encoded>
	securedDataToBeSigned := strings.Join(
		[]string{
			strconv.Itoa(int(device.SignatureCounter)),
			dataToBeSigned,
			device.LastSignature,
		},
		"_",
	)

	sig, err := algorithm.SignTransaction(device.EncodedPrivateKey, []byte(securedDataToBeSigned))
	if err != nil {
		// TODO: better error handling?
		return "", "", errors.New("failed to sign transaction")
	}

	// TODO: base64 encode signature
	// TODO: update counter and last signature
	// TODO: update

	return signature, "", nil
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
	Find(id uuid.UUID) (SignatureDevice, bool, error)
}
