package domain

import (
	"testing"

	"github.com/google/uuid"
)

func TestBuildSignatureDevice(t *testing.T) {
	t.Run("successfully builds RSA signature device", func(t *testing.T) {
		id := uuid.New()
		algorithm := rsaAlgorithm{}
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
	})

	t.Run("successfully builds ECC signature device", func(t *testing.T) {
		id := uuid.New()
		algorithm := eccAlgorithm{}
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
	})

	t.Run("sets label when provided", func(t *testing.T) {
		id := uuid.New()
		algorithm := rsaAlgorithm{}
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
