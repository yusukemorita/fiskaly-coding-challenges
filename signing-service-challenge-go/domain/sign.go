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
	repositoryProvider SignatureDeviceRepositoryProvider,
	dataToBeSigned string,
) (
	deviceFound bool,
	encodedSignature string,
	signedData string,
	err error,
) {
	txErr := repositoryProvider.WriteTx(func(repository SignatureDeviceRepository) error {
		device, ok, err := repository.Find(deviceID)
		if err != nil {
			return err
		}
		if !ok {
			deviceFound = false
			return nil
		}
		deviceFound = true

		signedData = SecureDataToBeSigned(device, dataToBeSigned)
		signature, err := device.Sign(signedData)
		if err != nil {
			return errors.New(fmt.Sprintf("failed to sign transaction: %s", err))
		}
		encodedSignature = base64.StdEncoding.EncodeToString(signature)

		err = repository.MarkSignatureCreated(device.ID, encodedSignature)
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
		encodedLastSignature := device.LastSignature
		components = append(components, encodedLastSignature)
	}

	return strings.Join(components, "_")
}
