package api_test

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/api"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/persistence"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

func TestCreateSignatureDeviceResponse(t *testing.T) {
	t.Run("fails when uuid is invalid", func(t *testing.T) {
		id := "invalid-uuid"
		algorithmName := crypto.RSAGenerator{}.AlgorithmName()

		signatureService := api.NewSignatureService(
			persistence.NewInMemorySignatureDeviceRepositoryProvider(
				persistence.NewInMemorySignatureDeviceRepository(),
			),
		)
		server := httptest.NewServer(api.NewServer("", signatureService).HTTPHandler())
		defer server.Close()

		response := sendJsonRequest(
			t,
			http.MethodPost,
			server.URL+"/api/v0/signature_devices",
			api.CreateSignatureDeviceRequest{
				ID:        id,
				Algorithm: algorithmName,
			},
		)

		// check status code
		expectedStatusCode := http.StatusBadRequest
		if response.StatusCode != expectedStatusCode {
			t.Errorf("expected status code: %d, got: %d", expectedStatusCode, response.StatusCode)
		}

		// check body
		body := readBody(t, response)
		expectedBody := `{"errors":["id is not a valid uuid"]}`
		if body != expectedBody {
			t.Errorf("expected: %s, got: %s", expectedBody, body)
		}
	})

	t.Run("fails when id already exists", func(t *testing.T) {
		id := uuid.New()

		// create existing device with the id
		generator := crypto.RSAGenerator{}
		keyPair, err := generator.Generate()
		if err != nil {
			t.Fatal(err)
		}
		repository := persistence.NewInMemorySignatureDeviceRepository()
		repository.Create(domain.SignatureDevice{
			ID:      id,
			KeyPair: keyPair,
		})

		signatureService := api.NewSignatureService(
			persistence.NewInMemorySignatureDeviceRepositoryProvider(repository),
		)
		server := httptest.NewServer(api.NewServer("", signatureService).HTTPHandler())
		defer server.Close()

		response := sendJsonRequest(
			t,
			http.MethodPost,
			server.URL+"/api/v0/signature_devices",
			api.CreateSignatureDeviceRequest{
				ID:        id.String(),
				Algorithm: generator.AlgorithmName(),
			},
		)

		// check status code
		expectedStatusCode := http.StatusBadRequest
		if response.StatusCode != expectedStatusCode {
			t.Errorf("expected status code: %d, got: %d", expectedStatusCode, response.StatusCode)
		}

		// check body
		body := readBody(t, response)
		expectedBody := `{"errors":["duplicate id"]}`
		if body != expectedBody {
			t.Errorf("expected: %s, got: %s", expectedBody, body)
		}
	})

	t.Run("fails when algorithm is invalid", func(t *testing.T) {
		id := uuid.New()
		algorithmName := "ABC"

		signatureService := api.NewSignatureService(
			persistence.NewInMemorySignatureDeviceRepositoryProvider(
				persistence.NewInMemorySignatureDeviceRepository(),
			),
		)
		server := httptest.NewServer(api.NewServer("", signatureService).HTTPHandler())
		defer server.Close()

		response := sendJsonRequest(
			t,
			http.MethodPost,
			server.URL+"/api/v0/signature_devices",
			api.CreateSignatureDeviceRequest{
				ID:        id.String(),
				Algorithm: algorithmName,
			},
		)

		// check status code
		expectedStatusCode := http.StatusBadRequest
		if response.StatusCode != expectedStatusCode {
			t.Errorf("expected status code: %d, got: %d", expectedStatusCode, response.StatusCode)
		}

		// check body
		body := readBody(t, response)
		expectedBody := `{"errors":["algorithm is not supported"]}`
		if body != expectedBody {
			t.Errorf("expected: %s, got: %s", expectedBody, body)
		}
	})

	t.Run("creates a SignatureDevice successfully", func(t *testing.T) {
		id := uuid.New()
		algorithmName := crypto.RSAGenerator{}.AlgorithmName()

		repository := persistence.NewInMemorySignatureDeviceRepository()
		signatureService := api.NewSignatureService(
			persistence.NewInMemorySignatureDeviceRepositoryProvider(repository),
		)
		server := httptest.NewServer(api.NewServer("", signatureService).HTTPHandler())
		defer server.Close()

		response := sendJsonRequest(
			t,
			http.MethodPost,
			server.URL+"/api/v0/signature_devices",
			api.CreateSignatureDeviceRequest{
				ID:        id.String(),
				Algorithm: algorithmName,
			},
		)

		// check status code
		expectedStatusCode := http.StatusCreated
		if response.StatusCode != expectedStatusCode {
			t.Errorf("expected status code: %d, got: %d", expectedStatusCode, response.StatusCode)
		}

		// check body
		createdDevice, ok, err := repository.Find(id)
		if err != nil {
			t.Fatal(err)
		}
		if !ok {
			t.Error("created device not found")
		}
		publicKey, err := createdDevice.KeyPair.EncodedPublicKey()
		if err != nil {
			t.Fatal(err)
		}
		compareResponseBodyData(
			t,
			response,
			api.CreateSignatureDeviceResponse{
				ID:        id.String(),
				Algorithm: algorithmName,
				Label:     "",
				PublicKey: publicKey,
			},
		)
	})

	t.Run("creates a SignatureDevice with a label successfully", func(t *testing.T) {
		id := uuid.New()
		algorithmName := "RSA"
		label := "my RSA key"

		repository := persistence.NewInMemorySignatureDeviceRepository()
		signatureService := api.NewSignatureService(
			persistence.NewInMemorySignatureDeviceRepositoryProvider(repository),
		)
		server := httptest.NewServer(api.NewServer("", signatureService).HTTPHandler())
		defer server.Close()

		response := sendJsonRequest(
			t,
			http.MethodPost,
			server.URL+"/api/v0/signature_devices",
			api.CreateSignatureDeviceRequest{
				ID:        id.String(),
				Algorithm: algorithmName,
				Label:     label,
			},
		)

		// check status code
		expectedStatusCode := http.StatusCreated
		if response.StatusCode != expectedStatusCode {
			t.Errorf("expected status code: %d, got: %d", expectedStatusCode, response.StatusCode)
		}

		// check body
		createdDevice, ok, err := repository.Find(id)
		if err != nil {
			t.Fatal(err)
		}
		if !ok {
			t.Error("created device not found")
		}
		publicKey, err := createdDevice.KeyPair.EncodedPublicKey()
		if err != nil {
			t.Fatal(err)
		}
		compareResponseBodyData(
			t,
			response,
			api.CreateSignatureDeviceResponse{
				ID:        id.String(),
				Algorithm: algorithmName,
				Label:     label,
				PublicKey: publicKey,
			},
		)
	})
}

func TestSignTransaction(t *testing.T) {
	t.Run("returns not found when device with id does not exist", func(t *testing.T) {
		id := uuid.NewString()

		signatureService := api.NewSignatureService(
			persistence.NewInMemorySignatureDeviceRepositoryProvider(
				persistence.NewInMemorySignatureDeviceRepository(),
			),
		)
		testServer := httptest.NewServer(api.NewServer(":8888", signatureService).HTTPHandler())
		defer testServer.Close()

		response := sendJsonRequest(
			t,
			http.MethodPost,
			fmt.Sprintf("%s/api/v0/signature_devices/%s/signatures", testServer.URL, id),
			api.SignTransactionRequest{DataToBeSigned: "some-data"},
		)

		// check status code
		expectedStatusCode := http.StatusNotFound
		if response.StatusCode != expectedStatusCode {
			t.Errorf("expected status code: %d, got: %d", expectedStatusCode, response.StatusCode)
		}

		// check body
		body := readBody(t, response)
		expectedBody := `{"errors":["signature device not found"]}`
		diff := cmp.Diff(body, expectedBody)
		if diff != "" {
			t.Errorf("unexpected diff: %s", diff)
		}
	})

	t.Run("successfully signs data with device (algorithm: RSA, counter = 0)", func(t *testing.T) {
		id := "64ff796e-fcde-499a-a03d-82dd1f89e8e5"
		base64EncodedID := "NjRmZjc5NmUtZmNkZS00OTlhLWEwM2QtODJkZDFmODllOGU1"
		dataToSign := "some-data"
		device, err := domain.BuildSignatureDevice(uuid.MustParse(id), crypto.RSAGenerator{})
		if err != nil {
			t.Fatal(err)
		}

		repository := persistence.NewInMemorySignatureDeviceRepository()
		err = repository.Create(device)
		if err != nil {
			t.Fatal(err)
		}

		signatureService := api.NewSignatureService(
			persistence.NewInMemorySignatureDeviceRepositoryProvider(repository),
		)
		testServer := httptest.NewServer(api.NewServer(":8888", signatureService).HTTPHandler())
		defer testServer.Close()

		response := sendJsonRequest(
			t,
			http.MethodPost,
			fmt.Sprintf("%s/api/v0/signature_devices/%s/signatures", testServer.URL, id),
			api.SignTransactionRequest{DataToBeSigned: dataToSign},
		)

		// check status code
		expectedStatusCode := http.StatusOK
		if response.StatusCode != expectedStatusCode {
			t.Errorf("expected status code: %d, got: %d", expectedStatusCode, response.StatusCode)
		}

		// unmarshal body
		body := readBody(t, response)
		jsonBody := struct {
			Data api.SignTransactionResponse `json:"data"`
		}{}
		err = json.Unmarshal([]byte(body), &jsonBody)
		if err != nil {
			t.Errorf("unexpected response body format: %s", err)
		}

		// check signature is verifiable
		keyPair := device.KeyPair.(*crypto.RSAKeyPair)
		digest, err := crypto.ComputeHashDigest([]byte(jsonBody.Data.SignedData))
		if err != nil {
			t.Fatal(err)
		}
		decodedSignature, err := base64.StdEncoding.DecodeString(jsonBody.Data.Signature)
		if err != nil {
			t.Fatal(err)
		}
		err = rsa.VerifyPSS(keyPair.Public, crypto.HashFunction, digest, decodedSignature, nil)
		if err != nil {
			t.Errorf("verification of signed data and signature failed. err: %s, signed data: %s, signature: %s", err, jsonBody.Data.SignedData, jsonBody.Data.Signature)
		}

		// check signed_data is correct format
		expectedSignedData := fmt.Sprintf("0_%s_%s", dataToSign, base64EncodedID)
		if jsonBody.Data.SignedData != expectedSignedData {
			t.Errorf("expected signed data: %s, got: %s", expectedSignedData, jsonBody.Data.SignedData)
		}

		// check persisted data
		device, ok, err := repository.Find(uuid.MustParse(id))
		if err != nil {
			t.Fatal(err)
		}
		if !ok {
			t.Fatal("device not found")
		}
		if device.SignatureCounter != 1 {
			t.Errorf("device signature counter should be incremented to 1, got: %d", device.SignatureCounter)
		}
		if device.LastSignature != jsonBody.Data.Signature {
			t.Errorf("device last signature should be updated to %s, got: %s", jsonBody.Data.Signature, device.LastSignature)
		}
	})

	t.Run("successfully signs data (algorithm: RSA, counter > 0)", func(t *testing.T) {
		id := "64ff796e-fcde-499a-a03d-82dd1f89e8e5"
		dataToSign := "some-data"

		// create a device that has been used once
		device, err := domain.BuildSignatureDevice(uuid.MustParse(id), crypto.RSAGenerator{})
		if err != nil {
			t.Fatal(err)
		}
		device.SignatureCounter = 1
		device.LastSignature = "last-signature-base-64-encoded"
		repository := persistence.NewInMemorySignatureDeviceRepository()
		err = repository.Create(device)
		if err != nil {
			t.Fatal(err)
		}

		signatureService := api.NewSignatureService(
			persistence.NewInMemorySignatureDeviceRepositoryProvider(repository),
		)
		testServer := httptest.NewServer(api.NewServer(":8888", signatureService).HTTPHandler())
		defer testServer.Close()

		response := sendJsonRequest(
			t,
			http.MethodPost,
			fmt.Sprintf("%s/api/v0/signature_devices/%s/signatures", testServer.URL, id),
			api.SignTransactionRequest{DataToBeSigned: dataToSign},
		)

		// check status code
		expectedStatusCode := http.StatusOK
		if response.StatusCode != expectedStatusCode {
			t.Errorf("expected status code: %d, got: %d", expectedStatusCode, response.StatusCode)
		}

		// unmarshal body
		body := readBody(t, response)
		jsonBody := struct {
			Data api.SignTransactionResponse `json:"data"`
		}{}
		err = json.Unmarshal([]byte(body), &jsonBody)
		if err != nil {
			t.Errorf("unexpected response body format: %s", err)
		}

		// check signature is verifiable
		keyPair := device.KeyPair.(*crypto.RSAKeyPair)
		digest, err := crypto.ComputeHashDigest([]byte(jsonBody.Data.SignedData))
		if err != nil {
			t.Fatal(err)
		}
		decodedSignature, err := base64.StdEncoding.DecodeString(jsonBody.Data.Signature)
		if err != nil {
			t.Fatal(err)
		}
		err = rsa.VerifyPSS(keyPair.Public, crypto.HashFunction, digest, decodedSignature, nil)
		if err != nil {
			t.Errorf("verification of signed data and signature failed. err: %s, signed data: %s, signature: %s", err, jsonBody.Data.SignedData, jsonBody.Data.Signature)
		}

		// check signed_data is correct format
		expectedSignedData := fmt.Sprintf("1_%s_%s", dataToSign, device.LastSignature)
		if jsonBody.Data.SignedData != expectedSignedData {
			t.Errorf("expected signed data: %s, got: %s", expectedSignedData, jsonBody.Data.SignedData)
		}

		// check persisted data
		device, ok, err := repository.Find(uuid.MustParse(id))
		if err != nil {
			t.Fatal(err)
		}
		if !ok {
			t.Fatal("device not found")
		}
		if device.SignatureCounter != 2 {
			t.Errorf("device signature counter should be incremented to 2, got: %d", device.SignatureCounter)
		}
		if device.LastSignature != jsonBody.Data.Signature {
			t.Errorf("device last signature should be updated to %s, got: %s", jsonBody.Data.Signature, device.LastSignature)
		}
	})

	t.Run("successfully signs data with device (algorithm: ECC, counter = 0)", func(t *testing.T) {
		id := "64ff796e-fcde-499a-a03d-82dd1f89e8e5"
		base64EncodedID := "NjRmZjc5NmUtZmNkZS00OTlhLWEwM2QtODJkZDFmODllOGU1"
		dataToSign := "some-data"
		device, err := domain.BuildSignatureDevice(uuid.MustParse(id), crypto.ECCGenerator{})
		if err != nil {
			t.Fatal(err)
		}

		repository := persistence.NewInMemorySignatureDeviceRepository()
		err = repository.Create(device)
		if err != nil {
			t.Fatal(err)
		}

		signatureService := api.NewSignatureService(
			persistence.NewInMemorySignatureDeviceRepositoryProvider(repository),
		)
		testServer := httptest.NewServer(api.NewServer(":8888", signatureService).HTTPHandler())
		defer testServer.Close()

		response := sendJsonRequest(
			t,
			http.MethodPost,
			fmt.Sprintf("%s/api/v0/signature_devices/%s/signatures", testServer.URL, id),
			api.SignTransactionRequest{DataToBeSigned: dataToSign},
		)

		// check status code
		expectedStatusCode := http.StatusOK
		if response.StatusCode != expectedStatusCode {
			t.Errorf("expected status code: %d, got: %d", expectedStatusCode, response.StatusCode)
		}

		// unmarshal body
		body := readBody(t, response)
		jsonBody := struct {
			Data api.SignTransactionResponse `json:"data"`
		}{}
		err = json.Unmarshal([]byte(body), &jsonBody)
		if err != nil {
			t.Errorf("unexpected response body format: %s", err)
		}

		// check signature is verifiable
		keyPair := device.KeyPair.(*crypto.ECCKeyPair)
		digest, err := crypto.ComputeHashDigest([]byte(jsonBody.Data.SignedData))
		if err != nil {
			t.Fatal(err)
		}
		decodedSignature, err := base64.StdEncoding.DecodeString(jsonBody.Data.Signature)
		if err != nil {
			t.Fatal(err)
		}
		result := ecdsa.VerifyASN1(keyPair.Public, digest, decodedSignature)
		if !result {
			t.Errorf("verification of signed data and signature failed. err: %s, signed data: %s, signature: %s", err, jsonBody.Data.SignedData, jsonBody.Data.Signature)
		}

		// check signed_data is correct format
		expectedSignedData := fmt.Sprintf("0_%s_%s", dataToSign, base64EncodedID)
		if jsonBody.Data.SignedData != expectedSignedData {
			t.Errorf("expected signed data: %s, got: %s", expectedSignedData, jsonBody.Data.SignedData)
		}

		// check persisted data
		device, ok, err := repository.Find(uuid.MustParse(id))
		if err != nil {
			t.Fatal(err)
		}
		if !ok {
			t.Fatal("device not found")
		}
		if device.SignatureCounter != 1 {
			t.Errorf("device signature counter should be incremented to 1, got: %d", device.SignatureCounter)
		}
		if device.LastSignature != jsonBody.Data.Signature {
			t.Errorf("device last signature should be updated to %s, got: %s", jsonBody.Data.Signature, device.LastSignature)
		}
	})

	t.Run("successfully signs data (algorithm: ECC, counter > 0)", func(t *testing.T) {
		id := "64ff796e-fcde-499a-a03d-82dd1f89e8e5"
		dataToSign := "some-data"

		// create a device that has been used once
		device, err := domain.BuildSignatureDevice(uuid.MustParse(id), crypto.ECCGenerator{})
		if err != nil {
			t.Fatal(err)
		}
		device.SignatureCounter = 1
		device.LastSignature = "last-signature-base-64-encoded"
		repository := persistence.NewInMemorySignatureDeviceRepository()
		err = repository.Create(device)
		if err != nil {
			t.Fatal(err)
		}

		signatureService := api.NewSignatureService(
			persistence.NewInMemorySignatureDeviceRepositoryProvider(repository),
		)
		testServer := httptest.NewServer(api.NewServer(":8888", signatureService).HTTPHandler())
		defer testServer.Close()

		response := sendJsonRequest(
			t,
			http.MethodPost,
			fmt.Sprintf("%s/api/v0/signature_devices/%s/signatures", testServer.URL, id),
			api.SignTransactionRequest{DataToBeSigned: dataToSign},
		)

		// check status code
		expectedStatusCode := http.StatusOK
		if response.StatusCode != expectedStatusCode {
			t.Errorf("expected status code: %d, got: %d", expectedStatusCode, response.StatusCode)
		}

		// unmarshal body
		body := readBody(t, response)
		jsonBody := struct {
			Data api.SignTransactionResponse `json:"data"`
		}{}
		err = json.Unmarshal([]byte(body), &jsonBody)
		if err != nil {
			t.Errorf("unexpected response body format: %s", err)
		}

		// check signature is verifiable
		digest, err := crypto.ComputeHashDigest([]byte(jsonBody.Data.SignedData))
		if err != nil {
			t.Fatal(err)
		}
		decodedSignature, err := base64.StdEncoding.DecodeString(jsonBody.Data.Signature)
		if err != nil {
			t.Fatal(err)
		}
		keyPair := device.KeyPair.(*crypto.ECCKeyPair)
		result := ecdsa.VerifyASN1(keyPair.Public, digest, decodedSignature)
		if !result {
			t.Errorf("verification of signed data and signature failed. err: %s, signed data: %s, signature: %s", err, jsonBody.Data.SignedData, jsonBody.Data.Signature)
		}

		// check signed_data is correct format
		expectedSignedData := fmt.Sprintf("1_%s_%s", dataToSign, device.LastSignature)
		if jsonBody.Data.SignedData != expectedSignedData {
			t.Errorf("expected signed data: %s, got: %s", expectedSignedData, jsonBody.Data.SignedData)
		}

		// check persisted data
		device, ok, err := repository.Find(uuid.MustParse(id))
		if err != nil {
			t.Fatal(err)
		}
		if !ok {
			t.Fatal("device not found")
		}
		if device.SignatureCounter != 2 {
			t.Errorf("device signature counter should be incremented to 2, got: %d", device.SignatureCounter)
		}
		if device.LastSignature != jsonBody.Data.Signature {
			t.Errorf("device last signature should be updated to %s, got: %s", jsonBody.Data.Signature, device.LastSignature)
		}
	})
}

func TestFindSignatureDevice(t *testing.T) {
	t.Run("returns not found when device with id does not exist", func(t *testing.T) {
		id := uuid.NewString()

		repository := persistence.NewInMemorySignatureDeviceRepository()
		signatureService := api.NewSignatureService(
			persistence.NewInMemorySignatureDeviceRepositoryProvider(repository),
		)
		testServer := httptest.NewServer(api.NewServer("", signatureService).HTTPHandler())
		defer testServer.Close()

		response := sendJsonRequest(
			t,
			http.MethodGet,
			fmt.Sprintf("%s/api/v0/signature_devices/%s", testServer.URL, id),
		)

		// check status code
		expectedStatusCode := http.StatusNotFound
		if response.StatusCode != expectedStatusCode {
			t.Errorf("expected status code: %d, got: %d", expectedStatusCode, response.StatusCode)
		}

		// check body
		body := readBody(t, response)
		expectedBody := `{"errors":["signature device not found"]}`
		diff := cmp.Diff(body, expectedBody)
		if diff != "" {
			t.Errorf("unexpected diff: %s", diff)
		}
	})

	t.Run("returns device when device with id exists", func(t *testing.T) {
		// create a device
		label := "my ecc key"
		device, err := domain.BuildSignatureDevice(
			uuid.New(),
			crypto.ECCGenerator{},
			label,
		)
		if err != nil {
			t.Fatal(err)
		}
		repository := persistence.NewInMemorySignatureDeviceRepository()
		err = repository.Create(device)
		if err != nil {
			t.Fatal(err)
		}

		signatureService := api.NewSignatureService(
			persistence.NewInMemorySignatureDeviceRepositoryProvider(repository),
		)
		testServer := httptest.NewServer(api.NewServer("", signatureService).HTTPHandler())
		defer testServer.Close()

		response := sendJsonRequest(
			t,
			http.MethodGet,
			fmt.Sprintf("%s/api/v0/signature_devices/%s", testServer.URL, device.ID.String()),
		)

		// check status code
		expectedStatusCode := http.StatusOK
		if response.StatusCode != expectedStatusCode {
			t.Errorf("expected status code: %d, got: %d", expectedStatusCode, response.StatusCode)
		}

		// check body
		publicKey, err := device.KeyPair.EncodedPublicKey()
		if err != nil {
			t.Fatal(err)
		}
		compareResponseBodyData(
			t,
			response,
			api.FindSignatureDeviceResponse{
				ID:        device.ID.String(),
				Label:     label,
				PublicKey: publicKey,
				Algorithm: "ECC",
			},
		)
	})
}

func TestListSignatureDevices(t *testing.T) {
	// create an ecc device
	eccDevice, err := domain.BuildSignatureDevice(
		uuid.MustParse("e9af3524-a2ab-4671-b3c5-d7fdc10b511e"),
		crypto.ECCGenerator{},
		"my ecc key",
	)
	if err != nil {
		t.Fatal(err)
	}
	repository := persistence.NewInMemorySignatureDeviceRepository()
	err = repository.Create(eccDevice)
	if err != nil {
		t.Fatal(err)
	}

	// create an rsa device
	rsaDevice, err := domain.BuildSignatureDevice(
		uuid.MustParse("f7e820e9-7bf5-4c41-a463-b038cb9336a0"),
		crypto.RSAGenerator{},
		"my rsa key",
	)
	if err != nil {
		t.Fatal(err)
	}
	err = repository.Create(rsaDevice)
	if err != nil {
		t.Fatal(err)
	}

	signatureService := api.NewSignatureService(
		persistence.NewInMemorySignatureDeviceRepositoryProvider(repository),
	)
	testServer := httptest.NewServer(api.NewServer("", signatureService).HTTPHandler())
	defer testServer.Close()

	response := sendJsonRequest(
		t,
		http.MethodGet,
		fmt.Sprintf("%s/api/v0/signature_devices", testServer.URL),
	)

	// check status code
	expectedStatusCode := http.StatusOK
	if response.StatusCode != expectedStatusCode {
		t.Errorf("expected status code: %d, got: %d", expectedStatusCode, response.StatusCode)
	}

	// check body
	rsaPublicKey, err := rsaDevice.KeyPair.EncodedPublicKey()
	if err != nil {
		t.Fatal(err)
	}
	eccPublicKey, err := eccDevice.KeyPair.EncodedPublicKey()
	if err != nil {
		t.Fatal(err)
	}
	expectedBody := api.ListSignatureDevicesResponse{
		// ecc device id must come first, as the ids should be sorted alphabetically
		{
			ID:        eccDevice.ID.String(),
			Label:     eccDevice.Label,
			Algorithm: eccDevice.KeyPair.AlgorithmName(),
			PublicKey: eccPublicKey,
		},
		{
			ID:        rsaDevice.ID.String(),
			Label:     rsaDevice.Label,
			Algorithm: rsaDevice.KeyPair.AlgorithmName(),
			PublicKey: rsaPublicKey,
		},
	}
	compareResponseBodyData(t, response, expectedBody)
}

func compareResponseBodyData(t *testing.T, response *http.Response, expectedData any) {
	t.Helper()

	body := readBody(t, response)
	expectedResponse := struct {
		Data any `json:"data"`
	}{
		Data: expectedData,
	}
	expectedBody, err := json.MarshalIndent(expectedResponse, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	diff := cmp.Diff(body, string(expectedBody))
	if diff != "" {
		t.Errorf("unexpected diff: %s", diff)
	}
}
