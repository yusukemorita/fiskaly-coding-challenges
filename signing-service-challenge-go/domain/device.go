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
	LastSignature string
	// track how many signatures have been created with this device
	SignatureCounter uint
}

func (device SignatureDevice) SignTransaction(dataToBeSigned string) ([]byte, error) {
	return device.Algorithm.SignTransaction(
		device.EncodedPrivateKey,
		[]byte(dataToBeSigned),
	)
}

func (device SignatureDevice) ExtendDataToBeSigned(data string) string {
	if device.SignatureCounter == 0 {
		// when the device has not yet been used, the `lastSignature` is blank,
		// so use the device ID instead
		encodedID := base64.StdEncoding.EncodeToString([]byte(device.ID.String()))

		return strings.Join(
			[]string{
				strconv.Itoa(int(device.SignatureCounter)),
				data,
				encodedID,
			},
			"_",
		)
	}

	encodedLastSignature := base64.StdEncoding.EncodeToString([]byte(device.LastSignature))

	// The resulting string (secured_data_to_be_signed) should follow this format:
	// <signature_counter>_<data_to_be_signed>_<last_signature_base64_encoded>
	return strings.Join(
		[]string{
			strconv.Itoa(int(device.SignatureCounter)),
			data,
			encodedLastSignature,
		},
		"_",
	)
}

// func SignTransaction(
// 	device SignatureDevice,
// 	deviceRepository SignatureDeviceRepository,
// 	dataToBeSigned string,
// ) (
// 	signature string,
// 	signedData string,
// 	err error,
// ) {

// 	sig, err := device.SignTransaction(securedDataToBeSigned)
// 	if err != nil {
// 		// TODO: better error handling?
// 		return "", "", errors.New("failed to sign transaction")
// 	}

// 	// TODO: base64 encode signature
// 	// TODO: update counter and last signature
// 	// TODO: update

// 	return signature, "", nil
// }

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
