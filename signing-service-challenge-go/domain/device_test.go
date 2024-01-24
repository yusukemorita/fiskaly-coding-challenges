package domain

import (
	"testing"

	"github.com/google/uuid"
)

type MockSignatureAlgorithm struct {
	encodedPrivateKey []byte
}

func (device MockSignatureAlgorithm) Name() string {
	return "MOCK"
}

func (device MockSignatureAlgorithm) GenerateEncodedPrivateKey() ([]byte, error) {
	return device.encodedPrivateKey, nil
}

func TestBuildSignatureDevice(t *testing.T) {
	t.Run("successfully builds signature device", func(t *testing.T) {
		id := uuid.New()
		algorithm := MockSignatureAlgorithm{encodedPrivateKey: []byte("MOCK_KEY")}
		device, err := BuildSignatureDevice(id, algorithm)

		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}

		if device.id != id {
			t.Errorf("expected id: %s, got: %s", id, device.id.String())
		}

		if device.algorithmName != algorithm.Name() {
			t.Errorf("expected algorithm: %s, got: %s", algorithm.Name(), device.algorithmName)
		}

		if device.signatureCounter != 0 {
			t.Errorf("expected initial signature counter value to be 0, got: %d", device.signatureCounter)
		}

		if device.lastSignature != "" {
			t.Errorf("expected initial last signature value to be blank, got: %s", device.lastSignature)
		}

		if string(device.encodedPrivateKey) != string(algorithm.encodedPrivateKey) {
			t.Errorf("expected encoded private key: %s, got: %s", algorithm.encodedPrivateKey, device.lastSignature)
		}

		if device.label != "" {
			t.Errorf("expected label be blank when not provided, got: %s", device.label)
		}
	})

	t.Run("sets label when provided", func(t *testing.T) {
		id := uuid.New()
		algorithm := MockSignatureAlgorithm{encodedPrivateKey: []byte("MOCK_KEY")}
		label := "some-label"
		device, err := BuildSignatureDevice(
			id,
			algorithm,
			"some-label",
		)

		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}

		if device.label != label {
			t.Errorf("expected label: %s, got: %s", label, device.label)
		}
	})
}
