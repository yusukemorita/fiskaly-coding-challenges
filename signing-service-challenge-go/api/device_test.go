package api_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/api"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/persistence"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

func TestCreateSignatureDeviceResponse(t *testing.T) {
	t.Run("returns MethodNotAllowed when not POST", func(t *testing.T) {
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

		// TODO: check that encoded key is valid
	})
}
