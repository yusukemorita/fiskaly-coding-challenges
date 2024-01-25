package api

import (
	"net/http"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
)

type SignatureService struct {
	signatureDeviceRepository domain.SignatureDeviceRepository
}

func NewSignatureService(signatureDeviceRepository domain.SignatureDeviceRepository) SignatureService {
	return SignatureService{
		signatureDeviceRepository: signatureDeviceRepository,
	}
}

// TODO: REST endpoints ...
type CreateSignatureDeviceResponse struct {
	signatureDeviceId string
}

func (s *SignatureService) CreateSignatureDevice(response http.ResponseWriter, request *http.Request) {
	WriteAPIResponse(response, http.StatusOK, "")
}
