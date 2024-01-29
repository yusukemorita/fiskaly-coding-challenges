package api

import (
	"encoding/json"
	"net/http"
	"sort"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type SignatureService struct {
	signatureDeviceRepository domain.SignatureDeviceRepository
}

func NewSignatureService(repository domain.SignatureDeviceRepository) SignatureService {
	return SignatureService{
		signatureDeviceRepository: repository,
	}
}

type ApiSignatureDevice struct {
	ID        string `json:"id"`
	Label     string `json:"label"`
	PublicKey string `json:"public_key"`
	Algorithm string `json:"algorithm"`
}

type CreateSignatureDeviceResponse = ApiSignatureDevice

type CreateSignatureDeviceRequest struct {
	ID        string `json:"id"`
	Algorithm string `json:"algorithm"`
	Label     string `json:"label"` // optional
}

func (s *SignatureService) CreateSignatureDevice(response http.ResponseWriter, request *http.Request) {
	var requestBody CreateSignatureDeviceRequest
	err := json.NewDecoder(request.Body).Decode(&requestBody)
	if err != nil {
		WriteErrorResponse(response, http.StatusBadRequest, []string{
			"invalid json",
		})
		return
	}

	id, err := uuid.Parse(requestBody.ID)
	if err != nil {
		WriteErrorResponse(response, http.StatusBadRequest, []string{
			"id is not a valid uuid",
		})
		return
	}

	mutex := s.signatureDeviceRepository.Mutex()
	// acquire a write lock here, as there is a write lock later on in
	// this function
	mutex.Lock()
	defer mutex.Unlock()
	_, ok, err := s.signatureDeviceRepository.Find(id)
	if err != nil {
		WriteInternalError(response)
		return
	}
	if ok {
		WriteErrorResponse(response, http.StatusBadRequest, []string{
			"duplicate id",
		})
		return
	}

	generator, found := crypto.FindKeyPairGenerator(requestBody.Algorithm)
	if !found {
		WriteErrorResponse(response, http.StatusBadRequest, []string{
			"algorithm is not supported",
		})
		return
	}

	device, err := domain.BuildSignatureDevice(id, generator, requestBody.Label)
	if err != nil {
		// In a real application, this error would be logged and sent to an error notification service
		WriteInternalError(response)
		return
	}

	err = s.signatureDeviceRepository.Create(device)
	if err != nil {
		// In a real application, this error would be logged and sent to an error notification service
		WriteInternalError(response)
		return
	}

	publicKey, err := device.KeyPair.EncodedPublicKey()
	if err != nil {
		WriteInternalError(response)
		return
	}

	responseBody := CreateSignatureDeviceResponse{
		ID:        device.ID.String(),
		Label:     device.Label,
		PublicKey: publicKey,
		Algorithm: device.KeyPair.AlgorithmName(),
	}
	WriteAPIResponse(response, http.StatusCreated, responseBody)
}

type SignTransactionRequest struct {
	DataToBeSigned string `json:"data_to_be_signed"`
}

type SignTransactionResponse struct {
	Signature  string `json:"signature"`
	SignedData string `json:"signed_data"`
}

func (s *SignatureService) SignTransaction(response http.ResponseWriter, request *http.Request) {
	deviceIDString := chi.URLParam(request, "deviceID")
	deviceID, err := uuid.Parse(deviceIDString)
	if err != nil {
		WriteErrorResponse(response, http.StatusBadRequest, []string{
			"id is not a valid uuid",
		})
		return
	}

	var requestBody SignTransactionRequest
	err = json.NewDecoder(request.Body).Decode(&requestBody)
	if err != nil {
		WriteErrorResponse(response, http.StatusBadRequest, []string{
			"invalid json",
		})
		return
	}

	deviceFound, encodedSignature, signedData, err := domain.SignTransaction(
		deviceID,
		s.signatureDeviceRepository,
		requestBody.DataToBeSigned,
	)
	if err != nil {
		WriteInternalError(response)
		return
	}
	if !deviceFound {
		WriteErrorResponse(response, http.StatusNotFound, []string{
			"signature device not found",
		})
		return
	}

	WriteAPIResponse(
		response,
		http.StatusOK,
		SignTransactionResponse{
			Signature:  encodedSignature,
			SignedData: signedData,
		},
	)
}

type FindSignatureDeviceResponse = ApiSignatureDevice

func (s *SignatureService) FindSignatureDevice(response http.ResponseWriter, request *http.Request) {
	deviceIDString := chi.URLParam(request, "deviceID")
	deviceID, err := uuid.Parse(deviceIDString)
	if err != nil {
		WriteErrorResponse(response, http.StatusBadRequest, []string{
			"id is not a valid uuid",
		})
		return
	}

	mutex := s.signatureDeviceRepository.Mutex()
	mutex.RLock()
	defer mutex.RUnlock()
	device, ok, err := s.signatureDeviceRepository.Find(deviceID)
	if err != nil {
		WriteInternalError(response)
		return
	}
	if !ok {
		WriteErrorResponse(response, http.StatusNotFound, []string{
			"signature device not found",
		})
		return
	}

	publicKey, err := device.KeyPair.EncodedPublicKey()
	if err != nil {
		WriteInternalError(response)
		return
	}

	WriteAPIResponse(
		response,
		http.StatusOK,
		FindSignatureDeviceResponse{
			ID:        device.ID.String(),
			Label:     device.Label,
			PublicKey: publicKey,
			Algorithm: device.KeyPair.AlgorithmName(),
		},
	)
}

type ListSignatureDevicesResponse = []ApiSignatureDevice

func (s *SignatureService) ListSignatureDevice(response http.ResponseWriter, request *http.Request) {
	mutex := s.signatureDeviceRepository.Mutex()
	mutex.RLock()
	defer mutex.RUnlock()
	devices, err := s.signatureDeviceRepository.List()
	if err != nil {
		WriteInternalError(response)
		return
	}

	// sort devices by ID
	// This is just so that the order remains stable, what the order is determined by
	// is irrelevant
	sort.SliceStable(devices, func(a, b int) bool {
		deviceA := devices[a]
		deviceB := devices[b]
		return deviceA.ID.String() < deviceB.ID.String()
	})

	responseBody := ListSignatureDevicesResponse{}
	for _, device := range devices {
		publicKey, err := device.KeyPair.EncodedPublicKey()
		if err != nil {
			WriteInternalError(response)
			return
		}
		responseBody = append(responseBody, ApiSignatureDevice{
			ID:        device.ID.String(),
			Label:     device.Label,
			Algorithm: device.KeyPair.AlgorithmName(),
			PublicKey: publicKey,
		})
	}

	WriteAPIResponse(response, http.StatusOK, responseBody)
}
