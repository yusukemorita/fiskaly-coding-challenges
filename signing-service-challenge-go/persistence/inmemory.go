package persistence

import (
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/google/uuid"
)

type InMemorySignatureDeviceRepositoryProvider struct {
	repository InMemorySignatureDeviceRepository
	mutex      *sync.RWMutex
}

// Use when any of the repository methods in do() write
func (provider InMemorySignatureDeviceRepositoryProvider) WriteTx(do func(domain.SignatureDeviceRepository) error) error {
	provider.mutex.Lock()
	defer provider.mutex.Unlock()
	return do(provider.repository)
}

// Use when none of the repository methods in do() write
func (provider InMemorySignatureDeviceRepositoryProvider) ReadTx(do func(domain.SignatureDeviceRepository) error) error {
	provider.mutex.RLock()
	defer provider.mutex.RUnlock()
	return do(provider.repository)
}

func NewInMemorySignatureDeviceRepositoryProvider(repository InMemorySignatureDeviceRepository) InMemorySignatureDeviceRepositoryProvider {
	return InMemorySignatureDeviceRepositoryProvider{
		repository: repository,
		mutex:      &sync.RWMutex{},
	}
}

type InMemorySignatureDeviceRepository struct {
	devices map[uuid.UUID]domain.SignatureDevice
}

func (repository InMemorySignatureDeviceRepository) Create(device domain.SignatureDevice) error {
	_, ok := repository.devices[device.ID]
	if ok {
		return errors.New(fmt.Sprintf("duplicate id: %s", device.ID))
	}

	repository.devices[device.ID] = device
	return nil
}

func (repository InMemorySignatureDeviceRepository) MarkSignatureCreated(deviceID uuid.UUID, newSignature string) error {
	device, ok := repository.devices[deviceID]
	if !ok {
		return errors.New("cannot update signature device that does not exist")
	}

	device.SignatureCounter++
	device.LastSignature = newSignature
	repository.devices[deviceID] = device

	log.Printf("updated device id: %s counter: %d", device.ID, device.SignatureCounter)

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
