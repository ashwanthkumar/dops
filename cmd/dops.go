package ops

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

// Dops is the entry point command in the application
var Dops = &cobra.Command{
	Use:   "dops",
	Short: "Client to pull docker images via Torrent",
	Long:  `Client to pull docker images via Torrent`,
}

// CommandHandler is the wrapper interface that all commands to be implement as part of their "Run"
type CommandHandler func(args []string) error

// AttachHandler is a wrapper method for all commands that needs to be exposed
func AttachHandler(handler CommandHandler) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		err := handler(args)
		if err != nil {
			log.Printf("[Error] %s", err.Error())
			os.Exit(1)
		}
	}
}
