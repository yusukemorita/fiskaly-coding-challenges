package persistence

import (
	"errors"
	"fmt"
	"sync"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/google/uuid"
)

type InMemorySignatureDeviceRepository struct {
	devices map[uuid.UUID]domain.SignatureDevice
	lock    *sync.RWMutex
}

func (repository InMemorySignatureDeviceRepository) Create(device domain.SignatureDevice) error {
	_, ok := repository.devices[device.ID]
	if ok {
		return errors.New(fmt.Sprintf("duplicate id: %s", device.ID))
	}

	repository.devices[device.ID] = device
	return nil
}

func (repository InMemorySignatureDeviceRepository) Update(device domain.SignatureDevice) error {
	_, ok := repository.devices[device.ID]
	if !ok {
		return errors.New("cannot update signature device that does not exist")
	}

	repository.devices[device.ID] = device
	return nil
}

func (repository InMemorySignatureDeviceRepository) Find(id uuid.UUID) (domain.SignatureDevice, bool, error) {
	device, ok := repository.devices[id]
	if !ok {
		return domain.SignatureDevice{}, false, nil
	}

	return device, true, nil
}

// Order is not guaranteed
// ref: https://go.dev/blog/maps
// > When iterating over a map with a range loop, the iteration order is not specified and is not guaranteed
// > to be the same from one iteration to the next. If you require a stable iteration order you must maintain
// > a separate data structure that specifies that order.
func (repository InMemorySignatureDeviceRepository) List() ([]domain.SignatureDevice, error) {
	allDevices := []domain.SignatureDevice{}

	for _, device := range repository.devices {
		allDevices = append(allDevices, device)
	}

	return allDevices, nil
}

func NewInMemorySignatureDeviceRepository() InMemorySignatureDeviceRepository {
	return InMemorySignatureDeviceRepository{
		devices: map[uuid.UUID]domain.SignatureDevice{},
	}
}
