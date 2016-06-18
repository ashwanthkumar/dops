package main

import (
	"fmt"
	"os"

	"github.com/ashwanthkumar/dops/cmd"
)

// Version of the app
var Version = "dev-build"

func main() {
	if err := cmd.Dops.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
