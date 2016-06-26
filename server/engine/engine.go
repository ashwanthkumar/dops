package server

import (
	"fmt"

	"github.com/ashwanthkumar/dops/config"
	"github.com/ashwanthkumar/dops/server/engine/docker"
	"github.com/ashwanthkumar/dops/server/torrent"
)

// ContainerEngine - Implementation responsible for downloading images, loading them into the
// engine, etc. Currently we only have doker implementation
type ContainerEngine interface {
	Init(config *config.Config, torrentMgr *torrent.Manager) error
	DownloadImage(image string, config *config.Config) error
}

// Since ContainerEngines are stateful objects, we need to re-use them across
var engines map[string]ContainerEngine

func GetEngine(config *config.Config, name string, torrentMgr *torrent.Manager) (ContainerEngine, error) {
	engine, exist := engines[name]
	if exist {
		return engine, nil
	}

	switch name {
	case "docker":
		engine := docker.DockerEngine{}
		if err := engine.Init(config, torrentMgr); err != nil {
			return nil, err
		}
		engines["docker"] = engine
		return engine, nil
	default:
		return nil, fmt.Errorf("%s is not a valid implementation\n", name)
	}
}
