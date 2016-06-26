package storage

import (
	"fmt"
	"sync"
	"time"
)

// InMemoryState is a storage implementation thats hold all data in memory
type InMemoryState struct {
	data map[string]string
	lock sync.Mutex
}

func NewInMemoryState(config interface{}) *InMemoryState {
	return &InMemoryState{
		data: make(map[string]string),
	}
}

func (in *InMemoryState) IsImageExist(image string) (bool, error) {
	in.lock.Lock()
	defer in.lock.Unlock()

	_, exist := in.data[image]
	return exist, nil
}

func (in *InMemoryState) IsLayerSeeding(image, digest string) (bool, error) {
	in.lock.Lock()
	defer in.lock.Unlock()

	key := fmt.Sprintf("%s_%s_seeding", image, digest)
	_, exist := in.data[key]
	return exist, nil
}

func (in *InMemoryState) GetSeedStartTime(image, digest string) (time.Time, error) {
	in.lock.Lock()
	defer in.lock.Unlock()

	key := fmt.Sprintf("%s_%s_seeding", image, digest)
	ts, _ := in.data[key]
	return time.Parse(time.RFC3339, ts)
}

func (in *InMemoryState) StartedSeeding(image, digest string, ts time.Time) error {
	in.lock.Lock()
	defer in.lock.Unlock()

	in.data[image] = image
	key := fmt.Sprintf("%s_%s_seeding", image, digest)
	in.data[key] = ts.Format(time.RFC3339)
	return nil
}

func (in *InMemoryState) StoppedSeeding(image, digest string, ts time.Time) error {
	in.lock.Lock()
	defer in.lock.Unlock()

	key := fmt.Sprintf("%s_%s_seeding", image, digest)
	delete(in.data, key)
	return nil
}

func (in *InMemoryState) Remove(image string) error {
	in.lock.Lock()
	defer in.lock.Unlock()

	delete(in.data, image)
	return nil
}
