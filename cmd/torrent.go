package cmd

import (
	"time"

	"github.com/anacrolix/torrent/metainfo"
	"github.com/spf13/cobra"
)

// Torrent creates torrents out of files present
var Torrent = &cobra.Command{
	Use:   "torrent",
	Short: "Torrent docker images via BitTorrent Protocol",
	Long:  `Torrent docker images via BitTorrent Protocol`,
	Run:   AttachHandler(doTorrent),
}

func init() {
	Dops.AddCommand(Torrent)
}

func doTorrent(args []string) error {
	path := args[0]
	_, err := createTorrent(path)

	return err
}

func createTorrent(path string) (*metainfo.MetaInfo, error) {
	torrent := new(metainfo.MetaInfo)
	// TODO - Compute the Piece Length in order to keep the total pieces optimal
	// Ref - http://web.archive.org/web/20160619024144/https://torrentfreak.com/how-to-make-the-best-torrents-081121/
	// Read the section on "PIECE SIZE"
	torrent.Info.PieceLength = 1 * 1024 * 1024 // 1MB chunks
	err := torrent.Info.BuildFromFilePath(path)
	if err != nil {
		return nil, err
	}
	torrent.Info.UpdateBytes()
	torrent.Comment = "Created from " + path
	torrent.CreatedBy = "DOps - https://github.com/ashwanthkumar/dops"
	torrent.CreationDate = time.Now().Unix()
	return torrent, nil
}
