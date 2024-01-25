package api_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/api"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/persistence"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

func TestCreateSignatureDeviceResponse(t *testing.T) {
	t.Run("fails when method is not POST", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/api/v0/signature_devices", nil)
		responseRecorder := httptest.NewRecorder()

		service := api.NewSignatureService(persistence.NewInMemorySignatureDeviceRepository())
		service.CreateSignatureDevice(responseRecorder, request)

		expectedStatusCode := http.StatusMethodNotAllowed
		if responseRecorder.Code != expectedStatusCode {
			t.Errorf("expected status code: %d, got: %d", expectedStatusCode, responseRecorder.Code)
		}

		body := responseRecorder.Body.String()
		expectedBody := `{"errors":["Method Not Allowed"]}`
		if body != expectedBody {
			t.Errorf("expected: %s, got: %s", expectedBody, body)
		}
	})

	t.Run("fails when uuid is invalid", func(t *testing.T) {
		id := "invalid-uuid"
		algorithm := "RSA"
		request := httptest.NewRequest(
			http.MethodPost,
			"/api/v0/signature_devices",
			strings.NewReader(fmt.Sprintf(`
			{
				"id": "%s",
				"algorithm": "%s"
			}`, id, algorithm)),
		)
		request.Header.Set("Content-Type", "application/json")
		responseRecorder := httptest.NewRecorder()

		repository := persistence.NewInMemorySignatureDeviceRepository()
		service := api.NewSignatureService(repository)
		service.CreateSignatureDevice(responseRecorder, request)

		// check status code
		expectedStatusCode := http.StatusBadRequest
		if responseRecorder.Code != expectedStatusCode {
			t.Errorf("expected status code: %d, got: %d", expectedStatusCode, responseRecorder.Code)
		}

		// check body
		body := responseRecorder.Body.String()
		expectedBody := `{"errors":["id is not a valid uuid"]}`
		if body != expectedBody {
			t.Errorf("expected: %s, got: %s", expectedBody, body)
		}
	})

	t.Run("fails when algorithm is invalid", func(t *testing.T) {
		id := uuid.New()
		algorithm := "ABC"
		request := httptest.NewRequest(
			http.MethodPost,
			"/api/v0/signature_devices",
			strings.NewReader(fmt.Sprintf(`
			{
				"id": "%s",
				"algorithm": "%s"
			}`, id, algorithm)),
		)
		request.Header.Set("Content-Type", "application/json")
		responseRecorder := httptest.NewRecorder()

		repository := persistence.NewInMemorySignatureDeviceRepository()
		service := api.NewSignatureService(repository)
		service.CreateSignatureDevice(responseRecorder, request)

		// check status code
		expectedStatusCode := http.StatusBadRequest
		if responseRecorder.Code != expectedStatusCode {
			t.Errorf("expected status code: %d, got: %d", expectedStatusCode, responseRecorder.Code)
		}

		// check body
		body := responseRecorder.Body.String()
		expectedBody := `{"errors":["algorithm is not supported"]}`
		if body != expectedBody {
			t.Errorf("expected: %s, got: %s", expectedBody, body)
		}
	})

	t.Run("creates a SignatureDevice successfully", func(t *testing.T) {
		id := uuid.New()
		algorithm := "RSA"
		request := httptest.NewRequest(
			http.MethodPost,
			"/api/v0/signature_devices",
			strings.NewReader(fmt.Sprintf(`
			{
				"id": "%s",
				"algorithm": "%s"
			}`, id, algorithm)),
		)
		request.Header.Set("Content-Type", "application/json")
		responseRecorder := httptest.NewRecorder()

		repository := persistence.NewInMemorySignatureDeviceRepository()
		service := api.NewSignatureService(repository)
		service.CreateSignatureDevice(responseRecorder, request)

		// check status code
		expectedStatusCode := http.StatusCreated
		if responseRecorder.Code != expectedStatusCode {
			t.Errorf("expected status code: %d, got: %d", expectedStatusCode, responseRecorder.Code)
		}

		// check body
		body := responseRecorder.Body.String()
		expectedBody := fmt.Sprintf(`{
  "data": {
    "signatureDeviceId": "%s"
  }
}`, id)
		diff := cmp.Diff(body, expectedBody)
		if diff != "" {
			t.Errorf("unexpected diff: %s", diff)
		}

		// check persisted data
		device, found, err := repository.Find(id)
		if err != nil {
			t.Error(err)
		}
		if !found {
			t.Error("expected device with id to be found")
		}
		if device.ID != id {
			t.Errorf("id not persisted correctly. expected: %s, got: %s", id, device.ID)
		}
		if device.AlgorithmName != algorithm {
			t.Errorf("algorithm not persisted correctly. expected: %s, got: %s", algorithm, device.AlgorithmName)
		}
		if device.Label != "" {
			t.Errorf("label not persisted correctly. expected blank string, got: %s", device.Label)
		}
		_, err = crypto.NewRSAMarshaler().Unmarshal(device.EncodedPrivateKey)
		if err != nil {
			t.Errorf("decode of generated private key failed: %s", err)
		}
	})

	t.Run("creates a SignatureDevice with a label successfully", func(t *testing.T) {
		id := uuid.New()
		algorithm := "RSA"
		label := "my RSA key"
		request := httptest.NewRequest(
			http.MethodPost,
			"/api/v0/signature_devices",
			strings.NewReader(fmt.Sprintf(`
			{
				"id": "%s",
				"algorithm": "%s",
				"label": "%s"
			}`, id, algorithm, label)),
		)
		request.Header.Set("Content-Type", "application/json")
		responseRecorder := httptest.NewRecorder()

		repository := persistence.NewInMemorySignatureDeviceRepository()
		service := api.NewSignatureService(repository)
		service.CreateSignatureDevice(responseRecorder, request)

		// check status code
		expectedStatusCode := http.StatusCreated
		if responseRecorder.Code != expectedStatusCode {
			t.Errorf("expected status code: %d, got: %d", expectedStatusCode, responseRecorder.Code)
		}

		// check body
		body := responseRecorder.Body.String()
		expectedBody := fmt.Sprintf(`{
  "data": {
    "signatureDeviceId": "%s"
  }
}`, id)
		diff := cmp.Diff(body, expectedBody)
		if diff != "" {
			t.Errorf("unexpected diff: %s", diff)
		}

		// check persisted data
		device, found, err := repository.Find(id)
		if err != nil {
			t.Error(err)
		}
		if !found {
			t.Error("expected device with id to be found")
		}
		if device.ID != id {
			t.Errorf("id not persisted correctly. expected: %s, got: %s", id, device.ID)
		}
		if device.AlgorithmName != algorithm {
			t.Errorf("algorithm not persisted correctly. expected: %s, got: %s", algorithm, device.AlgorithmName)
		}
		if device.Label != label {
			t.Errorf("label not persisted correctly. expected: %s, got: %s", label, device.Label)
		}
		_, err = crypto.NewRSAMarshaler().Unmarshal(device.EncodedPrivateKey)
		if err != nil {
			t.Errorf("decode of generated private key failed: %s", err)
		}
	})
}
