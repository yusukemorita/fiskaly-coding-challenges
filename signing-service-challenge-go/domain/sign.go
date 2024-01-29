package domain

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

func SignTransaction(
	deviceID uuid.UUID,
	deviceRepository SignatureDeviceRepository,
	dataToBeSigned string,
) (
	deviceFound bool,
	base64EncodedSignature string,
	signedData string,
	err error,
) {
	device, ok, err := deviceRepository.Find(deviceID)
	if err != nil {
		return false, "", "", err
	}
	if !ok {
		return false, "", "", nil
	}

	// 1. Read `lastSignature` and `signatureCounter` from `device`
	//    (these values cannot change until update is complete)
	//    If data was persisted in MySQL, for example, a locking read would be used
	securedDataToBeSigned := SecureDataToBeSigned(device, dataToBeSigned)

	// 2. Use the data read in 1. to create the signature
	signature, err := device.Sign(securedDataToBeSigned)
	if err != nil {
		return false, "", "", errors.New(fmt.Sprintf("failed to sign transaction: %s", err))
	}
	encodedSignature := base64.StdEncoding.EncodeToString(signature)

	// 3. Update the device, and release the lock
	device.Base64EncodedLastSignature = encodedSignature
	device.SignatureCounter++
	err = deviceRepository.Update(device)
	if err != nil {
		return false, "", "", errors.New(fmt.Sprintf("failed to update signature device: %s", err))
	}

	return true, encodedSignature, securedDataToBeSigned, nil
}

func SecureDataToBeSigned(device SignatureDevice, data string) string {
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
		encodedLastSignature := device.Base64EncodedLastSignature
		components = append(components, encodedLastSignature)
	}

	return strings.Join(components, "_")
}
