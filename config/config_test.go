package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConfigToTrackerAnnounceURLFromListenAddr(t *testing.T) {
	config := Config{
		ListenAddr: "my-server:8080",
	}

	announceURL, err := config.TrackerAnnounceURL()
	assert.NoError(t, err)
	assert.Equal(t, "http://my-server:8080/announce", announceURL)
}

func TestConfigToTrackerAnnounceURLFromPublicHost(t *testing.T) {
	config := Config{
		ListenAddr: ":8080",
		PublicHost: "public-fqdn",
	}

	announceURL, err := config.TrackerAnnounceURL()
	assert.NoError(t, err)
	assert.Equal(t, "http://public-fqdn:8080/announce", announceURL)
}

func TestConfigToTorrentConfig(t *testing.T) {
	config := Config{
		Debug:      true,
		DataDir:    "/foo/bar/data-dir",
		DisableDHT: true,
	}

	tConfig := config.ToTorrentConfig()
	assert.Equal(t, "/foo/bar/data-dir", tConfig.DataDir)
	assert.Equal(t, true, tConfig.NoDHT)
	assert.Equal(t, true, tConfig.Debug)
}

func TestConfigImageSeedDuration(t *testing.T) {
	config := Config{
		SeedDuration: "200h",
	}
	duration, err := config.ImageSeedDuration()
	assert.NoError(t, err)
	assert.Equal(t, time.Duration(200)*time.Hour, duration)
}

func TestConfigImageSeedDurationWhenEmpty(t *testing.T) {
	config := Config{
		SeedDuration: "",
	}
	duration, err := config.ImageSeedDuration()
	assert.NoError(t, err)
	assert.Equal(t, DefaultSeedDuration, duration)
}
