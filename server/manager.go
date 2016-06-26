package server

import (
	"time"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
)

// TorrentManager - Manages all the torrents which needs to be seeded
type TorrentManager struct {
	client      *torrent.Client
	config      *Config
	announceURL string
}

// Init - Initializes the TorrentManager - Should be called only once during the life-time
func (t *TorrentManager) Init(config *Config) (err error) {
	t.config = config
	t.client, err = torrent.NewClient(config.ToTorrentConfig())
	t.announceURL, err = config.TrackerAnnounceURL()
	return err
}

// Start starts the TorrentManager
func (t *TorrentManager) Start() {
	// TODO - Have a reaper to check the client and remove the torrent's after the seed duration
	// go t.startTorrentReaper()
}

func (t *TorrentManager) startTorrentReaper() {
	for _ = range time.After(10 * time.Minute) {
		for _ = range t.client.Torrents() {
			// TODO Check if the torrent t, has been seeding for over config.ImageSeedDuration
			// if so, drop it, else keep it running
		}
	}
}

// AddForSeeding - Creates a torrent out for the file and start seeding it
// it returns the MagnetURI and an error object
func (t *TorrentManager) AddForSeeding(path string) (string, error) {
	metaInfo, err := CreateTorrent(path, t.announceURL)
	if err != nil {
		return "", err
	}
	torrent, err := t.client.AddTorrent(metaInfo)
	if err != nil {
		return "", err
	}
	<-torrent.GotInfo()
	torrent.DownloadAll()
	return torrent.Metainfo().Magnet().String(), nil
}

// Stop - Stops the underlying torrent client
func (t *TorrentManager) Stop() error {
	t.client.Close()
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
