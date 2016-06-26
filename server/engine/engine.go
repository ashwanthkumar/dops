package server

import (
	"fmt"

	"github.com/anacrolix/missinggo/pubsub"
	"github.com/ashwanthkumar/dops/config"
	"github.com/ashwanthkumar/dops/server/engine/docker"
)

// ContainerEngine - Implementation responsible for downloading images, loading them into the
// engine, etc. Currently we only have doker implementation
type ContainerEngine interface {
	Init(config *config.Config) error
	// ScheduleImage schedules downloading of all the layers for the image in the Background
	// it returns the list of layers (Hex value) that'll be downloaded. One can subscribe
	// to the layer downloads using SubscribeDownloads()
	ScheduleImage(image string, config *config.Config) ([]string, error)
	SubscribeDownloads() *pubsub.Subscription
}

// Since ContainerEngines are stateful objects, we need to re-use them across
var engines map[string]ContainerEngine

// GetEngine returns an implementation of ContainerEngine identified by name
func GetEngine(name string, config *config.Config) (ContainerEngine, error) {
	engine, exist := engines[name]
	if exist {
		return engine, nil
	}

	switch name {
	case docker.Name:
		engine := &docker.DockerEngine{}
		if err := engine.Init(config); err != nil {
			return nil, err
		}
		engines[docker.Name] = engine
		return engine, nil
	default:
		return nil, fmt.Errorf("%s is not a valid implementation\n", name)
	}
}
