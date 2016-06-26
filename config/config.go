package config

import (
	"fmt"
	"net/url"
	"time"

	"github.com/anacrolix/torrent"
	"github.com/ashwanthkumar/golang-utils/netutil"
)

const (
	// AnnounceURI is the URI underwhich the tracker's announce requests are handled
	// all generated torrents / magnet uris should have a full uri path to this
	AnnounceURI = "/announce"
	// ScrapeURI is the URI underwhich the tracker's scrape requests are handled
	ScrapeURI = "/scrape"
)

var (
	// DefaultSeedDuration is the default amount of duration the images are cached on the registry
	// after which they'll no longer be seeded actively, you can restart seeding by using `dopsctl cache`
	// Value is 7 * 24 hours (aka) 7 days
	DefaultSeedDuration = time.Duration(7*24) * time.Hour
)

// Config - Server settings
type Config struct {
	DataDir    string `json:"data-dir"`
	Debug      bool   `json:"debug,omitempty"`
	DisableDHT bool   `json:"disable-dht,omitempty"`
	// this is the host:port combination where the Tracker will also run
	ListenAddr string `json:"listen-addr"`
	// Public facing hostname of the instance running the DOPS Registry
	PublicHost string `json:"public-host,omitempty"`
	//  Duration for which we should be seeding
	SeedDuration  string      `json:"seed-duration,omitempty"`
	StorageType   string      `json:"storage-type,omitempty"`
	StorageConfig interface{} `json:"storage,omitempty"`
}

// ToTorrentConfig converts our config to anacrolix's torrent config representation
func (c *Config) ToTorrentConfig() *torrent.Config {
	return &torrent.Config{
		Debug:   c.Debug,
		DataDir: c.DataDir,
		NoDHT:   c.DisableDHT,
	}
}

// TrackerAnnounceURL - Generates a URL for the server tracker's announce url
func (c *Config) TrackerAnnounceURL() (string, error) {
	hostname, port, err := netutil.SplitHostPort(c.ListenAddr)
	if err != nil {
		return "", err
	}
	if hostname == "" {
		hostname = c.PublicHost
	}

	if hostname == "" {
		hostname, err = netutil.FullyQualifiedHostname()
		if err != nil {
			return "", err
		}
	}

	baseURL, err := url.Parse(fmt.Sprintf("http://%s:%d", hostname, port))
	if err != nil {
		return "", err
	}
	uri, err := url.Parse(fmt.Sprintf("//%s", AnnounceURI))
	if err != nil {
		return "", err
	}

	return baseURL.ResolveReference(uri).String(), nil
}

// ImageSeedDuration - Duration for which each FSLayer blob is seeded
func (c *Config) ImageSeedDuration() (time.Duration, error) {
	if c.SeedDuration != "" {
		return time.ParseDuration(c.SeedDuration)
	}
	return DefaultSeedDuration, nil
}
