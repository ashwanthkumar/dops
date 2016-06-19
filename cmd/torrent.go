package cmd

import (
	"fmt"
	"time"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/dht"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/spf13/cobra"
)

var port int32
var dataDir string
var disableTrackers bool
var debug bool
var seed bool
var noUpload bool
var dhtDisableBootstrap bool
var dhtBootstrapNodes []string

// Torrent creates torrents out of files present
var Torrent = &cobra.Command{
	Use:   "torrent",
	Short: "Torrent docker images via BitTorrent Protocol",
	Long:  `Torrent docker images via BitTorrent Protocol`,
	Run:   AttachHandler(doTorrent),
}

func init() {
	Dops.AddCommand(Torrent)
	Torrent.PersistentFlags().Int32Var(&port, "port", 50007, "Port number to listen to for the torrent client")
	Torrent.PersistentFlags().StringVar(&dataDir, "data-dir", "/tmp/dops", "Temp directory where the torrent files are downloaded")
	Torrent.PersistentFlags().BoolVar(&disableTrackers, "disable-trackers", false, "Disables announcing to Trackers and relies only on DHT")
	Torrent.PersistentFlags().BoolVar(&debug, "debug", false, "Enable Debugging mode")
	Torrent.PersistentFlags().BoolVar(&seed, "seed", true, "Enable seeding the torrent we've downloaded")
	Torrent.PersistentFlags().BoolVar(&noUpload, "no-upload", false, "Disable uploading chunks to peers, even while downloading")
	Torrent.PersistentFlags().BoolVar(&dhtDisableBootstrap, "dht-disable-bootstrap", false, "Disable bootstrapping from global servers even if given no BootstrapNodes")
	Torrent.PersistentFlags().StringSliceVar(&dhtBootstrapNodes, "dht-bootstrap-nodes", []string{}, "DHT Bootstrap nodes")
}

func doTorrent(args []string) error {
	path := args[0]
	cfg := DefaultTorrentConfig()
	client, err := torrent.NewClient(cfg)
	if err != nil {
		return err
	}
	torrent, err := CreateTorrent(path)
	if err != nil {
		return err
	}

	t, err := client.AddTorrent(torrent)
	if err != nil {
		return err
	}

	<-t.GotInfo()
	t.DownloadAll()

	fmt.Printf("%v\n", client.ListenAddr())

	client.WaitAll()

	return err
}

// CreateTorrent - Generates a torrent metainfo from a file path and recursively inside it
func CreateTorrent(path string) (*metainfo.MetaInfo, error) {
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
	torrent.Info.UpdateBytes()
	return torrent, nil
}

// DefaultTorrentConfig provides te default settings for a torrent client instance
func DefaultTorrentConfig() *torrent.Config {
	listenAddr := fmt.Sprintf(":%d", port)
	dhtConfig := dht.ServerConfig{
		NoDefaultBootstrap: dhtDisableBootstrap,
		BootstrapNodes:     dhtBootstrapNodes,
	}
	cfg := &torrent.Config{
		DataDir:         dataDir,
		DisableTrackers: disableTrackers,
		Seed:            seed,
		NoUpload:        noUpload,
		Debug:           debug,
		ListenAddr:      listenAddr,
		DHTConfig:       dhtConfig,
	}

	return cfg
}
