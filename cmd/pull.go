package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/cheggaaa/pb.v1"

	"golang.org/x/net/context"

	"github.com/ashwanthkumar/dops/docker"
	"github.com/spf13/cobra"
)

// Pull executes the torrent based fs layer downloads of a docker image
var Pull = &cobra.Command{
	Use:   "pull",
	Short: "Pull docker images via BitTorrent Protocol",
	Long:  `Pull docker images via BitTorrent Protocol`,
	Run:   AttachHandler(doPull),
}

func init() {
	Dops.AddCommand(Pull)
}

func doPull(args []string) error {
	image := args[0]
	insecure := false

	named, manifest, err := docker.DownloadManifest(image, insecure)
	if err != nil {
		return err
	}

	repo, err := docker.GetRepositoryClient(named, insecure, "pull")
	if err != nil {
		return err
	}

	ctx := context.Background()
	blobSvc := repo.Blobs(ctx)
	downloadPath := "/tmp/" + repo.Named().Name()
	if err := os.MkdirAll(downloadPath, os.ModePerm); err != nil {
		return err
	}

	for _, descriptor := range manifest.References() {
		descriptor, _ := blobSvc.Stat(ctx, descriptor.Digest)
		// TODO - See if we skip downloading if the file exists and matches the checksum
		f, err := os.Create(filepath.Join(downloadPath, descriptor.Digest.Hex()))
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		reader, err := blobSvc.Open(ctx, descriptor.Digest)
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		writer := bufio.NewWriter(f)
		bar := pb.New64(descriptor.Size).SetUnits(pb.U_BYTES).Prefix(descriptor.Digest.Hex()[0:15])
		bar.Start()
		proxyWriter := io.MultiWriter(writer, bar)

		if length, err := io.Copy(proxyWriter, reader); length != descriptor.Size || err != nil {
			if err != nil {
				return err
			}
			return fmt.Errorf("Download incomplete, expected %d but got only %d\n", descriptor.Size, length)
		}
		bar.Finish()
		// log.Printf("Finished downloading")
	}

	return nil
}
