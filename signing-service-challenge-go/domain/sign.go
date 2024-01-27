package domain

import (
	"encoding/base64"
	"errors"
	"fmt"
)

func SignTransaction(
	device SignatureDevice,
	deviceRepository SignatureDeviceRepository,
	dataToBeSigned string,
) (
	base64EncodedSignature string,
	signedData string,
	err error,
) {
	securedDataToBeSigned := device.SecureDataToBeSigned(dataToBeSigned)

	signature, err := device.SignTransaction(securedDataToBeSigned)
	if err != nil {
		return "", "", errors.New(fmt.Sprintf("failed to sign transaction: %s", err))
	}
	encodedSignature := base64.StdEncoding.EncodeToString(signature)

	device.Base64EncodedLastSignature = encodedSignature
	device.SignatureCounter++
	err = deviceRepository.Update(device)
	if err != nil {
		return "", "", errors.New(fmt.Sprintf("failed to update signature device: %s", err))
	}

	return encodedSignature, securedDataToBeSigned, nil
}
