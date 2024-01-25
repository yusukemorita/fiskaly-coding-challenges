package api

import (
	"encoding/json"
	"net/http"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/google/uuid"
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
	Id string `json:"signatureDeviceId"`
}

type CreateSignatureDeviceRequest struct {
	Id        string `json:"id"`
	Algorithm string `json:"algorithm"`
}

func (s *SignatureService) CreateSignatureDevice(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		WriteErrorResponse(response, http.StatusMethodNotAllowed, []string{
			http.StatusText(http.StatusMethodNotAllowed),
		})
		return
	}

	var requestBody CreateSignatureDeviceRequest
	err := json.NewDecoder(request.Body).Decode(&requestBody)
	if err != nil {
		// TODO: better response code?
		WriteErrorResponse(response, http.StatusBadRequest, []string{
			http.StatusText(http.StatusBadRequest),
		})
		return
	}

	id, err := uuid.Parse(requestBody.Id)
	// TODO: handle invalid uuid

	var algorithm domain.SignatureAlgorithm
	for _, alg := range crypto.SupportedAlgorithms {
		if alg.Name() == requestBody.Algorithm {
			algorithm = alg
			break
		}
	}
	// TODO: handle invalid algorithm

	device, err := domain.BuildSignatureDevice(id, algorithm)
	if err != nil {
		// In a real application, this error would be logged and sent to an error notification service
		WriteInternalError(response)
		return
	}

	// TODO: what if the id is duplicate?
	err = s.signatureDeviceRepository.Create(device)
	if err != nil {
		// In a real application, this error would be logged and sent to an error notification service
		WriteInternalError(response)
		return
	}

	responseBody := CreateSignatureDeviceResponse{
		Id: requestBody.Id,
	}
	WriteAPIResponse(response, http.StatusCreated, responseBody)
}
