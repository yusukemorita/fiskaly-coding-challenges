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
	repositoryProvider domain.SignatureDeviceRepositoryProvider
}

func NewSignatureService(p domain.SignatureDeviceRepositoryProvider) SignatureService {
	return SignatureService{
		repositoryProvider: p,
	}
}

type ApiSignatureDevice struct {
	ID               string `json:"id"`
	Label            string `json:"label"`
	PublicKey        string `json:"public_key"`
	Algorithm        string `json:"algorithm"`
	SignatureCounter uint   `json:"signature_counter"`
	LastSignature    string `json:"last_signature"`
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

	var idIsDuplicate bool
	var algorithmIsInvalid bool
	var device domain.SignatureDevice
	err = s.repositoryProvider.WriteTx(func(repository domain.SignatureDeviceRepository) error {
		_, ok, err := repository.Find(id)
		if err != nil {
			return err
		}
		if ok {
			idIsDuplicate = true
			return nil
		}

		generator, found := crypto.FindKeyPairGenerator(requestBody.Algorithm)
		if !found {
			algorithmIsInvalid = true
			return nil
		}

		device, err = domain.BuildSignatureDevice(id, generator, requestBody.Label)
		if err != nil {
			return err
		}

		if err := repository.Create(device); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		WriteInternalError(response)
		return
	}
	if idIsDuplicate {
		WriteErrorResponse(response, http.StatusBadRequest, []string{
			"duplicate id",
		})
		return
	}
	if algorithmIsInvalid {
		WriteErrorResponse(response, http.StatusBadRequest, []string{
			"algorithm is not supported",
		})
		return
	}

	publicKey, err := device.KeyPair.EncodedPublicKey()
	if err != nil {
		WriteInternalError(response)
		return
	}

	responseBody := CreateSignatureDeviceResponse{
		ID:               device.ID.String(),
		Label:            device.Label,
		PublicKey:        publicKey,
		Algorithm:        device.KeyPair.AlgorithmName(),
		SignatureCounter: device.SignatureCounter,
		LastSignature:    device.LastSignature,
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
		s.repositoryProvider,
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

	var device domain.SignatureDevice
	var deviceFound bool
	err = s.repositoryProvider.ReadTx(func(repository domain.SignatureDeviceRepository) error {
		device, deviceFound, err = repository.Find(deviceID)
		return err
	})
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

	publicKey, err := device.KeyPair.EncodedPublicKey()
	if err != nil {
		WriteInternalError(response)
		return
	}

	WriteAPIResponse(
		response,
		http.StatusOK,
		FindSignatureDeviceResponse{
			ID:               device.ID.String(),
			Label:            device.Label,
			PublicKey:        publicKey,
			Algorithm:        device.KeyPair.AlgorithmName(),
			SignatureCounter: device.SignatureCounter,
			LastSignature:    device.LastSignature,
		},
	)
}

type ListSignatureDevicesResponse = []ApiSignatureDevice

func (s *SignatureService) ListSignatureDevice(response http.ResponseWriter, request *http.Request) {
	var devices []domain.SignatureDevice
	err := s.repositoryProvider.ReadTx(func(repository domain.SignatureDeviceRepository) error {
		d, err := repository.List()
		if err != nil {
			return err
		}
		devices = d
		return nil
	})
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
			ID:               device.ID.String(),
			Label:            device.Label,
			Algorithm:        device.KeyPair.AlgorithmName(),
			PublicKey:        publicKey,
			SignatureCounter: device.SignatureCounter,
			LastSignature:    device.LastSignature,
		})
	}

	WriteAPIResponse(response, http.StatusOK, responseBody)
}
