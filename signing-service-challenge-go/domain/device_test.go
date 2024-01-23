package domain

import (
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestBuildSignatureDevice(t *testing.T) {
	t.Run("returns error when uuid is invalid", func(t *testing.T) {
		id := "invalid-value"
		_, err := BuildSignatureDevice(id, "RSA")

		if err == nil {
			t.Errorf("expected error when uuid is invalid, got nil")
		}

		if !strings.Contains(err.Error(), "invalid uuid:") {
			t.Errorf("expected invalid uuid error, got: %s", err)
		}
	})

	t.Run("returns error when algorithm is invalid", func(t *testing.T) {
		_, err := BuildSignatureDevice(uuid.NewString(), "ABC")

		if err == nil || err.Error() != "invalid algorithm" {
			t.Errorf("expected error: invalid algorithm, got: %s", err)
		}
	})

	t.Run("successfully builds RSA signature device", func(t *testing.T) {
		id := uuid.NewString()
		algorithmName := "RSA"
		device, err := BuildSignatureDevice(
			id,
			algorithmName,
		)

		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}

		if device.id.String() != id {
			t.Errorf("expected id: %s, got: %s", id, device.id.String())
		}

		if device.algorithmName != algorithmName {
			t.Errorf("expected algorithm: %s, got: %s", algorithmName, device.algorithmName)
		}

		if device.signatureCounter != 0 {
			t.Errorf("expected initial signature counter value to be 0, got: %d", device.signatureCounter)
		}

		if device.lastSignature != "" {
			t.Errorf("expected initial last signature value to be blank, got: %s", device.lastSignature)
		}
	})

	t.Run("successfully builds ECC signature device", func(t *testing.T) {
		id := uuid.NewString()
		algorithmName := "ECC"
		device, err := BuildSignatureDevice(
			id,
			algorithmName,
		)

		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}

		if device.id.String() != id {
			t.Errorf("expected id: %s, got: %s", id, device.id.String())
		}

		if device.algorithmName != algorithmName {
			t.Errorf("expected algorithm: %s, got: %s", algorithmName, device.algorithmName)
		}

		if device.signatureCounter != 0 {
			t.Errorf("expected initial signature counter value to be 0, got: %d", device.signatureCounter)
		}

		if device.lastSignature != "" {
			t.Errorf("expected initial last signature value to be blank, got: %s", device.lastSignature)
		}
	})

	t.Run("sets label when provided", func(t *testing.T) {
		id := uuid.NewString()
		algorithm := "RSA"
		label := "some-label"
		device, err := BuildSignatureDevice(
			id,
			string(algorithm),
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
