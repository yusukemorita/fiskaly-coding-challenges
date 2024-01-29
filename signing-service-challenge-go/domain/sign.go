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
	encodedSignature string,
	signedData string,
	err error,
) {
	txErr := deviceRepository.Tx(func() error {
		device, ok, err := deviceRepository.Find(deviceID)
		if err != nil {
			return err
		}
		if !ok {
			deviceFound = false
			return nil
		}
		deviceFound = true

		// 1. Read `lastSignature` and `signatureCounter` from `device`
		//    (these values cannot change until update is complete)
		//    If data was persisted in MySQL, for example, a locking read would be used
		signedData = SecureDataToBeSigned(device, dataToBeSigned)

		// 2. Use the data read in 1. to create the signature
		signature, err := device.Sign(signedData)
		if err != nil {
			return errors.New(fmt.Sprintf("failed to sign transaction: %s", err))
		}
		encodedSignature = base64.StdEncoding.EncodeToString(signature)

		// 3. Update the device, and release the lock
		device.Base64EncodedLastSignature = encodedSignature
		device.SignatureCounter++
		err = deviceRepository.Update(device)
		if err != nil {
			return errors.New(fmt.Sprintf("failed to update signature device: %s", err))
		}

		return nil
	})

	if txErr != nil {
		return false, "", "", txErr
	}

	return
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
