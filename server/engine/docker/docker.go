package docker

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"

	"github.com/ashwanthkumar/dops/config"
	"github.com/ashwanthkumar/dops/server/torrent"
	"github.com/ashwanthkumar/golang-utils/worker"
	"github.com/docker/distribution"
	"github.com/docker/docker/reference"
	"golang.org/x/net/context"
)

// DockerEngine is an implementation of Container Engine
type DockerEngine struct {
	DataDir    string
	TorrentMgr *torrent.Manager
	worker     worker.Pool
}

type DownloadWork struct {
	named        reference.Named
	insecure     bool
	downloadPath string
	descriptor   *distribution.Descriptor
	ctx          context.Context
}

func (e DockerEngine) Init(config *config.Config, torrentMgr *torrent.Manager) error {
	e.DataDir = config.DataDir
	e.TorrentMgr = torrentMgr
	e.worker = worker.Pool{
		MaxWorkers: 5 * runtime.NumCPU(), // TODO - Make this configurable
		Op: func(work worker.Request) error {
			request := work.(DownloadWork)

			repo, err := GetRepositoryClient(request.named, request.insecure, "pull")
			if err != nil {
				return err
			}

			ctx := context.Background()
			blobSvc := repo.Blobs(ctx)
			downloadPath := filepath.Join(e.DataDir, repo.Named().Name())
			if err := os.MkdirAll(downloadPath, os.ModePerm); err != nil {
				return err
			}

			info, _ := blobSvc.Stat(request.ctx, request.descriptor.Digest)
			// TODO - See if we can skip downloading if the file exists and matches the checksum
			f, err := os.Create(filepath.Join(request.downloadPath, info.Digest.Hex()))
			if err != nil {
				return err
			}
			reader, err := blobSvc.Open(ctx, info.Digest)
			if err != nil {
				return err
			}

			writer := bufio.NewWriter(f)
			if length, err := io.Copy(writer, reader); length != info.Size || err != nil {
				if err != nil {
					return err
				}
				return fmt.Errorf("Download incomplete for %s, expected %d but got only %d\n", info.Digest.String(), info.Size, length)
			}

			return nil
		},
	}
	return nil
}

func (e DockerEngine) DownloadImage(image string, config *config.Config) error {
	insecure := false // TODO - make this configurable

	named, manifest, err := DownloadManifest(image, insecure)
	if err != nil {
		return err
	}

	for _, descriptor := range manifest.References() {
		work := DownloadWork{
			named:      named,
			insecure:   insecure,
			descriptor: &descriptor,
		}
		e.worker.AddWork(work)
	}

	return nil
}
