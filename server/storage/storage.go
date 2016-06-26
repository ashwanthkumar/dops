package storage

import (
	"fmt"
	"time"

	"github.com/ashwanthkumar/dops/config"
)

// Storage handles all the state using a persistent storage underneath
type Storage interface {
	// Do we've the image locally cached?
	IsImageExist(image string) (bool, error)
	// Are we seeding the digest from the image?
	IsLayerSeeding(image, digest string) (bool, error)
	// Time when the seeding started for the layer,
	// you need to gaurd the call to this with #IsLayerSeeding
	GetSeedStartTime(image, digest string) (time.Time, error)
	// Started seeding the digest from the image at the given ts
	// If ts is nil, we set it as time.Now()
	StartedSeeding(image, digest string, ts time.Time) error
	// Stoped seeding the digest from the image at the given ts
	// If ts is nil, we set it as time.Now()
	StoppedSeeding(image, digest string, ts time.Time) error
	// Remove all references of the given image from our memory
	Remove(image string) error
}

func GetStorage(config *config.Config) (Storage, error) {
	switch config.StorageType {
	case "memory":
		return NewInMemoryState(config.StorageConfig), nil
	default:
		return nil, fmt.Errorf("%s is a valid storage type", config.StorageType)
	}
}
