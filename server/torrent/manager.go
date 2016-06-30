package torrent

import (
	"time"

	torrentlib "github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/ashwanthkumar/dops/config"
	"github.com/ashwanthkumar/dops/server/storage"
)

// Manager - Manages all the torrents which needs to be seeded
type Manager struct {
	config      *config.Config
	storage     storage.Storage
	client      *torrentlib.Client
	announceURL string
}

// New creates new Manager instance
func New(config *config.Config) (*Manager, error) {
	var err error
	manager := &Manager{}
	manager.config = config
	if manager.storage, err = storage.GetStorage(config); err != nil {
		return nil, err
	}
	if manager.client, err = torrentlib.NewClient(config.ToTorrentConfig()); err != nil {
		return nil, err
	}
	if manager.announceURL, err = config.TrackerAnnounceURL(); err != nil {
		return nil, err
	}
	return manager, nil
}

// Start starts the Manager
func (manager *Manager) Start() {
	// TODO - Have a reaper to check the client and remove the torrent's after the seed duration
	// go t.startTorrentReaper()
}

// AddForSeeding - Creates a torrent out for the file and start seeding it
// it returns the MagnetURI and an error object
func (manager *Manager) AddForSeeding(path string) (string, error) {
	metaInfo, err := CreateTorrent(path, manager.announceURL)
	if err != nil {
		return "", err
	}
	torrent, err := manager.client.AddTorrent(metaInfo)
	if err != nil {
		return "", err
	}
	<-torrent.GotInfo()
	torrent.DownloadAll()
	return torrent.Metainfo().Magnet().String(), nil
}

// Torrents returns all the torrents that the underlying client is managing
func (manager *Manager) Torrents() []*torrentlib.Torrent {
	return manager.client.Torrents()
}

// RemoveTorrent removes a torrent from the underlying torrent client
func (manager *Manager) RemoveTorrent(torrent *torrentlib.Torrent) bool {
	t, exist := manager.client.Torrent(torrent.InfoHash())
	if exist {
		t.Drop()
		return true
	}

	return false
}

// Stop - Stops the underlying torrent client
func (manager *Manager) Stop() error {
	manager.client.Close()
	return nil
}

// CreateTorrent - Generates a torrent metainfo from a file path and recursively inside it
func CreateTorrent(path, announceURL string) (*metainfo.MetaInfo, error) {
	torrent := new(metainfo.MetaInfo)
	// TODO - Compute the Piece Length in order to keep the total pieces optimal
	// Ref - http://web.archive.org/web/20160619024144/https://torrentfreak.com/how-to-make-the-best-torrents-081121/
	// Read the section on "PIECE SIZE"
	torrent.Info.PieceLength = 1 * 1024 * 1024 // 1MB chunks
	err := torrent.Info.BuildFromFilePath(path)
	if err != nil {
		return nil, err
	}
	torrent.Comment = "Created from " + path
	torrent.CreatedBy = "DOps - https://github.com/ashwanthkumar/dops"
	torrent.CreationDate = time.Now().Unix()
	torrent.Announce = announceURL
	torrent.Info.UpdateBytes()
	return torrent, nil
}
