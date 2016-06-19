package cmd

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/chihaya/chihaya"
	"github.com/chihaya/chihaya/server"
	"github.com/chihaya/chihaya/tracker"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"

	// Servers
	_ "github.com/ashwanthkumar/dops/tracker"

	// Middleware
	_ "github.com/chihaya/chihaya/middleware/deniability"
	_ "github.com/chihaya/chihaya/middleware/varinterval"
	_ "github.com/chihaya/chihaya/server/store/middleware/client"
	_ "github.com/chihaya/chihaya/server/store/middleware/infohash"
	_ "github.com/chihaya/chihaya/server/store/middleware/ip"
	_ "github.com/chihaya/chihaya/server/store/middleware/response"
	_ "github.com/chihaya/chihaya/server/store/middleware/swarm"
)

// Tracker starts an internal BitTorrent tracker which is supposed to be used in
// conjunction with the dopsctl command
var Tracker = &cobra.Command{
	Use:   "tracker",
	Short: "Start a BitTorrent Tracker",
	Long:  `Start a BitTorrent Tracker`,
	Run:   AttachHandler(startTracker),
}

func init() {
	Dops.AddCommand(Tracker)
}

func startTracker(args []string) error {
	configPath, err := homedir.Expand("~/.dops.yml")
	config, err := chihaya.OpenConfigFile(configPath)
	if err != nil {
		return err
	}
	tkr, err := tracker.NewTracker(&config.Tracker)
	if err != nil {
		return err
	}

	log.Printf("Starting BitTorrent Tracker and DOps torrent registry")
	p, err := server.StartPool(config.Servers, tkr)
	if err != nil {
		return err
	}

	shutdown := make(chan os.Signal)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown
	p.Stop()

	return nil
}
