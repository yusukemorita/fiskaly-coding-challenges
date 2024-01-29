package persistence

import (
	"testing"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

func TestCreate(t *testing.T) {
	t.Run("persists the device in memory", func(t *testing.T) {
		keyPair, err := crypto.RSAGenerator{}.Generate()
		if err != nil {
			t.Fatal(err)
		}
		device := domain.SignatureDevice{
			ID:      uuid.New(),
			KeyPair: keyPair,
			Label:   "my rsa key",
		}

		repository := NewInMemorySignatureDeviceRepository()

		if len(repository.devices) != 0 {
			t.Errorf("new repository should have 0 devices")
		}

		err = repository.Create(device)
		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}
		if len(repository.devices) != 1 {
			t.Errorf("expected repository to contain 1 device, got: %d", len(repository.devices))
		}

		persistedDevice, ok := repository.devices[device.ID]
		if !ok {
			t.Error("expected device with id to be persisted")
		}
		diff := cmp.Diff(persistedDevice, device)
		if diff != "" {
			t.Errorf("unexpected difference between original and persisted device: %s", diff)
		}
	})

	t.Run("does not persist when id is not unique", func(t *testing.T) {
		id := uuid.New()
		alreadyExistingDevice := domain.SignatureDevice{
			ID:    id,
			Label: "already existing rsa key",
		}
		duplicateIdDevice := domain.SignatureDevice{
			ID:    id,
			Label: "new rsa key",
		}

		repository := NewInMemorySignatureDeviceRepository()
		repository.devices[id] = alreadyExistingDevice
		if len(repository.devices) != 1 {
			t.Errorf("repository should contain 1 device")
		}

		err := repository.Create(duplicateIdDevice)
		if err == nil {
			t.Error("expected error")
		}
		if len(repository.devices) != 1 {
			t.Errorf("expected repository to contain 1 device, got: %d", len(repository.devices))
		}

		persistedDevice, ok := repository.devices[id]
		if !ok {
			t.Error("expected device with id to be present")
		}
		diff := cmp.Diff(persistedDevice, alreadyExistingDevice)
		if diff != "" {
			t.Errorf("expected persisted device to not have changed. diff: %s", diff)
		}
	})
}

func TestMarkSignatureCreated(t *testing.T) {
	t.Run("increments counter and updates last signature when device with id is found", func(t *testing.T) {
		id := uuid.New()
		keyPair, err := crypto.ECCGenerator{}.Generate()
		if err != nil {
			t.Fatal(err)
		}
		device := domain.SignatureDevice{
			ID:      id,
			KeyPair: keyPair,
			Label:   "my rsa key",
		}

		repository := NewInMemorySignatureDeviceRepository()
		repository.devices[device.ID] = device

		newSignature := "new-signature"
		err = repository.MarkSignatureCreated(id, newSignature)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		got, ok := repository.devices[device.ID]
		if !ok {
			t.Error("device not found")
		}
		if got.SignatureCounter != 1 {
			t.Errorf("expected counter to be incremented to 1, got %d", got.SignatureCounter)
		}
		if got.LastSignature != newSignature {
			t.Errorf("expected last signature to be updated to '%s', got '%s'", newSignature, got.LastSignature)
		}
	})

	t.Run("returns error when device with id is not found", func(t *testing.T) {
		id := uuid.New()
		repository := NewInMemorySignatureDeviceRepository()
		err := repository.MarkSignatureCreated(id, "some-signature")
		if err == nil {
			t.Error("expected error when updating non-existent device")
		}
	})
}

func TestFind(t *testing.T) {
	t.Run("returns the device when device with id exists", func(t *testing.T) {
		device := domain.SignatureDevice{
			ID:    uuid.New(),
			Label: "my rsa key",
		}

		repository := NewInMemorySignatureDeviceRepository()
		repository.devices[device.ID] = device

		foundDevice, found, err := repository.Find(device.ID)
		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}
		if !found {
			t.Error("expected device to be found")
		}
		diff := cmp.Diff(foundDevice, device)
		if diff != "" {
			t.Errorf("unexpected difference between original and found device: %s", diff)
		}
	})

	t.Run("returns false when device with id does not exist", func(t *testing.T) {
		id := uuid.New()
		repository := NewInMemorySignatureDeviceRepository()

		_, found, err := repository.Find(id)
		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}
		if found {
			t.Error("expected found: false")
		}
	})
}

func TestList(t *testing.T) {
	repository := NewInMemorySignatureDeviceRepository()

	// create rsa device
	rsaDevice, err := domain.BuildSignatureDevice(
		uuid.New(),
		crypto.RSAGenerator{},
		"my rsa key",
	)
	if err != nil {
		t.Fatal(err)
	}
	repository.devices[rsaDevice.ID] = rsaDevice

	// create ecc device
	eccDevice, err := domain.BuildSignatureDevice(
		uuid.New(),
		crypto.ECCGenerator{},
		"my ecc key",
	)
	if err != nil {
		t.Fatal(err)
	}
	repository.devices[eccDevice.ID] = eccDevice

	got, err := repository.List()
	if err != nil {
		t.Error(err)
	}
	if len(got) != 2 {
		t.Errorf("expected 2 devices, got %d", len(got))
	}

	if got[0] != rsaDevice && got[1] != rsaDevice {
		t.Error("expected got to contain rsa device")
	}
	if got[0] != eccDevice && got[1] != eccDevice {
		t.Error("expected got to contain ecc device")
	}
}
